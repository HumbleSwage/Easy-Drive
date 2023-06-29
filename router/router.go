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

		}
	}

	return ginRouter
}
