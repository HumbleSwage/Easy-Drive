package handler

import (
	"easy-drive/pkg/ctl"
	"easy-drive/pkg/e"
	"easy-drive/pkg/utils"
	"easy-drive/service"
	"easy-drive/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func UserRegisterHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserRegisterReq
		err := ctx.ShouldBind(&req)
		if err == nil {
			// 校验参数
			if !ValidateParams(ctx, &req) {
				return
			}

			// 定义返回
			var resp interface{}

			// 验证码校验
			if CaptchaVerify(ctx, 0, req.CheckCode) {
				// 校验成功
				l := service.GetUserSrv()
				resp, err = l.UserRegister(ctx.Request.Context(), &req)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, err)
					return
				}
			} else {
				// 校验失败
				code := e.VerificationCodeError
				resp = ctl.RespError(code)
			}

			// 响应返回
			ctx.JSON(http.StatusOK, resp)
			return
		}

		// 参数绑定失败
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
	}
}

func UserLoginHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UserLoginReq
		err := ctx.ShouldBind(&req)
		if err != nil {
			// 参数绑定失败
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		}

		// 校验参数
		if !ValidateParams(ctx, &req) {
			return
		}

		// 定义返回
		var resp interface{}

		// 校验验证码
		if CaptchaVerify(ctx, 0, req.CheckCode) {
			// 注册验证码通过
			l := service.GetUserSrv()
			resp, err = l.UserLogin(ctx, &req)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				return
			}
		} else {
			// 验证码未通过
			code := e.VerificationCodeError
			resp = ctl.RespError(code)
		}

		// 正常响应
		ctx.JSON(http.StatusOK, resp)
	}
}

func UserResetPasswordHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 数据模版
		var req types.UserResetPwdReq

		// 参数绑定
		err := ctx.ShouldBind(&req)
		if err != nil {
			// 绑定参数失败
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		}

		// 参数校验
		if !ValidateParams(ctx, &req) {
			return
		}

		// 定义返回
		var resp interface{}

		// 校验验证码
		if CaptchaVerify(ctx, 0, req.CheckCode) {
			// 校验成功
			l := service.GetUserSrv()
			resp, err = l.UserResetPwd(ctx.Request.Context(), &req)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, err)
				return
			}
		} else {
			// 验证码未通过
			code := e.VerificationCodeError
			resp = ctl.RespError(code)
		}

		// 返回正常响应
		ctx.JSON(http.StatusOK, resp)
		return

	}
}

func GetUserAvatarHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求参数
		userId := ctx.Param("userId")

		// 检验参数
		if strings.EqualFold(userId, "") {
			ctx.JSON(http.StatusBadRequest, ctl.RespError(e.InvalidParams))
			return
		}

		// 请求服务
		l := service.GetUserSrv()
		resp, err := l.GetUserAvatar(ctx, userId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		// 断言
		data, ok := resp.([]byte)
		if !ok {
			utils.LogrusObj.Error("断言数据出错")
			ctx.JSON(http.StatusInternalServerError, ctl.RespError())
			return
		}

		// 因为直接返回图片，所以需要设置返回头
		ctx.Writer.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
		ctx.Writer.Header().Set("Pragma", "no-cache")
		ctx.Writer.Header().Set("Expires", "0")
		//ctx.Writer.Header().Set("Content-Type", "application/octet-stream")

		_, err = ctx.Writer.Write(data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
	}
}

func GetUseSpaceHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		l := service.GetUserSrv()
		resp, err := l.GetUseSpace(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, resp)

	}
}

func LogoutHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		l := service.GetUserSrv()
		resp, err := l.UserLogout(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	}
}

func UpdateUserAvatarHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 解析头像文件
		file, fileHeader, err := ctx.Request.FormFile("avatar")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 绑定数据
		var req types.UpdateUserAvatarReq
		req.File = file
		req.FileSize = fileHeader.Size

		// 调用服务
		l := service.GetUserSrv()
		resp, err := l.UpdateUserAvatar(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 成功响应
		ctx.JSON(http.StatusOK, resp)
	}
}

func UpdatePasswordHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UpdatePasswordReq
		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 检验参数
		if !ValidateParams(ctx, &req) {
			return
		}

		// 定义返回
		var resp interface{}

		// 调用服务
		l := service.GetUserSrv()
		resp, err = l.UpdatePassword(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
		}

		// 响应返回
		ctx.JSON(http.StatusOK, resp)
	}
}
