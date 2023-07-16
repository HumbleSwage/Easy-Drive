package service

import (
	"easy-drive/consts"
	"easy-drive/pkg/ctl"
	"easy-drive/pkg/e"
	"easy-drive/pkg/utils"
	"easy-drive/pkg/utils/fileUtil"
	"easy-drive/repositry/dao"
	"easy-drive/types"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"sync"
)

var RecycleSrv *RecycleService
var RecycleSrvOnce sync.Once

type RecycleService struct {
}

func GetRecycleSrv() *RecycleService {
	RecycleSrvOnce.Do(func() {
		RecycleSrv = &RecycleService{}
	})
	return RecycleSrv
}

func (rs *RecycleService) LoadRecycleListService(ctx *gin.Context, req *types.LoadReRecycleListReq) (resp interface{}, err error) {
	code := e.Success

	// session对象
	session := sessions.Default(ctx)

	// 获取用户id
	var userId string
	if u := session.Get(consts.UserInfo); u != nil {
		userId = u.(string)
	} else {
		code = e.UserSessionExpiration
		return ctl.RespError(code), nil
	}

	// 分页
	// 单页数量
	pageNoStr := strings.TrimSpace(req.PageNo)
	if pageNoStr == "" {
		pageNoStr = "0"
	}
	pageNumStr := strings.TrimSpace(req.PageSize)
	if pageNumStr == "" {
		pageNumStr = "15"
	}
	pageNum, err := strconv.Atoi(pageNoStr)
	pageSize, err := strconv.Atoi(pageNumStr)
	if err != nil {
		utils.LogrusObj.Error("分页参数转换错误:", err)
		return ctl.RespError(), nil
	}

	// 获取用户回收站中的文件
	fileDao := dao.NewFileDao(ctx.Request.Context())
	files, totalCount, err := fileDao.GetDelFile(userId, pageNum, pageSize)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("查询用户回收站的文件出错：", err)
			return ctl.RespError(), nil
		}
	}

	// 响应数据
	fileArray := make([]*types.FileInfoResp, 0)
	for _, file := range files {
		d := &types.FileInfoResp{
			FileId:         file.FileId,
			FilePid:        file.FileId,
			FileSize:       file.Size,
			FileName:       file.Name,
			FileCover:      file.Cover,
			CreateTime:     file.CreatedAt.Format("2006-01-02 15:04:05"),
			LastUpdateTime: file.UpdatedAt.Format("2006-01-02 15:04:05"),
			RecoveryTime:   file.RestoredAt.Format("2006-01-02 15:04:05"),
			FolderType:     cast.ToInt(file.IsDirectory),
			FileCategory:   file.Category,
			FileType:       file.Type,
			Status:         file.Status,
		}
		fileArray = append(fileArray, d)
	}

	// 绑定返回体
	data := types.LoadDataListResp{
		TotalCount: totalCount,
		PageSize:   int64(pageSize),
		PageNo:     int64(pageNum),
		PageTotal:  int64(len(fileArray)),
		List:       fileArray,
	}

	return ctl.RespSuccessWithData(data), nil
}

func (rs *RecycleService) RecoverFileService(ctx *gin.Context, req *types.RecoverFileReq) (resp interface{}, err error) {
	code := e.Success

	// session对象
	session := sessions.Default(ctx)

	// 获取用户id
	var userId string
	if u := session.Get(consts.UserInfo); u != nil {
		userId = u.(string)
	} else {
		code = e.UserSessionExpiration
		return ctl.RespError(code), nil
	}

	// 获取fileId
	fileArray := strings.Split(req.FileIds, ",")

	// 更新文件状态
	fileDao := dao.NewFileDao(ctx.Request.Context())
	for _, fileId := range fileArray {
		file, err := fileDao.GetFileById(fileId, userId)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				utils.LogrusObj.Error("用户id和fileId查询文件出错：", err)
				return ctl.RespError(), nil
			}
		}

		// 不是删除文件
		if file.Flag == 2 {
			continue
		}

		// 是文件夹
		if file.IsDirectory {
			// 获取文件夹下面所有文件
			filesInFolder, err := fileDao.GetFileInSameFolder(userId, []string{file.FileId})
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					utils.LogrusObj.Error("查询文件夹下文件出错：", err)
					return ctl.RespError(), nil
				}
			}

			// 更新这些文件的状态
			for _, f := range filesInFolder {
				f.Flag = consts.NormalFile.Index()

				// 更新文件
				err := fileDao.UpdateFile(f)
				if err != nil {
					utils.LogrusObj.Error("更新文件出错：", err)
					return ctl.RespError(), nil
				}
			}
		}

		// 更新文件状态
		file.Flag = consts.NormalFile.Index()

		// 检查文件或文件夹重命名的问题
		f, err := fileDao.GetFileByFileName(file.Name, file.ParentId, userId)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				utils.LogrusObj.Error("根据文件名查询文件出错:")
				return ctl.RespError(), nil
			}
		} else {
			// 检查到已有文件名
			fileNameNoExt := fileUtil.GetFileNameWithoutExtension(f.Name)

			// 解析文件名和当前索引
			currentIndex := fileUtil.GetFileNameIndex(f.Name)
			newIndex := currentIndex + 1

			// 构造新的文件名
			newFileName := fmt.Sprintf("%s(%d).txt", fileNameNoExt, newIndex)

			// 检查新文件名是否已存在
			for {
				f, err = fileDao.GetFileByFileName(newFileName, file.ParentId, userId)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						break
					} else {
						utils.LogrusObj.Error("根据文件名查询文件出错:", err)
						return ctl.RespError(), nil
					}
				}
				// 递增索引，重新构造新文件名
				newIndex++
				newFileName = fmt.Sprintf("%s(%d).txt", fileNameNoExt, newIndex)
			}
			// 绑定数据
			file.Name = newFileName
		}

		// 更新文件
		err = fileDao.UpdateFile(file)
		if err != nil {
			utils.LogrusObj.Error("更新文件出错：", err)
			return ctl.RespError(), nil
		}
	}

	return ctl.RespSuccess(), nil
}
