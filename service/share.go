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

var ShareSrv *ShareService
var ShareSrvOnce sync.Once

type ShareService struct {
}

func GetShareSrv() *ShareService {
	ShareSrvOnce.Do(func() {
		ShareSrv = &ShareService{}
	})
	return ShareSrv
}

func (ss *ShareService) LoadShareListService(ctx *gin.Context, req *types.LoadShareListReq) (resp interface{}, err error) {
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

	shareDao := dao.NewShareDao(ctx)
	// 获取用户所有分享的id
	shares, count, err := shareDao.GetShareFileByUserId(userId, pageNum, pageSize)
	if err != nil {
		utils.LogrusObj.Error("获取用户分享文件出错：", err)
		return ctl.RespError(), nil
	}

	// 遍历分享文件
	shareArray := make([]*types.ShareInfoResp, 0)
	fileDao := dao.NewFileDaoByDB(shareDao.DB)
	for _, s := range shares {
		// 查询文件信息
		file, err := fileDao.GetFileById(s.FileId, s.UserId)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				utils.LogrusObj.Error("查询文件出错：", err)
				return ctl.RespError(), nil
			}
		}
		l := &types.ShareInfoResp{
			ShareId:      s.ShareId,
			FileId:       s.FileId,
			UserId:       s.UserId,
			ValidType:    s.ValidType,
			ExpireTime:   s.ExpireTime.Format("2006-01-02 15:04:05"),
			ShareTime:    s.ShareTime.Format("2006-01-02 15:04:05"),
			Code:         s.Code,
			HitCount:     s.HitCount,
			FileName:     file.Name,
			FolderType:   cast.ToInt(file.IsDirectory),
			FileCategory: file.Category,
			FileType:     file.Type,
			FileCover:    file.Cover,
		}
		shareArray = append(shareArray, l)
	}

	// 返回数据
	data := &types.LoadDataListResp{
		TotalCount: count,
		PageSize:   int64(pageSize),
		PageNo:     int64(pageNum),
		PageTotal:  int64(len(shareArray)),
		List:       shareArray,
	}

	return ctl.RespSuccessWithData(data), nil
}

func (ss *ShareService) ShareFileService(ctx *gin.Context, req *types.ShareFileReq) (resp interface{}, err error) {
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

	// 生成分享id
	shareId := commonUtil.GenerateShareId()

	// 生成分享码
	shareCode := ""
	if req.Code == "" {
		shareCode = commonUtil.GenerateShareCode()
	} else {
		shareCode = req.Code
	}

	// 过期时间
	vType := cast.ToInt(req.ValidType)
	expire := ""
	if vType == consts.OneDay.Index() {
		expire = consts.OneDay.String()
	} else if vType == consts.SevenDay.Index() {
		expire = consts.SevenDay.String()
	} else if vType == consts.ThirtyDay.Index() {
		expire = consts.ThirtyDay.String()
	} else if vType == consts.NoExpire.Index() {
		expire = consts.NoExpire.String()
	}
	var expireNum int
	if expire != "" {
		expireNum = cast.ToInt(expire)
	}

	// 设置过期时间
	t := time.Now()
	expireTime := t.AddDate(0, 0, expireNum)
	expireDuration := time.Duration(expireNum) * 24 * time.Hour

	// 创建缓存
	key := cache.ShareFileKey(shareId)
	rClient := cache.RedisClient
	value := fmt.Sprintf("%s_%s_%s_%s_%s", userId, req.FileId, shareCode, t.Format("2006-01-02 15:04:05"), expireTime.Format("2006-01-02 15:04:05"))
	_, err = rClient.Set(key, value, expireDuration).Result()
	if err != nil {
		utils.LogrusObj.Error("设置分享缓存出错：", err)
		return ctl.RespError(), nil
	}

	// 创建数据
	newShare := &model.Share{
		ShareId:    shareId,
		FileId:     req.FileId,
		UserId:     userId,
		ValidType:  vType,
		ExpireTime: &expireTime,
		ShareTime:  &t,
		Code:       shareCode,
		HitCount:   0,
	}
	shareDao := dao.NewShareDao(ctx.Request.Context())
	shareDao.Create(newShare)

	// 响应数据
	data := types.ShareInfoResp{
		ShareId:    shareId,
		FileId:     req.FileId,
		UserId:     userId,
		ValidType:  vType,
		ExpireTime: expireTime.Format("2006-01-02 15:04:05"),
		ShareTime:  t.Format("2006-01-02 15:04:05"),
		Code:       shareCode,
	}

	return ctl.RespSuccessWithData(data), nil
}

