package handler

import (
	"easy-drive/service"
	"easy-drive/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadDataListHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 数据模版
		var req types.LoadDataListReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 参数校验
		if !ValidateParams(ctx, &req) {
			return
		}

		// 调用服务
		l := service.GetFileSrv()
		resp, err := l.LoadDataService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 正常响应
		ctx.JSON(http.StatusOK, resp)
	}
}

func UploadFileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 数据模版
		var req types.UploadFileReq

		// 解析文件
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 获取表单字段值
		req.FileId = ctx.Request.FormValue("fileId")
		req.FileName = ctx.Request.FormValue("fileName")
		req.FilePid = ctx.Request.FormValue("filePid")
		req.FileMd5 = ctx.Request.FormValue("fileMd5")
		req.ChunkIndex = ctx.Request.FormValue("chunkIndex")
		req.Chunks = ctx.Request.FormValue("chunks")
		req.File = file

		// 参数校验
		if !ValidateParams(ctx, &req) {
			return
		}

		// 调用服务
		l := service.GetFileSrv()
		resp, err := l.UploadFileService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 正常响应
		ctx.JSON(http.StatusOK, resp)
	}
}
