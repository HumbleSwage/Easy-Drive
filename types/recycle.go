package types

type LoadReRecycleListReq struct {
	PageNo   string `json:"pageNo" form:"pageNo"`
	PageSize string `json:"pageSize" form:"pageSize"`
}

type RecoverFileReq struct {
	FileIds string `json:"fileIds" form:"fileIds" validate:"required"`
}
