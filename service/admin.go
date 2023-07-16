package service

import (
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

var AdminSrv *AdminService
var AdminSrvOnce sync.Once

type AdminService struct {
}

func GetAdminSrv() *AdminService {
	AdminSrvOnce.Do(func() {
		AdminSrv = &AdminService{}
	})
	return AdminSrv
}

func (as *AdminService) GetSysSettingService(ctx *gin.Context, req *types.SystemSettingReq) (resp interface{}, err error) {
	code := e.Success

	systemDao := dao.NewSystemDao(ctx.Request.Context())

	// 响应数据
	system := &types.SystemSettingResp{}
	title, err := systemDao.GetSystemSettingById(consts.UserRegisterTiTleSettingId)
	content, err := systemDao.GetSystemSettingById(consts.UserRegisterContentSettingId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("获取用户注册页面模版出错：", err)
			return ctl.RespError(), nil
		}
	}

	// 绑定注册内容
	system.RegisterEmailTitle = title.Text
	system.RegisterEmailContent = content.Text

	// 获取用户初始内容
	space, err := systemDao.GetSystemSettingById(consts.UserInitSpaceId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("获取用户初始容量出错：", err)
			return ctl.RespError(), nil
		}
	}
	system.UserInitUseSpace = space.Text

	return ctl.RespSuccessWithData(system, code), nil
}

func (as *AdminService) SaveSysSettingService(ctx *gin.Context, req *types.SystemSettingReq) (resp interface{}, err error) {
	systemDao := dao.NewSystemDao(ctx.Request.Context())

	// 更新注册邮件title
	setting, err := systemDao.GetSystemSettingById(consts.UserRegisterTiTleSettingId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("查询注册邮件title出错：", err)
			return ctl.RespError(), nil
		}
	}
	setting.Text = req.RegisterEmailTitle
	systemDao.Updates(&setting)

	// 更新注册邮件content
	setting, err = systemDao.GetSystemSettingById(consts.UserRegisterContentSettingId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("查询注册邮件content出错：", err)
			return ctl.RespError(), nil
		}
	}
	setting.Text = req.RegisterEmailContent
	systemDao.Updates(&setting)

	// 更新初始内存用量
	setting, err = systemDao.GetSystemSettingById(consts.UserRegisterContentSettingId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("查询注册邮件content出错：", err)
			return ctl.RespError(), nil
		}
	}
	setting.Text = req.UserInitUseSpace
	systemDao.Updates(&setting)

	return ctl.RespSuccess(), nil
}

func (as *AdminService) LoadUserListService(ctx *gin.Context, req *types.LoadUserListReq) (resp interface{}, err error) {
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

	// 处理请求status
	var status []string
	if req.Status == "" {
		status = []string{"0", "1"}
	} else {
		status = []string{req.Status}
	}

	// 查询用户
	var users []*model.User
	var count int64
	userDao := dao.NewUserDao(ctx.Request.Context())
	if req.NickNameFuzzy == "" {
		// 精度查找模式
		users, count, err = userDao.SelectUser(status, pageNum, pageSize)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				utils.LogrusObj.Error("管理端查询所有用户时出错：", err)
				return ctl.RespError(), nil
			}
		}
	} else {
		// 模糊搜寻模式
		users, count, err = userDao.SelectUserByFuzzy(status, req.NickNameFuzzy, pageNum, pageSize)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				utils.LogrusObj.Error("管理端模糊搜寻用户时出错：", err)
				return ctl.RespError(), nil
			}
		}
	}

	d := make([]*types.LoadUserInfo, 0)
	// 遍历绑定user
	for _, user := range users {
		if user.UserId == userId {
			// 排除管理员自身
			continue
		}
		l := &types.LoadUserInfo{
			UserId:        user.UserId,
			NickName:      user.NickName,
			Email:         user.Email,
			JoinTime:      user.CreatedAt.Format("2006-01-02 15:04:05"),
			LastLoginTime: user.LastLoginTime.Format("2006-01-02 15:04:05"),
			Status:        cast.ToInt(user.Status),
			UseSpace:      user.UseSpace,
			TotalSpace:    user.TotalSpace,
		}
		d = append(d, l)
	}

	// 设置响应体
	data := &types.ListInfoResp{
		TotalCount: count,
		PageSize:   pageSize,
		PageNo:     pageNum,
		PageTotal:  len(d),
		List:       d,
	}
	return ctl.RespSuccessWithData(data, code), nil
}

func (as *AdminService) UpdateUserStatusService(ctx *gin.Context, req *types.UpdateUserStatusReq) (resp interface{}, err error) {
	code := e.Success

	// 获取用户
	userDao := dao.NewUserDao(ctx.Request.Context())
	user, err := userDao.GetUserByUserId(req.UserId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("根据用户id查询用户出错：", err)
			return ctl.RespError(), nil
		}
	}

	// 更新用户状态
	user.Status = cast.ToInt8(req.Status)

	// 数据入库
	err = userDao.UpdateUser(user)
	if err != nil {
		utils.LogrusObj.Error("更新用户状态时出错：", err)
		return ctl.RespError(), nil
	}

	return ctl.RespSuccess(code), nil
}

