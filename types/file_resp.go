package types

type LoadDataListResp struct {
	TotalCount int64       `json:"totalCount" form:"totalCount"`
	PageSize   int64       `json:"pageSize" form:"pageSize"`
	PageNo     int64       `json:"pageNo" form:"pageNo"`
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
	RecoveryTime   string `json:"recoveryTime" form:"recoveryTime"`
	FolderType     int    `json:"folderType" form:"folderType"`
	FileCategory   int    `json:"fileCategory" form:"fileCategory"`
	FileType       int    `json:"fileType" form:"fileType"`
	Status         int    `json:"status" form:"status"`
}

type UploadFileResp struct {
	FileId string `json:"fileId" form:"fileId"`
	Status string `json:"status" form:"status"`
}

type GetFolderInfoResp struct {
	FileName string `json:"fileName"`
	FileId   string `json:"fileId"`
}

type DownloadFileResp struct {
	DownloadCode string `json:"downloadCode"`
	FileName     string `json:"fileName"`
	FilePath     string `json:"filePath"`
	Data         []byte `json:"data"`
}
