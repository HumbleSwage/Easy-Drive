package types

type LoadDataListResp struct {
	TotalCount int64       `json:"totalCount" form:"totalCount"`
	PageSize   int64       `json:"pageSize" form:"pageSize"`
	PageTotal  int64       `json:"pageTotal" form:"pageTotal"`
	List       interface{} `json:"list" form:"list"`
}

type FileInfoResp struct {
	FileId         string `json:"fileId" form:"fileId"`
	FilePid        string `json:"filePid" form:"FilePid"`
	FileSize       int64  `json:"fileSize" form:"fileSize"`
	FileName       string `json:"fileName" form:"fileName"`
	FileCover      string `json:"fileCover" form:"fileCover"`
	CreateTime     string `json:"createTime" form:"createTime"`
	LastUpdateTime string `json:"lastUpdateTime" form:"lastUpdateTime"`
	FolderType     int    `json:"folderType" form:"folderType"`
	FileCategory   int    `json:"fileCategory" form:"fileCategory"`
	FileType       int    `json:"fileType" form:"fileType"`
	Status         int    `json:"status" form:"status"`
}

type UploadFileResp struct {
	FileId string `json:"fileId" form:"fileId"`
	Status string `json:"status" form:"status"`
}
