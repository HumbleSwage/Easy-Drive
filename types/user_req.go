package types

import "mime/multipart"

type UserRegisterReq struct {
	Email     string `json:"email" form:"email" validate:"required"`
	NickName  string `json:"nickName" form:"nickName" validate:"required"`
	Password  string `json:"password" form:"password" validate:"required"`
	EmailCode string `json:"emailCode" form:"emailCode" validate:"required"`
	CheckCode string `json:"checkCode" form:"checkCode" validate:"required"`
}

type UserLoginReq struct {
	Email     string `json:"email" form:"email" validate:"required"`
	Password  string `json:"password" form:"password" validate:"required"`
	CheckCode string `json:"checkCode" form:"checkCode" validate:"required"`
}

type UserResetPwdReq struct {
	Email     string `json:"email" form:"email" validate:"required"`
	Password  string `json:"password" form:"password" validate:"required"`
	CheckCode string `json:"checkCode" form:"checkCode" validate:"required"`
	EmailCode string `json:"emailCode" form:"emailCode" validate:"required"`
}

type UpdateUserAvatarReq struct {
	File     multipart.File `json:"file" form:"file"`
	FileSize int64          `json:"fileSize" form:"fileSize"`
}

type UpdatePasswordReq struct {
	Password string `json:"password" form:"password" validate:"required"`
}