func (as *AdminService) UpdateUserSpaceService(ctx *gin.Context, req *types.UpdateUserSpaceReq) (resp interface{}, err error) {
	code := e.Success

	// 获取用户
	userDao := dao.NewUserDao(ctx.Request.Context())
	user, err := userDao.GetUserByUserId(req.UserId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("根据用户id查询用户出错：", err)
			return ctl.RespError(), nil
		}
	}

	// 判断分配空间大小
	if cast.ToInt64(req.ChangeSpace) > consts.UserLimitedSpace {
		code = e.OverLimitUserSpaceError
		return ctl.RespError(code), nil
	}

	// 更新用户空间
	user.TotalSpace = cast.ToInt64(req.ChangeSpace) * 1024 * 1024

	// 数据入库
	err = userDao.UpdateUser(user)
	if err != nil {
		utils.LogrusObj.Error("更新用户状态时出错：", err)
		return ctl.RespError(), nil
	}

	return ctl.RespSuccess(code), nil
}

func (as *AdminService) LoadFileListService(ctx *gin.Context, req *types.LoadFileListReq) (resp interface{}, err error) {
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

	// 确定搜索模式
	var files []*model.File
	var count int64
	fileDao := dao.NewFileDao(ctx.Request.Context())

	if req.FileNameFuzzy == "" {
		// 搜索全部文件
		files, count, err = fileDao.SelectFile(req.FilePid, pageNum, pageSize)
	} else {
		// 模糊搜索文件
		files, count, err = fileDao.SelectFileByFuzzyName(req.FilePid, req.FileNameFuzzy, pageNum, pageSize)
	}

	// 绑定数据
	userDao := dao.NewUserDaoByDB(fileDao.DB)
	fileArray := make([]*types.LoadFileInfo, 0)
	for _, file := range files {
		user, err := userDao.GetUserByUserId(file.UserId)
		if err != nil {
			utils.LogrusObj.Error("用户id查询用户出错：", err)
			return ctl.RespError(), nil
		}
		l := &types.LoadFileInfo{
			FileId:         file.FileId,
			FilePid:        file.ParentId,
			UserId:         file.UserId,
			UserName:       user.UserName,
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

		fileArray = append(fileArray, l)
	}

	// 返回数据
	data := &types.ListInfoResp{
		TotalCount: count,
		PageSize:   pageSize,
		PageNo:     pageNum,
		PageTotal:  len(fileArray),
		List:       fileArray,
	}

	return ctl.RespSuccessWithData(data, code), nil
}

func (as *AdminService) GetFolderInfoService(ctx *gin.Context, req *types.AdminGetFolderReq) (resp interface{}, err error) {
	code := e.Success

	// 切割路径
	folderIdArray := strings.Split(req.Path, "/")

	// 查询文件
	fileDao := dao.NewFileDao(ctx.Request.Context())
	folders, err := fileDao.SelectFolderInfo(folderIdArray)
	if err != nil {
		utils.LogrusObj.Error("管理端查询用户文件时出错：", err)
		return
	}

	// 绑定数据
	data := make([]*types.AdminGetFolderInfoResp, 0)
	for _, folder := range folders {
		r := &types.AdminGetFolderInfoResp{
			FileName: folder.Name,
			FileId:   folder.FileId,
		}
		data = append(data, r)
	}

	return ctl.RespSuccessWithData(data, code), nil
}

func (as *AdminService) GetFileService(ctx *gin.Context, userId, fileId string) (resp interface{}, err error) {
	// 查询文件
	fileDao := dao.NewFileDao(ctx.Request.Context())
	file, err := fileDao.GetFileById(fileId, userId)
	if err != nil {
		utils.LogrusObj.Error("查询文件出错：", err)
		return ctl.RespError(), nil
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

	return data, err
}

func (as *AdminService) GetVideoService(ctx *gin.Context, userId, fileId string) (resp interface{}, err error) {
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

func (as *AdminService) CreateDownloadUrlService(ctx *gin.Context, userId, fileId string) (resp interface{}, err error) {
	code := e.Success

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

func (as *AdminService) DownloadFileService(ctx *gin.Context, downloadCode string) (resp interface{}, err error) {
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

func (as *AdminService) DelFileService(ctx *gin.Context, req *types.AdminDelFileReq) (resp interface{}, err error) {
	code := e.Success

	// 获取fileIdArray
	fileIdAndUserIds := strings.Split(req.FileIdAndUserIds, ",")

	// 遍历fileIdArray
	fileDao := dao.NewFileDao(ctx.Request.Context())
	for _, udAndFd := range fileIdAndUserIds {
		userId, fileId := strings.Split(udAndFd, "_")[0], strings.Split(udAndFd, "_")[1]
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
