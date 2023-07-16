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
	"os"
	"path/filepath"
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

func GetCoverHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取参数值
		coverPath := ctx.Param("filePath")
		if coverPath == "" {
			ctx.JSON(http.StatusBadRequest, ctl.RespError())
			return
		}

		// 文件绝对路径
		workDir, _ := os.Getwd()
		coverPath = filepath.Join(workDir, coverPath)

		// 读取文件
		fileInfo, err := os.Stat(coverPath)
		if err != nil {
			utils.LogrusObj.Error("获取文件封面状态出错:", err)
			ctx.JSON(http.StatusInternalServerError, ctl.RespError())
			return
		}

		// 创建缓冲流
		data := make([]byte, fileInfo.Size())

		// 打开文件
		file, err := os.Open(coverPath)
		if err != nil {
			utils.LogrusObj.Error("打开文件封面出错:", err)
			ctx.JSON(http.StatusInternalServerError, ctl.RespError())
			return
		}

		// 读取文件
		if _, err = file.Read(data); err != nil {
			utils.LogrusObj.Error("读取文件封面出错:", err)
			ctx.JSON(http.StatusInternalServerError, ctl.RespError())
			return
		}

		// 设置响应头
		ctx.Writer.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
		ctx.Writer.Header().Set("Pragma", "no-cache")
		ctx.Writer.Header().Set("Expires", "0")

		// 响应
		_, err = ctx.Writer.Write(data)
		if err != nil {
			utils.LogrusObj.Error("响应文件封面出错:", err)
			ctx.JSON(http.StatusInternalServerError, ctl.RespError())
			return
		}

		ctx.JSON(http.StatusOK, ctl.RespSuccess())
	}
}

func GetVideoInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求参数
		fileId := ctx.Param("fileId")

		// 检验参数
		if fileId == "" {
			ctx.JSON(http.StatusBadRequest, ctl.RespError())
			return
		}

		// 调用服务
		l := service.GetFileSrv()
		resp, err := l.GetVideoInfoService(ctx, fileId)
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

		// 响应信息
		ctx.JSON(http.StatusOK, ctl.RespSuccess())
	}
}

func GetFileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求参数
		fileId := ctx.Param("fileId")

		// 检验参数
		if fileId == "" {
			ctx.JSON(http.StatusBadRequest, ctl.RespError())
			return
		}

		// 调用服务
		l := service.GetFileSrv()
		resp, err := l.GetFileService(ctx, fileId)
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

func NewFolderHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 数据模版
		var req *types.NewFolderReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 调用服务
		l := service.GetFileSrv()
		resp, err := l.NewFolderService(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 正常响应
		ctx.JSON(http.StatusOK, resp)
	}
}

func GetFolderInfoHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 数据模版
		var req types.GetFolderInfoReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ctl.RespError())
			return
		}

		// 数据校验
		if !ValidateParams(ctx, req) {
			return
		}

		// 调用服务
		l := service.GetFileSrv()
		resp, err := l.GetFolderInfoService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 正常响应
		ctx.JSON(http.StatusOK, resp)
	}
}

func RenameHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.RenameReq

		// 绑定参数
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		}

		// 参数校验
		if !ValidateParams(ctx, req) {
			return
		}

		// 调用服务
		l := service.GetFileSrv()
		resp, err := l.RenameService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ctl.RespError())
			return
		}

		// 成功响应
		ctx.JSON(http.StatusOK, resp)

	}
}

func LoadAllFolderHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.LoadAllFolderReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 参数校验
		if !ValidateParams(ctx, req) {
			return
		}

		// 调用服务
		l := service.GetFileSrv()
		resp, err := l.LoadFolderService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应数据
		ctx.JSON(http.StatusOK, resp)
	}
}

func ChangeFileFolderHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ChangeFileFolderReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 参数校验
		if !ValidateParams(ctx, req) {
			return
		}

		// 调用服务
		l := service.GetFileSrv()
		resp, err := l.ChangeFileFolderService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 成功响应
		ctx.JSON(http.StatusOK, resp)

	}
}

func CreateDownloadUrlHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求fileId
		fileId := ctx.Param("fileId")

		// 校验参数
		if fileId == "" {
			ctx.JSON(http.StatusBadRequest, errors.New("缺少请求参数"))
			return
		}
		// 请求服务
		l := service.GetFileSrv()
		resp, err := l.CreateDownloadUrlService(ctx, fileId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应正常信息
		ctx.JSON(http.StatusOK, resp)
	}
}

func DownloadFileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取请求参数
		downloadCode := ctx.Param("downloadCode")

		// 校验参数
		if downloadCode == "" {
			ctx.JSON(http.StatusBadRequest, errors.New("请求参数为空"))
			return
		}

		// 调用服务
		l := service.GetFileSrv()
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

func DelFileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.DelFileReq

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
		l := service.GetFileSrv()
		resp, err := l.DelFileService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 正常响应
		ctx.JSON(http.StatusOK, resp)
	}
}