func (ss *ShareService) CancelShareService(ctx *gin.Context, req *types.CancelShareReq) (resp interface{}, err error) {
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

	// 分割id
	shareIds := strings.Split(req.ShareIds, ",")

	// 删除数据
	shareDao := dao.NewShareDao(ctx.Request.Context())
	for _, id := range shareIds {
		err = shareDao.DeleteShareFile(id, userId)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				utils.LogrusObj.Error("删除共享文件出错：", err)
				return ctl.RespError(), nil
			}
		}
	}

	// 返回数据
	return ctl.RespSuccess(), nil
}

func (ss *ShareService) ShareLoginInfoService(ctx *gin.Context, req *types.ShowShareReq) (resp interface{}, err error) {
	code := e.Success

	// 获取缓存key
	key := cache.ShareFileKey(req.ShareId)
	rClient := cache.RedisClient

	// 获取缓冲value
	shareInfo, err := rClient.Get(key).Result()
	if err != nil || shareInfo == "" {
		code = e.ShareFileExpired
		return ctl.RespError(code), nil
	}

	// 拆分元素
	shareInfos := strings.Split(shareInfo, "_")
	fmt.Println(shareInfos)
	shareUserId := shareInfos[0]
	fileId := shareInfos[1]
	shareTime := shareInfos[3]
	expireTime := shareInfos[4]

	// 新建一个返回对象
	r := &types.ShareLoginResp{
		ShareTime:  shareTime,
		ExpireTime: expireTime,
		FileId:     fileId,
	}

	// session对象
	session := sessions.Default(ctx)

	// 获取用户id
	var userId string
	if u := session.Get(consts.UserInfo); u != nil {
		userId = u.(string)
		// 判断是否为当前用户
		if userId == shareUserId {
			r.CurrentUser = true
		} else {
			r.CurrentUser = false
		}
	} else {
		r.CurrentUser = false
	}

	// 查询fileName
	fileDao := dao.NewFileDao(ctx.Request.Context())
	file, err := fileDao.GetFileByFileId(fileId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("查询文件出错：", err)
			return ctl.RespError(), nil
		}
	}
	r.FileName = file.Name

	// 查询用户
	userDao := dao.NewUserDaoByDB(fileDao.DB)
	user, err := userDao.GetUserByUserId(shareUserId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("查询用户出错：", err)
			return ctl.RespError(), nil
		}
	}
	r.NickName = user.NickName
	r.UserId = user.UserId

	return ctl.RespSuccessWithData(r, code), nil
}

func (ss *ShareService) ShareInfoService(ctx *gin.Context, req *types.ShowShareReq) (resp interface{}, err error) {
	code := e.Success

	// 获取缓存key
	key := cache.ShareFileKey(req.ShareId)
	rClient := cache.RedisClient

	fmt.Println(key)
	// 获取缓冲value
	shareInfo, err := rClient.Get(key).Result()
	if err != nil || shareInfo == "" {
		code = e.ShareFileExpired
		return ctl.RespError(code), nil
	}

	// 拆分元素
	shareInfos := strings.Split(shareInfo, "_")
	userId := shareInfos[0]
	fileId := shareInfos[1]
	shareTime := shareInfos[3]
	expireTime := shareInfos[4]

	// 查询文件
	fileDao := dao.NewFileDao(ctx.Request.Context())
	file, err := fileDao.GetFileByFileId(fileId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("查询文件出错：", err)
			return ctl.RespError(), nil
		}
	}

	// 查询用户
	userDao := dao.NewUserDaoByDB(fileDao.DB)
	user, err := userDao.GetUserByUserId(file.UserId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("查询用户出错：", err)
			return ctl.RespError(), nil
		}
	}

	// 绑定数据
	data := &types.ShareLoginResp{
		ShareTime:   shareTime,
		ExpireTime:  expireTime,
		NickName:    user.NickName,
		FileName:    file.Name,
		FileId:      fileId,
		UserId:      userId,
		CurrentUser: false,
	}

	return ctl.RespSuccessWithData(data, code), nil
}

