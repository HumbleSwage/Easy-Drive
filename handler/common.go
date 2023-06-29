package handler

import (
	"easy-drive/pkg/ctl"
	"easy-drive/pkg/e"
	"easy-drive/types"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorResponse(err error) *types.Response {
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return ctl.RespError(e.ErrorJsonType)
	}
	return ctl.RespError(e.InvalidParams)
}

func ValidateParams(ctx *gin.Context, req interface{}) bool {
	validate := types.GetValidate()
	err := validate.Struct(req)
	if err != nil {
		// 参数校验失败
		status := e.ParameterValidationError
		r := types.Response{
			Code:  status,
			Info:  e.GetMsg(status),
			Error: err.Error(),
		}
		ctx.JSON(http.StatusOK, r)
		return false
	}
	return true
}
