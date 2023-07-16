package service

import (
	"context"
	"easy-drive/conf"
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
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
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

	// 处理默认值
	if req.FilePid == "" {
		req.FilePid = "0"
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
	switch req.Category {
	case consts.VideoCategory.String():
		files, totalCount, err = fileDao.GetUserFilesByCategory(user, req.FilePid, consts.VIDEO.Index(), pageNum, pageSize)
	case consts.MusicCategory.String():
		files, totalCount, err = fileDao.GetUserFilesByCategory(user, req.FilePid, consts.MUSIC.Index(), pageNum, pageSize)
	case consts.ImageCategory.String():
		files, totalCount, err = fileDao.GetUserFilesByCategory(user, req.FilePid, consts.IMAGE.Index(), pageNum, pageSize)
	case consts.DocCategory.String():
		files, totalCount, err = fileDao.GetUserFilesByCategory(user, req.FilePid, consts.DocCategory.Index(), pageNum, pageSize)
	case consts.OthersCategory.String():
		files, totalCount, err = fileDao.GetUserFilesByCategory(user, req.FilePid, consts.OthersCategory.Index(), pageNum, pageSize)
	default:
		// 默认查询所有
		files, totalCount, err = fileDao.GetUserFiles(user.UserId, req.FilePid, pageNum, pageSize)
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
		// 正常文件，没有在回收站中
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
		PageNo:     int64(pageNum),
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
		// 头分块
		newFile.FileId = commonUtil.GenerateFileID()

		// 根据请求文件名查询文件:TODO:有点臃肿，后面单独封装
		file, err := fileDao.GetFileByFileName(req.FileName, req.FilePid, userId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newFile.Name = req.FileName
			} else {
				utils.LogrusObj.Error("根据文件名查询文件出错:", err)
				return ctl.RespError(), nil
			}
		} else {
			fileNameNoExt := fileUtil.GetFileNameWithoutExtension(file.Name)
			fmt.Println(fileNameNoExt)

			// 解析文件名和当前索引
			currentIndex := fileUtil.GetFileNameIndex(file.Name)
			newIndex := currentIndex + 1

			// 构造新的文件名
			newFileName := fmt.Sprintf("%s(%d).txt", fileNameNoExt, newIndex)

			// 检查新文件名是否已存在：TODO：这个地方的效率有问题，反复访问数据库，建议先按照某种规则先查询出来，然后在服务端遍历，减少访问数据库的方式
			for {
				file, err = fileDao.GetFileByFileName(newFileName, req.FilePid, userId)
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
			newFile.Name = newFileName
		}
	} else {
		newFile.FileId = req.FileId
		newFile.Name = req.FileName
	}

	// 获取文件类型
	t, flag := fileUtil.FindTypeKey(req.FileName)
	if !flag {
		utils.LogrusObj.Error("获取文件类型出错:", err)
		return ctl.RespError(), nil
	}
	newFile.Type = t

	// 获取文件真实目录
	filePath := fileUtil.GetFilePath(newFile.FileId, userId, req.FileName)

	// 判断是否为妙传
	if chunkIndex == 0 {
		// 查询文件
		file, err := fileDao.SelectFileByMd5(req.FileMd5, consts.NormalFile.Index())
		if err != nil && err != gorm.ErrRecordNotFound {
			utils.LogrusObj.Error("通过文件md5值查询文件出错:", err)
			return ctl.RespError(), nil
		}

		// 文件秒传
		if file.FileId != "" {
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
	key := cache.UserStoreSpaceKey(userId, newFile.FileId)

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
		tempFilePath := fileUtil.GetTempFilePath(newFile.FileId, userId)

		// 临时存储
		err := ctx.SaveUploadedFile(req.File, tempFilePath)
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
		tempFilePath := fileUtil.GetTempFilePath(req.FileId, userId)

		// 临时存储
		err = ctx.SaveUploadedFile(req.File, tempFilePath)
		if err != nil {
			utils.LogrusObj.Error("文件分片上传到临时文件夹出错:", err)
			return ctl.RespError(), nil
		}

		// 合并分片到指定位置
		err = fileUtil.MergeChunks(userId, req.FileId, filePath)
		if err != nil {
			utils.LogrusObj.Error("合并分片出错:", err)
		} else {
			utils.LogrusObj.Infoln("文件合并完成:", filePath)
		}
	}

	// 绑定数据
	newFile.Path = filePath
	newFile.Flag = consts.NormalFile.Index()
	newFile.Category, flag = fileUtil.FindCategoryKey(req.FileName)
	if !flag {
		utils.LogrusObj.Error("设置file的category出错:", err)
	}
	if newFile.Category == consts.VIDEO.Index() || newFile.Category == consts.IMAGE.Index() {
		newFile.Status = consts.Transfer.Index() // 正在转码中
	} else {
		newFile.Status = consts.Using.Index() // 正在转码中
	}

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
	if newFile.Category == consts.VIDEO.Index() || newFile.Category == consts.IMAGE.Index() {
		go func(ctx context.Context) {
			// 更新文件状态
			newFile.Status = consts.Using.Index()

			// 缩略图名称
			thumbnailName := fileUtil.GenerateThumbnailName(filePath)
			utils.LogrusObj.Infoln("缩略图名称:", thumbnailName)

			// 生成缩略图
			err = fileUtil.GenerateThumbnail(filePath, thumbnailName)
			if err != nil {
				newFile.Status = consts.TransferFailed.Index()
				utils.LogrusObj.Error("创建缩略图失败:", err)
			}

			// 设置封面地址
			newFile.Cover = thumbnailName

			// 视屏切割
			if newFile.Category == consts.VIDEO.Index() {
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

func (fs *FileService) GetVideoInfoService(ctx *gin.Context, fileId string) (resp interface{}, err error) {
	// 获取用户
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

	// 获取ts
	extension := fileUtil.GetFileExtension(fileId)
	chunkId := ""
	if extension == "ts" {
		fileId, chunkId = strings.Split(fileId, "_")[0], strings.Split(fileId, "_")[1]
	}

	// 获取文件
	fileDao := dao.NewFileDao(ctx.Request.Context())
	file, err := fileDao.GetFileById(fileId, userId)
	if err != nil {
		utils.LogrusObj.Error("通过文件id和用户id获取文件出错:", err)
		return nil, err
	}

	// 获取m3u8地址
	mFilePath := ""
	uConfig := conf.Conf.UploadPath
	workDir, _ := os.Getwd()
	mFileDir := fileUtil.GetFileNameWithoutExtension(file.Path)
	if extension == "ts" {
		mFilePath = filepath.Join(workDir, uConfig.VideoPath, "user_"+userId, mFileDir, fileId+"_"+chunkId)
		utils.LogrusObj.Infoln("视频分块地址：", mFilePath)
	} else {
		mFilePath = filepath.Join(workDir, uConfig.VideoPath, "user_"+userId, mFileDir, "index.m3u8")
		utils.LogrusObj.Infoln("index文件地址：", mFilePath)
	}

	// 获取文件基本信息
	var fileInfo os.FileInfo
	if mFilePath != "" {
		fileInfo, err = os.Stat(mFilePath)
		if err != nil {
			utils.LogrusObj.Error("文件不存在:", err)
			return nil, err
		}
	} else {
		utils.LogrusObj.Error("文件地址为空:", err)
		return nil, err
	}

	// 打开文件
	fStream, err := os.Open(mFilePath)
	if err != nil {
		utils.LogrusObj.Error("打开文件出错:", err)
		return nil, err
	}

	// 读取数据流
	data := make([]byte, fileInfo.Size())
	if _, err := fStream.Read(data); err != nil {
		utils.LogrusObj.Error("读取文件流出错:", err)
		return nil, err
	}

	return data, nil
}

func (fs *FileService) GetFileService(ctx *gin.Context, fileId string) (resp interface{}, err error) {
	// 获取用户
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

	// 获取文件
	fileDao := dao.NewFileDao(ctx.Request.Context())
	file, err := fileDao.GetFileById(fileId, userId)
	if err != nil {
		utils.LogrusObj.Error("通过id获取文件出错:", err)
		return nil, err
	}

	// 读取地址
	workDir, _ := os.Getwd()
	filePath := filepath.Join(workDir, file.Path)

	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		utils.LogrusObj.Error("获取文件信息:", err)
		return nil, err
	}

	// 创建数据流
	data := make([]byte, fileInfo.Size())

	// 打开文件
	fStream, err := os.Open(filePath)
	if err != nil {
		utils.LogrusObj.Error("打开文件出错:", err)
		return nil, err
	}

	// 读取文件
	if _, err = fStream.Read(data); err != nil {
		utils.LogrusObj.Error("读取文件出错:", err)
		return nil, err
	}

	return data, nil
}

func (fs *FileService) NewFolderService(ctx *gin.Context, req *types.NewFolderReq) (resp interface{}, err error) {
	// 获取用户
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

	// 检查文件夹是否命名重复
	fileDao := dao.NewFileDao(ctx.Request.Context())
	count, err := fileDao.CheckFolderName(userId, req.FilePid, req.FileName)
	if err != nil {
		utils.LogrusObj.Error("检查文件夹名是否重复出错:", err)
		return ctl.RespError(), nil
	}

	// 命名
	if count != 0 {
		code = e.FileNameExistsError
		return ctl.RespSuccess(code), nil
	}

	// 创建数据
	newFile := &model.File{
		FileId:      commonUtil.GenerateFileID(),
		UserId:      userId,
		ParentId:    req.FilePid,
		Name:        req.FileName,
		Status:      consts.Using.Index(),
		IsDirectory: true,
		Flag:        consts.NormalFile.Index(),
	}

	// 数据入库
	err = fileDao.CreateFile(newFile)
	if err != nil {
		utils.LogrusObj.Error("创建文件出错:", err)
		return nil, err
	}

	return ctl.RespSuccess(), nil
}

func (fs *FileService) GetFolderInfoService(ctx *gin.Context, req *types.GetFolderInfoReq) (resp interface{}, err error) {
	// 获取用户
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

	// 切割路径
	folderIdArray := strings.Split(req.Path, "/")

	// 获取指定文件下的所有文件
	fileDao := dao.NewFileDao(ctx.Request.Context())
	files, err := fileDao.GetFolderInfo(userId, folderIdArray)
	if err != nil {
		utils.LogrusObj.Infoln("搜索文件夹下文件出错:", err)
		return ctl.RespError(), nil
	}

	// 绑定数据
	data := make([]*types.GetFolderInfoResp, 0)
	for _, file := range files {
		r := &types.GetFolderInfoResp{
			FileId:   file.FileId,
			FileName: file.Name,
		}
		data = append(data, r)
	}

	return ctl.RespSuccessWithData(data), nil
}

func (fs *FileService) RenameService(ctx *gin.Context, req *types.RenameReq) (resp interface{}, err error) {
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

	// 获取文件
	fileDao := dao.NewFileDao(ctx.Request.Context())
	file, err := fileDao.GetFileById(req.FileId, userId)

	// 检查命名
	file, err = fileDao.GetFileByFileName(req.FileName, file.ParentId, userId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("根据文件名称查询文件时出错:", err)
			return nil, err
		}
	}

	// 文件名已有
	if file.FileId != "" {
		code = e.FileNameExistsError
		return ctl.RespError(code), nil
	}

	// 查询文件
	file, err = fileDao.GetFileById(req.FileId, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.FileNotExistsError
			return ctl.RespError(code), nil
		} else {
			utils.LogrusObj.Error("根据文件id和用户id查询文件时出错:", err)
			return nil, err
		}
	}

	// 更新数据
	file.Name = req.FileName
	err = fileDao.UpdateFile(file)
	if err != nil {
		utils.LogrusObj.Error("更新文件名时出错:", err)
		return nil, err
	}

	// 返回正常响应
	return ctl.RespSuccess(), nil
}

func (fs *FileService) LoadFolderService(ctx *gin.Context, req *types.LoadAllFolderReq) (resp interface{}, err error) {
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

	// 获取所有文件夹
	fileDao := dao.NewFileDao(ctx.Request.Context())
	files, err := fileDao.GetFolder(userId, req.FilePid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("查询文件下所有文件夹出错:", err)
			return ctl.RespError(), nil
		}
	}

	// 返回数据容器
	data := make([]*types.FileInfoResp, 0)

	// 获取currentFileIds:TODO:该处的逻辑还有点混乱
	if req.CurrentFileIds != "" {
		// 切割所选的文件夹
		currentFileIds := strings.Split(req.CurrentFileIds, ",")

		// 排除所选文件夹
		flag := false
		for _, file := range files {
			// 判断是否id包含在currentFileIds
			for _, cIds := range currentFileIds {
				if file.FileId == cIds {
					flag = true
					break
				}
			}

			// 该id没有包含在currentFileIds中
			if !flag {
				d := &types.FileInfoResp{
					FileId:         file.FileId,
					FilePid:        file.ParentId,
					FileSize:       file.Size,
					FileName:       file.Name,
					FileCover:      file.Cover,
					LastUpdateTime: file.UpdatedAt.Format("2006-01-02 15:04:05"),
					FolderType:     cast.ToInt(file.IsDirectory),
					Status:         file.Status,
				}

				// 添加
				data = append(data, d)
			}
		}
	} else {
		// 绑定数据
		for _, file := range files {
			d := &types.FileInfoResp{
				FileId:         file.FileId,
				FilePid:        file.ParentId,
				FileSize:       file.Size,
				FileName:       file.Name,
				FileCover:      file.Cover,
				LastUpdateTime: file.UpdatedAt.Format("2006-01-02 15:04:05"),
				FolderType:     cast.ToInt(file.IsDirectory),
				Status:         file.Status,
			}

			// 添加
			data = append(data, d)
		}
	}

	return ctl.RespSuccessWithData(data), nil
}

func (fs *FileService) ChangeFileFolderService(ctx *gin.Context, req *types.ChangeFileFolderReq) (resp interface{}, err error) {
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

	// 分割FileIds
	fileIds := strings.Split(req.FileIds, ",")

	// 查询文件
	fileDao := dao.NewFileDao(ctx.Request.Context())
	files, err := fileDao.GetFiles(userId, fileIds)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("查询文件出错:", err)
			return ctl.RespError(), nil
		}
	}

	// 更改文件
	for _, file := range files {
		// 更新数据
		file.ParentId = req.FilePid
		err = fileDao.UpdateFile(file)
		if err != nil {
			utils.LogrusObj.Error("更新文件出错:", err)
			return ctl.RespError(), nil
		}
	}

	return ctl.RespSuccess(), nil
}