func (ss *ShareService) CheckShareCodeService(ctx *gin.Context, req *types.CheckShareReq) (resp interface{}, err error) {
	code := e.Success

	// 获取缓存key
	key := cache.ShareFileKey(req.ShareId)
	rClient := cache.RedisClient

	// 获取缓冲value
	shareInfo, err := rClient.Get(key).Result()
	if err != nil || shareInfo == "" {
		code = e.ShareFileExpired
		return ctl.RespError(code), nil
	}

	// 拆分元素
	shareInfos := strings.Split(shareInfo, "_")
	shareCode := shareInfos[2]
	fmt.Println(shareCode)

	// 判断shareCode是否相当
	if shareCode != req.ShareCode {
		code = e.ShareCodeError
		return ctl.RespError(code), nil
	}

	return ctl.RespSuccess(code), nil
}

func (ss *ShareService) LoadShareService(ctx *gin.Context, req *types.LoadShareReq) (resp interface{}, err error) {
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

	// 获取缓存key
	key := cache.ShareFileKey(req.ShareId)
	rClient := cache.RedisClient

	// 获取缓冲value
	shareInfo, err := rClient.Get(key).Result()
	if err != nil || shareInfo == "" {
		code = e.ShareFileExpired
		return ctl.RespError(code), nil
	}

	// 拆分元素
	shareInfos := strings.Split(shareInfo, "_")
	fileId := shareInfos[1]

	// 查询文件
	fileDao := dao.NewFileDao(ctx.Request.Context())

	// 绑定数据
	var data *types.ShareListResp
	array := make([]*types.LoadShareResp, 0)
	if req.FilePid == "0" {
		file, err := fileDao.GetFileByFileId(fileId)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				utils.LogrusObj.Error("查询文件错误：", err)
				return ctl.RespError(), nil
			}
		}
		d := &types.LoadShareResp{
			FileId:         file.FileId,
			FilePid:        file.ParentId,
			FileSize:       file.Size,
			FileName:       file.Name,
			FileCover:      file.Cover,
			LastUpdateTime: file.UpdatedAt.Format("2006-01-02 15:04:05"),
			FolderType:     cast.ToInt(file.IsDirectory),
			FileCategory:   file.Category,
			FileType:       file.Type,
			Status:         file.Status,
		}
		array = append(array, d)
		data = &types.ShareListResp{
			TotalCount: 1,
			PageSize:   pageSize,
			PageNo:     pageNum,
			PageTotal:  len(array),
			List:       array,
		}
	} else {
		files, count, err := fileDao.SelectFile(req.FilePid, pageNum, pageSize)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				utils.LogrusObj.Error("查询文件错误：", err)
				return ctl.RespError(), nil
			}
		}
		for _, file := range files {
			d := &types.LoadShareResp{
				FileId:         file.FileId,
				FilePid:        file.ParentId,
				FileSize:       file.Size,
				FileName:       file.Name,
				FileCover:      file.Cover,
				LastUpdateTime: file.UpdatedAt.Format("2006-01-02 15:04:05"),
				FolderType:     cast.ToInt(file.IsDirectory),
				FileCategory:   file.Category,
				FileType:       file.Type,
				Status:         file.Status,
			}
			array = append(array, d)
		}
		data = &types.ShareListResp{
			TotalCount: count,
			PageSize:   pageSize,
			PageNo:     pageNum,
			PageTotal:  len(array),
			List:       array,
		}
	}

	return ctl.RespSuccessWithData(data, code), nil
}

func (ss *ShareService) GetShareFolderInfoService(ctx *gin.Context, req *types.ShareFolderInfoReq) (resp interface{}, err error) {
	code := e.Success

	// 获取缓存key
	key := cache.ShareFileKey(req.ShareId)
	rClient := cache.RedisClient

	// 获取缓冲value
	shareInfo, err := rClient.Get(key).Result()
	if err != nil || shareInfo == "" {
		code = e.ShareFileExpired
		return ctl.RespError(code), nil
	}

	// 拆分元素
	shareInfos := strings.Split(shareInfo, "_")
	userId := shareInfos[0]
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
	data := make([]*types.ShareFolderInfoResp, 0)
	for _, file := range files {
		r := &types.ShareFolderInfoResp{
			FileId:   file.FileId,
			FileName: file.Name,
		}
		data = append(data, r)
	}
	return ctl.RespSuccessWithData(data, code), nil
}

