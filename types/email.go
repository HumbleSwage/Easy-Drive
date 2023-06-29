package types

type EmailServiceReq struct {
	Email     string `json:"email" form:"email" validate:"required"`
	CheckCode string `json:"checkCode" form:"checkCode" validate:"required,min=4"`
	Type      int    `json:"type" form:"type"`
}
