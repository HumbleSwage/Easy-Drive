package handler

import (
	"easy-drive/pkg/ctl"
	"easy-drive/pkg/utils"
	"easy-drive/service"
	"easy-drive/types"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetSysSettingHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.SystemSettingReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 调用服务
		l := service.GetAdminSrv()
		resp, err := l.GetSysSettingService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 正常响应数据
		ctx.JSON(http.StatusOK, resp)
	}
}

func SaveSysSettingHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.SystemSettingReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 调用服务
		l := service.GetAdminSrv()
		resp, err := l.SaveSysSettingService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应成功数据
		ctx.JSON(http.StatusOK, resp)
	}

}

func LoadUserListHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.LoadUserListReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 调用服务
		l := service.GetAdminSrv()
		resp, err := l.LoadUserListService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应数据
		ctx.JSON(http.StatusOK, resp)
	}
}

func UpdateUserStatusHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UpdateUserStatusReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 校验参数
		if !ValidateParams(ctx, req) {
			return
		}

		// 调用服务
		l := service.GetAdminSrv()
		resp, err := l.UpdateUserStatusService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, resp)

	}
}

func UpdateUserSpaceHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UpdateUserSpaceReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 校验参数
		if !ValidateParams(ctx, req) {
			return
		}

		// 调用服务
		l := service.GetAdminSrv()
		resp, err := l.UpdateUserSpaceService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 正常响应
		ctx.JSON(http.StatusOK, resp)
	}
}

func LoadFileListHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.LoadFileListReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 调用服务
		l := service.GetAdminSrv()
		resp, err := l.LoadFileListService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 返回数据
		ctx.JSON(http.StatusOK, resp)
	}
}

func AdminGetFolderHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.AdminGetFolderReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 校验参数
		if !ValidateParams(ctx, req) {
			return
		}

		// 调用服务
		l := service.GetAdminSrv()
		resp, err := l.GetFolderInfoService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 正常响应
		ctx.JSON(http.StatusOK, resp)
	}
}

func AdminGetFileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("userId")
		fileId := ctx.Param("fileId")

		// 调用服务
		l := service.GetAdminSrv()
		resp, err := l.GetFileService(ctx, userId, fileId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		// 响应数据
		_, err = ctx.Writer.Write(resp.([]byte))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, ctl.RespSuccess())
	}
}

func AdminGetVideoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("userId")
		fileId := ctx.Param("fileId")

		// 调用服务
		l := service.GetAdminSrv()
		resp, err := l.GetVideoService(ctx, userId, fileId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		// 响应数据
		_, err = ctx.Writer.Write(resp.([]byte))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应信息
		ctx.JSON(http.StatusOK, ctl.RespSuccess())
	}
}

func AdminDownloadUrlHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("userId")
		fileId := ctx.Param("fileId")

		// 校验参数
		if fileId == "" {
			ctx.JSON(http.StatusBadRequest, errors.New("缺少请求参数"))
			return
		}

		// 请求服务
		l := service.GetAdminSrv()
		resp, err := l.CreateDownloadUrlService(ctx, userId, fileId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应正常信息
		ctx.JSON(http.StatusOK, resp)
	}
}

func AdminDownloadHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求参数
		downloadCode := ctx.Param("downloadCode")

		// 校验参数
		if downloadCode == "" {
			ctx.JSON(http.StatusBadRequest, errors.New("请求参数为空"))
			return
		}

		// 调用服务
		l := service.GetAdminSrv()
		resp, err := l.DownloadFileService(ctx, downloadCode)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 获取fileName
		downloadFileResp, ok := resp.(*types.DownloadFileResp)
		if !ok {
			utils.LogrusObj.Error("断言出错")
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 设置返回头
		ctx.Writer.Header().Set("Content-Disposition", "application/x-msdownload; charset=UTF-8")
		value := fmt.Sprintf("%s\"%s\"", "attachment;filename=", downloadFileResp.FileName)
		ctx.Writer.Header().Set("Content-Disposition", value)
		_, err = ctx.Writer.Write(downloadFileResp.Data)
		if err != nil {
			utils.LogrusObj.Error("写文件流出错：", err)
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 正常响应
		ctx.JSON(http.StatusOK, ctl.RespSuccess())
	}
}

func AdminDelFileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.AdminDelFileReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 校验数据
		if !ValidateParams(ctx, req) {
			return
		}

		// 调用服务
		l := service.GetAdminSrv()
		resp, err := l.DelFileService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 正常响应
		ctx.JSON(http.StatusOK, resp)
	}
}
