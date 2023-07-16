package types

type LoadShareListReq struct {
	PageNo   string `json:"pageNo" form:"pageNo"`
	PageSize string `json:"pageSize" form:"pageSize"`
}

type ShareInfoResp struct {
	ShareId      string `json:"shareId" form:"shareId"`
	FileId       string `json:"fileId" form:"fileId"`
	UserId       string `json:"userId" form:"userId"`
	ValidType    int    `json:"validType" form:"validType"`
	ExpireTime   string `json:"expireTime" form:"expireTime"`
	ShareTime    string `json:"shareTime" form:"shareTime"`
	Code         string `json:"code" form:"code"`
	HitCount     int    `json:"hitCount" form:"hitCount"`
	FileName     string `json:"fileName" form:"fileName"`
	FolderType   int    `json:"folderType" form:"folderType"`
	FileCategory int    `json:"fileCategory" form:"fileCategory"`
	FileType     int    `json:"fileType" form:"fileType"`
	FileCover    string `json:"fileCover" form:"fileCover"`
}

type ShareListResp struct {
	TotalCount int64       `json:"totalCount" form:"totalCount"`
	PageSize   int         `json:"pageSize" form:"pageSize"`
	PageNo     int         `json:"pageNo" form:"pageNo"`
	PageTotal  int         `json:"pageTotal" form:"pageTotal"`
	List       interface{} `json:"list" form:"list"`
}

type ShareFileReq struct {
	FileId    string `json:"fileId" form:"fileId" validate:"required"`
	ValidType string `json:"validType" form:"validType" validate:"required"`
	Code      string `json:"code" form:"code"`
}

type CancelShareReq struct {
	ShareIds string `json:"shareIds" form:"shareIds" validate:"required"`
}

type ShowShareReq struct {
	ShareId string `json:"shareId" form:"shareId" validate:"required"`
}

type ShareLoginResp struct {
	ShareTime   string `json:"shareTime" form:"shareTime"`
	ExpireTime  string `json:"expireTime" form:"expireTime"`
	NickName    string `json:"nickName" form:"nickName"`
	FileName    string `json:"fileName" form:"fileName"`
	CurrentUser bool   `json:"currentUser" form:"currentUser"`
	FileId      string `json:"fileId" form:"fileId"`
	Avatar      string `json:"avatar" form:"avatar"`
	UserId      string `json:"userId" form:"userId"`
}

type ShareFileResp struct {
	FileName string `json:"fileName"`
}

type CheckShareReq struct {
	ShareId   string `json:"shareId" form:"shareId" validate:"required"`
	ShareCode string `json:"code" form:"code" validate:"required"`
}

type LoadShareReq struct {
	PageNo   string `json:"pageNo" form:"pageNo"`
	PageSize string `json:"pageSize" form:"pageSize"`
	ShareId  string `json:"shareId" form:"shareId" validate:"required"`
	FilePid  string `json:"filePid" form:"filePid" validate:"required"`
}

type LoadShareResp struct {
	FileId         string `json:"fileId" form:"fileId"`
	FilePid        string `json:"filePid" form:"filePid"`
	FileSize       int64  `json:"fileSize" form:"fileSize"`
	FileName       string `json:"fileName" form:"fileName"`
	FileCover      string `json:"fileCover" form:"fileCover"`
	LastUpdateTime string `json:"lastUpdateTime" form:"lastUpdateTime"`
	FolderType     int    `json:"folderType" form:"folderType"`
	FileCategory   int    `json:"fileCategory" form:"fileCategory"`
	FileType       int    `json:"fileType" form:"fileType"`
	Status         int    `json:"status" form:"status"`
}

type ShareFolderInfoReq struct {
	Path    string `json:"path" form:"path" validate:"required"`
	ShareId string `json:"shareId" form:"shareId"`
}

type ShareFolderInfoResp struct {
	FileName string `json:"fileName" form:"fileName"`
	FileId   string `json:"fileId" form:"fileId"`
}

type GetShareFileReq struct {
	ShareId string `json:"shareId" form:"shareId" validate:"required"`
	FileId  string `json:"fileId" form:"fileId" validate:"required"`
}

type ShareDownloadUrlReq struct {
	ShareId string `json:"shareId" form:"shareId" validate:"required"`
	FileId  string `json:"fileId" form:"fileId" validate:"required"`
}

type SaveShareReq struct {
	ShareId      string `json:"shareId" form:"shareId"`
	ShareFileIds string `json:"shareFileIds" form:"shareFileIds"`
	MyFolderId   string `json:"myFolderId" form:"myFolderId" validate:"required"`
}
