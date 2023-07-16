package model

import (
	"time"
)

type Share struct {
	ShareId    string
	FileId     string
	UserId     string
	ValidType  int
	ExpireTime *time.Time
	ShareTime  *time.Time
	Code       string
	HitCount   int
}
