package types

import "mime/multipart"

type LoadDataListReq struct {
	Category      string `json:"category" form:"category"`
	FilePid       string `json:"filePid" form:"filePid"`
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

type NewFolderReq struct {
	FilePid  string `json:"filePid" form:"filePid" validate:"required"`   // 文件父id
	FileName string `json:"fileName" form:"fileName" validate:"required"` // 目录名
}

type GetFolderInfoReq struct {
	Path    string `json:"path" form:"path" validate:"required"`
	ShareId string `json:"shareId" form:"shareId"`
}

type RenameReq struct {
	FileId   string `json:"fileId" form:"fileId" validate:"required"`
	FileName string `json:"fileName" form:"fileName" validate:"required"`
}

type LoadAllFolderReq struct {
	FilePid        string `json:"filePid" form:"filePid" validate:"required"`
	CurrentFileIds string `json:"currentFileIds" form:"currentFileIds"`
}

type ChangeFileFolderReq struct {
	FileIds string `json:"fileIds" form:"fileIds" validate:"required"`
	FilePid string `json:"filePid" form:"filePid" validate:"required"`
}

type DelFileReq struct {
	FileIds string `json:"fileIds" form:"fileIds" validate:"required"`
}
