package handler

import (
	"easy-drive/service"
	"easy-drive/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadRecycleListHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.LoadReRecycleListReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 调用服务
		l := service.GetRecycleSrv()
		resp, err := l.LoadRecycleListService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 响应成功
		ctx.JSON(http.StatusOK, resp)

	}
}

func RecoverFileHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.RecoverFileReq

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
		l := service.GetRecycleSrv()
		resp, err := l.RecoverFileService(ctx, &req)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// 正常响应
		ctx.JSON(http.StatusOK, resp)
	}

}
