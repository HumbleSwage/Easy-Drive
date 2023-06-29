package model

import (
	"crypto/md5"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"time"
)

type User struct {
	gorm.Model
	UserId        string `gorm:"unique"`
	UserName      string `gorm:"unique"`
	NickName      string
	Email         string `gorm:"unique"`
	Password      string
	LastLoginTime *time.Time
	Authority     bool  `gorm:"default=false"`
	Status        int8  `gorm:"default=1"`
	UseSpace      int64 `gorm:"default=0"`
	TotalSpace    int64
	Avatar        string
}

func (u *User) SetPassword(password string) (err error) {
	//bytes, err := bcrypt.GenerateFromPassword([]byte(password), consts.PasswordCost)
	//if err != nil {
	//	return
	//}
	//u.Password = string(bytes)
	//return
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(password))
	cipherStr := md5Ctx.Sum(nil)
	u.Password = fmt.Sprintf("%x", cipherStr)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	// 目前直接基于MD5做的密码验证
	//err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	//return err == nil
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(password))
	cipherStr := md5Ctx.Sum(nil)
	md5str := fmt.Sprintf("%x", cipherStr)
	return reflect.DeepEqual(md5str, u.Password)
}
