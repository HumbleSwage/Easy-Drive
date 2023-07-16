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

func LoadShareListHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.LoadShareListReq

		// 绑定参数
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 调用服务
		l := service.GetShareSrv()
		resp, err := l.LoadShareListService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 成功响应
		ctx.JSON(http.StatusOK, resp)
	}
}

func ShareFileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ShareFileReq

		// 绑定参数
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}

		// 校验参数
		if !ValidateParams(ctx, req) {
			return
		}

		// 调用服务
		l := service.GetShareSrv()
		resp, err := l.ShareFileService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 正常响应
		ctx.JSON(http.StatusOK, resp)

	}
}

func CancelShareHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.CancelShareReq

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
		l := service.GetShareSrv()
		resp, err := l.CancelShareService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应成功
		ctx.JSON(http.StatusOK, resp)
	}
}

func ShareLoginInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ShowShareReq

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
		l := service.GetShareSrv()
		resp, err := l.ShareLoginInfoService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应
		ctx.JSON(http.StatusOK, resp)

	}
}

func ShareInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ShowShareReq

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
		l := service.GetShareSrv()
		resp, err := l.ShareInfoService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应
		ctx.JSON(http.StatusOK, resp)

	}
}

func CheckShareCodeHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.CheckShareReq

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
		l := service.GetShareSrv()
		resp, err := l.CheckShareCodeService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 成功响应
		ctx.JSON(http.StatusOK, resp)

	}
}

func LoadShareHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.LoadShareReq

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
		l := service.GetShareSrv()
		resp, err := l.LoadShareService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 成功响应
		ctx.JSON(http.StatusOK, resp)

	}
}

func GetShareFolderInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ShareFolderInfoReq

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
		l := service.GetShareSrv()
		resp, err := l.GetShareFolderInfoService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 成功响应
		ctx.JSON(http.StatusOK, resp)
	}
}

func GetShareFileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 解析url参数
		shareId := ctx.Param("shareId")
		fileId := ctx.Param("fileId")

		// 校验参数
		if shareId == "" || fileId == "" {
			ctx.JSON(http.StatusBadRequest, errors.New("请求参数错误"))
			return
		}

		// 调用服务
		l := service.GetShareSrv()
		resp, err := l.GetShareFileService(ctx, shareId, fileId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应数据
		_, err = ctx.Writer.Write(resp.([]byte))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
	}
}

func GetShareVideoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 解析url参数
		shareId := ctx.Param("shareId")
		fileId := ctx.Param("fileId")

		// 校验参数
		if shareId == "" || fileId == "" {
			ctx.JSON(http.StatusBadRequest, errors.New("请求参数错误"))
			return
		}

		// 调用服务
		l := service.GetShareSrv()
		resp, err := l.GetShareVideoService(ctx, shareId, fileId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 断言对象
		if resp == nil {
			utils.LogrusObj.Error("数据响应为空：", err)
			ctx.JSON(http.StatusInternalServerError, ctl.RespError())
		}
		data, ok := resp.([]byte)
		if !ok {
			utils.LogrusObj.Error("数据断言出错")
			ctx.JSON(http.StatusInternalServerError, ctl.RespError())
			return
		}

		// 响应数据
		_, err = ctx.Writer.Write(data)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应信息
		ctx.JSON(http.StatusOK, ctl.RespSuccess())
	}
}

func CreateShareDownloadUrlHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 解析url参数
		shareId := ctx.Param("shareId")
		fileId := ctx.Param("fileId")

		// 参数校验
		if shareId == "" || fileId == "" {
			ctx.JSON(http.StatusBadRequest, errors.New("请求参数不足"))
			return
		}

		// 调用服务
		l := service.GetShareSrv()
		resp, err := l.CreateShareDownloadUrl(ctx, shareId, fileId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应正常信息
		ctx.JSON(http.StatusOK, resp)
	}
}

func ShareDownloadHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求路径中的参数
		shareCode := ctx.Param("downloadCode")

		// 校验参数
		if shareCode == "" {
			ctx.JSON(http.StatusBadRequest, errors.New("缺少请求参数"))
			return
		}

		// 调用服务
		l := service.GetShareSrv()
		resp, err := l.ShareDownloadService(ctx, shareCode)
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

func SaveShareHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.SaveShareReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 调用服务
		l := service.GetShareSrv()
		resp, err := l.SaveShareService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应成功
		ctx.JSON(http.StatusOK, resp)

	}
}