func (fs *FileService) CreateDownloadUrlService(ctx *gin.Context, fileId string) (resp interface{}, err error) {
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

	fileDao := dao.NewFileDao(ctx.Request.Context())
	// 获取文件
	files, err := fileDao.GetFiles(userId, []string{fileId})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.FileNotExistsError
			return ctl.RespError(code), nil
		}
		utils.LogrusObj.Error(err)
		return ctl.RespError(), nil
	}

	// 生成code码
	file := files[0]
	downloadCode := commonUtil.GenerateRandomNumber()

	// 获取key
	key := cache.DownloadFileKey(userId, downloadCode)

	// 获取redis客户端
	rClient := cache.RedisClient

	// 设置过期时间
	expiration := consts.DownloadExpiration * time.Minute
	fileInfo := fmt.Sprintf("%s:%s", file.Path, file.Name)
	err = rClient.Set(key, fileInfo, expiration).Err()
	if err != nil {
		utils.LogrusObj.Error(err)
		return ctl.RespError(), nil
	}

	// 响应成功信息
	return ctl.RespSuccessWithData(downloadCode, code), nil
}

func (fs *FileService) DownloadFileService(ctx *gin.Context, downloadCode string) (resp interface{}, err error) {
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

	// 获取redis
	rClient := cache.RedisClient

	// 获取key
	key := cache.DownloadFileKey(userId, downloadCode)

	// 获取文件路径
	fileInfo, err := rClient.Get(key).Result()
	if err != nil {
		utils.LogrusObj.Error("从redis获取目标文件错误：", err)
		return ctl.RespError(), nil
	}
	filePath, fileName := strings.Split(fileInfo, ":")[0], strings.Split(fileInfo, ":")[1]
	utils.LogrusObj.Infoln(filePath)
	utils.LogrusObj.Infoln(fileName)
	// 获取文件信息
	workDir, _ := os.Getwd()
	filePath = filepath.Join(workDir, filePath)
	fInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		utils.LogrusObj.Error("目标下载文件不存在：", err)
		return ctl.RespError(), nil
	}

	// 打开文件
	f, err := os.Open(filePath)
	if err != nil {
		utils.LogrusObj.Error("打开目标下载文件错误：", err)
		return ctl.RespError(), nil
	}

	// 读取文件
	data := make([]byte, fInfo.Size())
	_, err = f.Read(data)
	if err != nil {
		utils.LogrusObj.Error("读取目标下载文件错误：", err)
		return ctl.RespError(), nil
	}

	// 绑定返回数据
	d := &types.DownloadFileResp{
		FileName:     fileName,
		DownloadCode: downloadCode,
		FilePath:     filePath,
		Data:         data,
	}

	return d, nil
}

