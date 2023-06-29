package service

import (
	"context"
	"easy-drive/consts"
	"easy-drive/pkg/ctl"
	"easy-drive/pkg/e"
	"easy-drive/pkg/utils"
	"easy-drive/pkg/utils/commonUtil"
	"easy-drive/pkg/utils/fileUtil"
	"easy-drive/repositry/cache"
	"easy-drive/repositry/dao"
	"easy-drive/repositry/model"
	"easy-drive/types"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"sync"
)

var FileSrv *FileService
var FileSrvOnce sync.Once

type FileService struct {
}

func GetFileSrv() *FileService {
	FileSrvOnce.Do(func() {
		FileSrv = &FileService{}
	})
	return FileSrv
}

func (fs *FileService) LoadDataService(ctx *gin.Context, req *types.LoadDataListReq) (resp interface{}, err error) {
	code := e.Success

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

	// 获取用户空间
	userDao := dao.NewUserDao(ctx.Request.Context())
	user, err := userDao.GetUserByUserId(userId)
	if err != nil {
		utils.LogrusObj.Error("根据用户id获取用户出错：", err)
		return ctl.RespError(), nil
	}

	// 查询指定用户下的指定的文件
	fileDao := dao.NewFileDao(ctx.Request.Context())

	// 分页查询category
	var files []*model.File
	var totalCount int64
	// 正常状态
	switch req.Category {
	case consts.Video.String():
		files, err = fileDao.GetUserFilesByCategory(user, consts.Video.Index(), pageNum, pageSize, totalCount)
	case consts.Music.String():
		files, err = fileDao.GetUserFilesByCategory(user, consts.Music.Index(), pageNum, pageSize, totalCount)
	case consts.Image.String():
		files, err = fileDao.GetUserFilesByCategory(user, consts.Image.Index(), pageNum, pageSize, totalCount)
	case consts.Doc.String():
		files, err = fileDao.GetUserFilesByCategory(user, consts.Doc.Index(), pageNum, pageSize, totalCount)
	case consts.Others.String():
		files, err = fileDao.GetUserFilesByCategory(user, consts.Others.Index(), pageNum, pageSize, totalCount)
	default:
		// 默认查询所有
		files, err = fileDao.GetUserFiles(user, pageNum, pageSize, totalCount)
	}

	// 没有记录就返回空
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			utils.LogrusObj.Error("根据用户查询其指定类别文件时发生错误:", err)
			return ctl.RespError(), nil
		}
	}

	// 处理反应的响应体
	l := make([]*types.FileInfoResp, 0)
	for _, file := range files {
		f := &types.FileInfoResp{
			FileId:         file.FileId,
			FilePid:        file.FileId,
			FileSize:       file.Size,
			FileName:       file.Name,
			FileCover:      file.Cover,
			CreateTime:     file.CreatedAt.Format("2006-01-02 15:04:05"),
			LastUpdateTime: file.UpdatedAt.Format("2006-01-02 15:04:05"),
			FolderType:     cast.ToInt(file.IsDirectory),
			FileCategory:   file.Category,
			FileType:       file.Type,
			Status:         file.Status,
		}
		l = append(l, f)
	}

	// 绑定返回体
	data := types.LoadDataListResp{
		TotalCount: totalCount,
		PageSize:   int64(pageSize),
		PageTotal:  int64(len(l)),
		List:       l,
	}

	// 反应响应体
	return ctl.RespSuccessWithData(data), err
}

