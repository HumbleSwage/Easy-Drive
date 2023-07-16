package ctl

import (
	"easy-drive/pkg/e"
	"easy-drive/types"
)

func RespSuccess(codes ...int) *types.Response {
	code := e.Success
	if codes != nil {
		code = codes[0]
	}
	return &types.Response{
		Code:   code,
		Status: "请求成功",
		Info:   e.GetMsg(code),
	}
}

func RespSuccessWithData(data interface{}, codes ...int) *types.Response {
	code := e.Success
	if codes != nil {
		code = codes[0]
	}
	return &types.Response{
		Code:   code,
		Data:   data,
		Status: "请求成功",
		Info:   e.GetMsg(code),
	}
}

func RespError(codes ...int) *types.Response {
	code := e.Error
	if codes != nil {
		code = codes[0]
	}

	return &types.Response{
		Code:   code,
		Status: "请求失败",
		Info:   e.GetMsg(code),
	}
}
