package types

type SystemSettingReq struct {
	RegisterEmailTitle   string `json:"registerEmailTitle" form:"registerEmailTitle"`
	RegisterEmailContent string `json:"registerEmailContent" form:"registerEmailContent"`
	UserInitUseSpace     string `json:"userInitUseSpace" form:"registerEmailContent"`
}

type SystemSettingResp struct {
	RegisterEmailTitle   string `json:"registerEmailTitle" form:"registerEmailTitle"`
	RegisterEmailContent string `json:"registerEmailContent" form:"registerEmailContent"`
	UserInitUseSpace     string `json:"userInitUseSpace" form:"registerEmailContent"`
}

type LoadUserListReq struct {
	PageNo        string `json:"pageNo" form:"pageNo"`
	PageSize      string `json:"pageNum" form:"pageNum"`
	NickNameFuzzy string `json:"nickNameFuzzy" form:"nickNameFuzzy"`
	Status        string `json:"status" form:"status"`
}

type LoadUserInfo struct {
	UserId        string `json:"userId" form:"userId"`
	NickName      string `json:"nickName" form:"nickName"`
	Email         string `json:"email" form:"email"`
	Avatar        string `json:"avatar" form:"avatar"`
	JoinTime      string `json:"joinTime" form:"joinTime"`
	LastLoginTime string `json:"lastLoginTime" form:"lastLoginTime"`
	Status        int    `json:"status" form:"status"`
	UseSpace      int64  `json:"useSpace" form:"useSpace"`
	TotalSpace    int64  `json:"totalSpace" form:"totalSpace"`
}

type ListInfoResp struct {
	TotalCount int64       `json:"totalCount" form:"totalCount"`
	PageSize   int         `json:"pageSize" form:"pageSize"`
	PageNo     int         `json:"pageNo" form:"pageNo"`
	PageTotal  int         `json:"pageTotal" form:"pageTotal"`
	List       interface{} `json:"list" form:"list"`
}

type UpdateUserStatusReq struct {
	UserId string `json:"userId" form:"userId" validate:"required"`
	Status string `json:"status" form:"status" validate:"required"`
}

type UpdateUserSpaceReq struct {
	UserId      string `json:"userId" form:"userId" validate:"required"`
	ChangeSpace string `json:"changeSpace" form:"changeSpace" validate:"required"`
}

type LoadFileListReq struct {
	PageNo        string `json:"pageNo" form:"pageNo"`
	PageSize      string `json:"pageSize" form:"pageSize"`
	FileNameFuzzy string `json:"fileNameFuzzy" form:"fileNameFuzzy"`
	FilePid       string `json:"filePid" form:"filePid"`
}

type LoadFileInfo struct {
	FileId         string `json:"fileId" form:"fileId"`
	FilePid        string `json:"filePid" form:"FilePid"`
	UserId         string `json:"userId" form:"userId"`
	UserName       string `json:"userName" form:"userName"`
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

type AdminGetFolderReq struct {
	Path string `json:"path" form:"path"`
}

type AdminGetFolderInfoResp struct {
	FileName string `json:"fileName"`
	FileId   string `json:"fileId"`
}

type AdminDelFileReq struct {
	FileIdAndUserIds string `json:"fileIdAndUserIds" form:"fileIdAndUserIds" validate:"required"`
}