func (fs *FileService) UploadFileService(ctx *gin.Context, req *types.UploadFileReq) (resp interface{}, err error) {
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

	// 获取用户
	userDao := dao.NewUserDao(ctx.Request.Context())
	user, err := userDao.GetUserByUserId(userId)
	if err != nil {
		utils.LogrusObj.Error("根据用户id获取用户出错：", err)
		return ctl.RespError(), nil
	}
	fileDao := dao.NewFileDaoByDB(userDao.DB)

	// 参数转换
	chunkIndex, err := strconv.Atoi(req.ChunkIndex)
	chunks, err := strconv.Atoi(req.Chunks)
	if err != nil {
		utils.LogrusObj.Error(err)
		return ctl.RespError(), nil
	}

	// 绑定参数
	newFile := &model.File{
		UserId:   user.UserId,
		MD5:      req.FileMd5,
		ParentId: req.FilePid,
	}

	// 文件id
	if req.FileId == "" {
		newFile.FileId = commonUtil.GenerateFileID()
	} else {
		newFile.FileId = req.FileId
	}

	// 获取文件名：TODO：这个地方逻辑有一点奇怪
	filePath, fileName := fileUtil.GetUniqueFileName(req.FileName, userId)
	newFile.Name = fileName

	// 获取文件类型
	t, err := fileUtil.GetFileTypeNumber(fileName)
	if err != nil {
		utils.LogrusObj.Error("获取文件类型出错:", err)
		return ctl.RespError(), nil
	}
	newFile.Type = t

	// 判断是否分片
	if chunkIndex == 0 {
		// 查询文件
		file, err := fileDao.SelectFileByMd5(req.FileMd5, consts.NormalFile.Index())
		if err != nil && err != gorm.ErrRecordNotFound {
			utils.LogrusObj.Error("通过文件md5值查询文件出错:", err)
			return ctl.RespError(), nil
		}

		// 文件秒传
		if file.FileId != "" {
			utils.LogrusObj.Infoln("文件秒传")
			// 内存空间
			if user.UseSpace+file.Size > user.TotalSpace {
				code = e.UserStoreSpaceError
				return ctl.RespError(code), nil
			}

			// 绑定数据
			newFile.Cover = file.Cover
			newFile.IsDirectory = file.IsDirectory
			newFile.Category = file.Category
			newFile.Type = file.Type
			newFile.Size = file.Size
			newFile.Status = consts.Using.Index()
			newFile.Flag = consts.NormalFile.Index()
			newFile.Path = file.Path

			// 更新用户使用空间 ：TODO:判断系统磁盘使用空间
			user.UseSpace = user.UseSpace + file.Size

			// 更新数据: TODO：数据库事务
			if err := userDao.UpdateUser(user); err != nil {
				utils.LogrusObj.Error("更新用户入库出错:", err)
				return ctl.RespError(), nil
			}
			if err := fileDao.CreateFile(newFile); err != nil {
				utils.LogrusObj.Error("创建文件出错:", err)
				return ctl.RespError(), nil
			}

			// 响应体
			data := &types.UploadFileResp{
				FileId: newFile.FileId,
				Status: consts.FastUpload.String(),
			}

			return ctl.RespSuccessWithData(data), nil
		}
	}

	// 用户内存缓存
	key := cache.UserStoreSpaceKey(userId, fileName)

	// 文件分片上传到临时文件夹
	if chunkIndex < chunks-1 {
		// 获取文件缓存
		var value int64
		if cache.RedisClient.Exists(key).Val() == 1 {
			// 已有key，获取key对应value
			value, err = cache.RedisClient.Get(key).Int64()
			if err != nil {
				utils.LogrusObj.Error("获取分块文件缓存出错:", err)
			}
		} else {
			// 无key，第一个分块，创建缓存
			_, err = cache.RedisClient.Set(key, req.File.Size, 0).Result()
			if err != nil {
				utils.LogrusObj.Error("创建分块文件缓存出错:", err)
				return nil, err
			}
		}

		// 已有key:判断用户可用内存
		space := value + user.UseSpace + req.File.Size // 第一次value为0

		// 内存不足:TODO：如果前端能一次直接将文件内存大小传过来就不需要这个操作
		if space > user.TotalSpace {
			code = e.UserStoreSpaceError
			return ctl.RespError(code), nil
		}

		// value增加
		_, err = cache.RedisClient.Set(key, value+req.File.Size, 0).Result()
		if err != nil {
			utils.LogrusObj.Error("更新分块文件缓存出错:", err)
			return nil, err
		}

		// 临时文件存储的位置
		tempFilePath := fileUtil.GetTempFilePath(fileName, userId)
		// 临时存储
		err := ctx.SaveUploadedFile(req.File, tempFilePath)
		//_, err = fileUtil.ChunkedFileToLocalTemp(userId, req.FileName, req.File, fileHeader)
		if err != nil {
			utils.LogrusObj.Error("文件分片上传到临时文件夹出错:", err)
			return ctl.RespError(), nil
		}

		// 响应体
		data := &types.UploadFileResp{
			FileId: newFile.FileId,
			Status: consts.OnUpload.String(),
		}

		return ctl.RespSuccessWithData(data), nil
	}

	if chunks == 1 {
		// 无分片，直接将资源保存到指定位置
		err = ctx.SaveUploadedFile(req.File, filePath)
		if err != nil {
			utils.LogrusObj.Error("上传资源到本地出错：", err)
		}
	} else {
		// 存储文件
		tempFilePath := fileUtil.GetTempFilePath(fileName, userId)

		// 临时存储
		err = ctx.SaveUploadedFile(req.File, tempFilePath)
		if err != nil {
			utils.LogrusObj.Error("文件分片上传到临时文件夹出错:", err)
			return ctl.RespError(), nil
		}

		// 合并分片到指定位置
		err = fileUtil.MergeChunks(userId, fileName, filePath)
		if err != nil {
			utils.LogrusObj.Error("合并分片出错:", err)
		} else {
			utils.LogrusObj.Infoln("文件合并完成:", filePath)
		}
	}

	// 绑定数据
	newFile.Path = filePath
	newFile.Flag = consts.NormalFile.Index()
	newFile.Status = consts.Transfer.Index() // 正在转码中
	newFile.Category = fileUtil.GetFileType(fileName)

	// 更新内存
	var value int64
	if cache.RedisClient.Exists(key).Val() == 1 {
		// 如果是分片，就从缓存中取出所有内存总和
		value, err = cache.RedisClient.Get(key).Int64()
		if err != nil {
			utils.LogrusObj.Error("获取文件缓存出错:", err)
			return ctl.RespError(), nil
		}
		// 最后分块，删除缓存
		cache.RedisClient.Del(key)
	}
	fileSize := value + req.File.Size
	user.UseSpace = user.UseSpace + fileSize
	newFile.Size = fileSize

	// 数据入库:TODO:保证事务的一致性
	if err = userDao.UpdateUser(user); err != nil {
		utils.LogrusObj.Error("更新用户出错:", err)
		return ctl.RespError(), nil
	}
	if err = fileDao.CreateFile(newFile); err != nil {
		utils.LogrusObj.Error("创建文件出错:", err)
		return ctl.RespError(), nil
	}

	// 针对视频和图片：启动一个协程序完成文件转码
	if fileUtil.GetFileType(fileName) == consts.Video.Index() || fileUtil.GetFileType(fileName) == consts.Image.Index() {
		go func(ctx context.Context) {
			// 更新文件状态
			newFile.Status = consts.Using.Index()

			// 缩略图名称
			thumbnailName := fileUtil.GenerateThumbnailName(fileName)
			utils.LogrusObj.Infoln("缩略图名称:", thumbnailName)

			// 生成缩略图
			err = fileUtil.GenerateThumbnail(filePath, thumbnailName)
			if err != nil {
				newFile.Status = consts.TransferFailed.Index()
				utils.LogrusObj.Error("创建缩略图失败:", err)
			}

			// 视屏切割
			if fileUtil.GetFileType(fileName) == consts.Video.Index() {
				err = fileUtil.SplitVideo(filePath)
				if err != nil {
					newFile.Status = consts.TransferFailed.Index()
					utils.LogrusObj.Error("视屏切割出错:", err)
				}
			}

			// 更新数据库
			fileDao = dao.NewFileDao(ctx)
			err = fileDao.UpdateFile(newFile)
			if err != nil {
				utils.LogrusObj.Error("更新文件状态出错:", err)
				return
			}
		}(context.Background())
	}

	// 设置响应体
	data := &types.UploadFileResp{
		FileId: newFile.FileId,
		Status: consts.Uploaded.String(),
	}

	// 响应数据
	return ctl.RespSuccessWithData(data), err
}
