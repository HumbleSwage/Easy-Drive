package types

import "mime/multipart"

type LoadDataListReq struct {
	Category      string `json:"category" form:"category" validate:"required"`
	FilePid       string `json:"filePid" form:"filePid" validate:"required"`
	FileNameFuzzy string `json:"fileNameFuzzy" form:"fileNameFuzzy"`
	PageNo        string `json:"pageNo" form:"pageNo"`
	PageSize      string `json:"pageSize" form:"pageSize"`
}

type UploadFileReq struct {
	FileId     string                `json:"fileId" form:"fileId"`
	File       *multipart.FileHeader `json:"file" form:"file" validate:"required" binding:"-"`
	FileName   string                `json:"fileName" form:"fileName" validate:"required"`
	FilePid    string                `json:"filePid" form:"filePid" validate:"required"`
	FileMd5    string                `json:"fileMd5" form:"fileMd5" validate:"required"`
	ChunkIndex string                `json:"chunkIndex" form:"chunkIndex" validate:"required"`
	Chunks     string                `json:"chunks" form:"chunks" validate:"required"`
}
