package router

import (
	"easy-drive/handler"
	"easy-drive/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.Cors())
	ginRouter.Use(middleware.Session("something-very-secret"))
	api := ginRouter.Group("api")
	{
		// 测试
		api.GET("ping", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, "success")

		})

		// 验证码
		api.GET("checkCode", handler.CheckCode())
		api.POST("sendEmailCode", handler.SendEmailCode())

		// 用户
		v1 := api.Group("/")
		{
			v1.POST("register", handler.UserRegisterHandler())
			v1.POST("login", handler.UserLoginHandler())
			v1.POST("resetPwd", handler.UserResetPasswordHandler())
			v1.GET("getAvatar/:userId", handler.GetUserAvatarHandler())
			v1.POST("getUseSpace", handler.GetUseSpaceHandler())
			v1.POST("logout", handler.LogoutHandler())
			v1.POST("updateUserAvatar", handler.UpdateUserAvatarHandler())
			v1.POST("updatePassword", handler.UpdatePasswordHandler())
		}

		// 文件
		v2 := api.Group("file")
		{
			v2.POST("loadDataList", handler.LoadDataListHandler())
			v2.POST("uploadFile", handler.UploadFileHandler())
			v2.GET("getImage/*filePath", handler.GetCoverHandler())
			v2.GET("ts/getVideoInfo/:fileId", handler.GetVideoInfoHandler())
			v2.POST("getFile/:fileId", handler.GetFileHandler())
			v2.GET("getFile/:fileId", handler.GetFileHandler())
			v2.POST("newFoloder", handler.NewFolderHandler())
			v2.POST("getFolderInfo", handler.GetFolderInfoHandler())
			v2.POST("rename", handler.RenameHandler())
			v2.POST("loadAllFolder", handler.LoadAllFolderHandler())
			v2.POST("changeFileFolder", handler.ChangeFileFolderHandler())
			v2.POST("createDownloadUrl/:fileId", handler.CreateDownloadUrlHandler())
			v2.GET("download/:downloadCode", handler.DownloadFileHandler())
			v2.POST("delFile", handler.DelFileHandler())
		}

		// 回收站
		v3 := api.Group("recycle")
		{
			v3.POST("loadRecycleList", handler.LoadRecycleListHandler())
			v3.POST("recoverFile", handler.RecoverFileHandler())
		}

		// 分享
		v4 := api.Group("share")
		{
			v4.POST("loadShareList", handler.LoadShareListHandler())
			v4.POST("shareFile", handler.ShareFileHandler())
			v4.POST("cancelShare", handler.CancelShareHandler())
		}

		// 超级管理员
		v5 := api.Group("admin")
		{
			v5.POST("getSysSettings", handler.GetSysSettingHandler())
			v5.POST("saveSysSettings", handler.SaveSysSettingHandler())
			v5.POST("loadUserList", handler.LoadUserListHandler())
			v5.POST("updateUserStatus", handler.UpdateUserStatusHandler())
			v5.POST("updateUserSpace", handler.UpdateUserSpaceHandler())
			v5.POST("loadFileList", handler.LoadFileListHandler())
			v5.POST("getFolderInfo", handler.AdminGetFolderHandler())
			v5.POST("getFile/:userId/:fileId", handler.AdminGetFileHandler())
			v5.GET("ts/getVideoInfo/:userId/:fileId", handler.AdminGetVideoHandler())
			v5.POST("createDownloadUrl/:userId/:fileId", handler.AdminDownloadUrlHandler())
			v5.GET("download/:downloadCode", handler.AdminDownloadHandler())
			v5.POST("delFile", handler.AdminDelFileHandler())
		}

		// 获取外部分享

		v6 := api.Group("showShare")
		{
			v6.POST("getShareLoginInfo", handler.ShareLoginInfoHandler())
			v6.POST("getShareInfo", handler.ShareInfoHandler())
			v6.POST("checkShareCode", handler.CheckShareCodeHandler())
			v6.POST("loadFileList", handler.LoadShareHandler())
			v6.POST("getFolderInfo", handler.GetShareFolderInfoHandler())
			v6.GET("getFile/:shareId/:fileId", handler.GetShareFileHandler())
			v6.GET("ts/getVideoInfo/:shareId/:fileId", handler.GetShareVideoHandler())
			v6.POST("createDownloadUrl/:shareId/:fileId", handler.CreateShareDownloadUrlHandler())
			v6.GET("download/:downloadCode", handler.ShareDownloadHandler())
			v6.POST("saveShare", handler.SaveShareHandler())
		}
	}

	return ginRouter
}