func (ss *ShareService) GetShareFileService(ctx *gin.Context, shareId, fileId string) (resp interface{}, err error) {
	code := e.Success

	// 获取缓存key
	key := cache.ShareFileKey(shareId)
	rClient := cache.RedisClient

	// 获取缓冲value
	shareInfo, err := rClient.Get(key).Result()
	if err != nil || shareInfo == "" {
		code = e.ShareFileExpired
		utils.LogrusObj.Error(err)
		return ctl.RespError(code), nil
	}

	// 拆分元素
	shareInfos := strings.Split(shareInfo, "_")
	userId := shareInfos[0]
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

func (ss *ShareService) GetShareVideoService(ctx *gin.Context, shareId, fileId string) (resp interface{}, err error) {
	code := e.Success

	// 获取缓存key
	key := cache.ShareFileKey(shareId)
	rClient := cache.RedisClient

	// 获取缓冲value
	shareInfo, err := rClient.Get(key).Result()
	if err != nil || shareInfo == "" {
		code = e.ShareFileExpired
		return ctl.RespError(code), nil
	}

	// 拆分元素
	shareInfos := strings.Split(shareInfo, "_")
	userId := shareInfos[0]

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
		return ctl.RespError(), err
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

func (ss *ShareService) CreateShareDownloadUrl(ctx *gin.Context, shareId, fileId string) (resp interface{}, err error) {
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

	// 创建下载文件id
	downloadCode := commonUtil.GenerateRandomNumber()
	downloadKey := cache.DownloadFileKey(userId, downloadCode)

	// 获取缓存key
	shareKey := cache.ShareFileKey(shareId)
	rClient := cache.RedisClient

	// 获取缓冲value
	shareInfo, err := rClient.Get(shareKey).Result()
	if err != nil || shareInfo == "" {
		code = e.ShareFileExpired
		return ctl.RespError(code), nil
	}

	// 拆分元素
	shareInfos := strings.Split(shareInfo, "_")
	shareUserId := shareInfos[0]

	// 获取文件
	fileDao := dao.NewFileDao(ctx.Request.Context())
	files, err := fileDao.GetFiles(shareUserId, []string{fileId})
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

	// 设置过期时间
	expiration := consts.DownloadExpiration * time.Minute
	fileInfo := fmt.Sprintf("%s:%s", file.Path, file.Name)
	err = rClient.Set(downloadKey, fileInfo, expiration).Err()
	if err != nil {
		utils.LogrusObj.Error(err)
		return ctl.RespError(), nil
	}

	// 响应成功信息
	return ctl.RespSuccessWithData(downloadCode, code), nil
}

func (ss *ShareService) ShareDownloadService(ctx *gin.Context, downloadCode string) (resp interface{}, err error) {
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

	// 获取key
	key := cache.DownloadFileKey(userId, downloadCode)

	// 获取文件路径
	rClient := cache.RedisClient
	fileInfo, err := rClient.Get(key).Result()
	if err != nil {
		utils.LogrusObj.Error("从redis获取目标文件错误：", err)
		return ctl.RespError(), nil
	}
	filePath, fileName := strings.Split(fileInfo, ":")[0], strings.Split(fileInfo, ":")[1]

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

func (ss *ShareService) SaveShareService(ctx *gin.Context, req *types.SaveShareReq) (resp interface{}, err error) {
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

	key := cache.ShareFileKey(req.ShareId)
	rClient := cache.RedisClient

	// 获取缓冲value
	shareInfo, err := rClient.Get(key).Result()
	if err != nil || shareInfo == "" {
		code = e.ShareFileExpired
		return ctl.RespError(code), nil
	}

	// 拆分元素
	shareInfos := strings.Split(shareInfo, "_")
	shareUserId := shareInfos[0]

	// 判断是否为当前用户分享的文件
	if shareUserId == userId {
		code = e.UserSaveShareError
		return ctl.RespError(code), nil
	}

	// 分割Ids
	shareFileIds := strings.Split(req.ShareFileIds, ",")

	// 查找文件
	fileDao := dao.NewFileDao(ctx.Request.Context())
	files := make([]*model.File, 0)
	for _, fileId := range shareFileIds {
		file, err := fileDao.GetFileByFileId(fileId)
		if err != nil {
			utils.LogrusObj.Error("通过文件id查询文件出错：", err)
			return ctl.RespError(), nil
		}

		newFile := &model.File{
			FileId:      commonUtil.GenerateFileID(),
			UserId:      userId,
			MD5:         file.MD5,
			ParentId:    req.MyFolderId,
			Size:        file.Size,
			Cover:       file.Cover,
			Path:        file.Path,
			IsDirectory: file.IsDirectory,
			Category:    file.Category,
			Type:        file.Category,
			Status:      consts.Using.Index(),
			Flag:        consts.NormalFile.Index(),
		}

		// 重命名的问题:查询要保存到我的folder下是否已经有该名称
		checkFile, err := fileDao.GetFileByFileName(file.Name, req.MyFolderId, userId)
		filenameEx := fileUtil.GetFileExtension(checkFile.Name)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 保存的文件夹下还没有该名称的文件
				newFile.Name = file.Name
			} else {
				utils.LogrusObj.Error("根据文件名查询文件出错:", err)
				return ctl.RespError(), nil
			}
		} else {

			fileNameNoExt := fileUtil.GetFileNameWithoutExtension(checkFile.Name)
			currentIndex := fileUtil.GetFileNameIndex(checkFile.Name)
			newIndex := currentIndex + 1

			// 构造新的文件名
			var newFileName string
			if !checkFile.IsDirectory {
				// 不是文件夹
				newFileName = fmt.Sprintf("%s(%d).%s", fileNameNoExt, newIndex, filenameEx)
			} else {
				// 文件夹
				newFileName = fmt.Sprintf("%s(%d)", fileNameNoExt, newIndex)
			}

			// 检查新文件名是否已存在：TODO：这个地方的效率有问题，反复访问数据库，建议先按照某种规则先查询出来，然后在服务端遍历，减少访问数据库的方式
			for {
				checkFile, err = fileDao.GetFileByFileName(newFileName, req.MyFolderId, userId)
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
				if !checkFile.IsDirectory {
					// 不是文件夹
					newFileName = fmt.Sprintf("%s(%d).%s", fileNameNoExt, newIndex, filenameEx)
				} else {
					// 文件夹
					newFileName = fmt.Sprintf("%s(%d)", fileNameNoExt, newIndex)
				}
			}

			// 绑定不重复名称
			newFile.Name = newFileName
		}
		files = append(files, newFile)

		// 判断是否为文件
		if file.IsDirectory {
			// 该文件夹是文件夹，找到该项目下的所有文件
			err = findAllSubFiles(ctx, file.FileId, newFile.FileId, userId)
			if err != nil {
				utils.LogrusObj.Error("循环遍历分享文件所有文件时出错:", err)
				return ctl.RespError(), nil
			}
		}
	}

	utils.LogrusObj.Infoln(len(files))
	// 执行一次性的插入操作，将所有新文件插入数据库
	if len(files) > 0 {
		err := fileDao.Table("file").Create(&files).Error
		if err != nil {
			utils.LogrusObj.Error("插入数据错误:", err)
			return ctl.RespError(), nil
		}
	}
	return ctl.RespSuccess(), nil
}

func findAllSubFiles(ctx context.Context, fileId, currentPid, userId string) (err error) {
	fileDao := dao.NewFileDao(ctx)
	var files []*model.File
	err = fileDao.DB.Model(&model.File{}).Where("parent_id = ?", fileId).Find(&files).Error
	if err != nil {
		return err
	}

	for _, f := range files {
		newFile := &model.File{
			FileId:      commonUtil.GenerateFileID(),
			UserId:      userId,
			MD5:         f.MD5,
			ParentId:    currentPid,
			Size:        f.Size,
			Name:        f.Name,
			Cover:       f.Cover,
			Path:        f.Path,
			IsDirectory: f.IsDirectory,
			Category:    f.Category,
			Type:        f.Type,
			Status:      f.Status,
			Flag:        f.Flag,
			RestoredAt:  nil,
		}

		fileDao.DB.Model(&model.File{}).Create(&newFile)

		if newFile.IsDirectory {
			err = findAllSubFiles(ctx, f.FileId, newFile.FileId, userId)
			if err != nil {
				return err
			}
		}
	}

	return
}
