package test

import (
	"easy-drive/repositry/model"
	"fmt"
	"testing"
)

func TestCheckPassword(t *testing.T) {
	user := &model.User{
		UserId:   "4523",
		UserName: "wangwu",
		NickName: "wangwu",
		Email:    "86712526@qq.com",
	}

	user.Password = "78302615c8b79cac8df6d2607f8a83ee"
	fmt.Println(user.CheckPassword("123qwe!@#"))
}