func (fs *FileService) DelFileService(ctx *gin.Context, req *types.DelFileReq) (resp interface{}, err error) {
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

	fileDao := dao.NewFileDao(ctx.Request.Context())

	// 获取fileIdArray
	fileIdArray := strings.Split(req.FileIds, ",")

	// 遍历fileIdArray
	for _, fileId := range fileIdArray {
		// 获取目标
		file, err := fileDao.GetFileById(fileId, userId)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				utils.LogrusObj.Error("查找文件出错：", err)
				return ctl.RespError(), nil
			}
		}

		// 判断是否为文件夹
		if file.IsDirectory {
			// 文件夹，获取下面的的所有文件
			files, err := fileDao.GetFileInSameFolder(userId, []string{fileId})
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					utils.LogrusObj.Error("获取文件夹下所有文件出错：", err)
					return ctl.RespError(), nil
				}
			}

			// 更新其中每一个文件的状态
			for _, f := range files {
				// 更新删除状态
				f.Flag = consts.DeletedFile.Index()
				var t = time.Now()
				f.RestoredAt = &t

				// 更新数据库
				err = fileDao.UpdateFile(f)
				if err != nil {
					utils.LogrusObj.Error("更新文件删除状态出错：", err)
					return ctl.RespError(), nil
				}
			}
		}

		// 是文件，直接更新其delFlag状态
		file.Flag = consts.RestoreFile.Index()
		var t = time.Now()
		file.RestoredAt = &t

		// 更新数据库
		err = fileDao.UpdateFile(file)
		if err != nil {
			utils.LogrusObj.Error("更新文件删除状态出错：", err)
			return ctl.RespError(), nil
		}
	}

	// 正常响应
	return ctl.RespSuccess(code), nil
}
