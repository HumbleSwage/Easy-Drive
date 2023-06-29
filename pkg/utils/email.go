package utils

import (
	"easy-drive/consts"
	"math/rand"
	"time"
)

func GenerateEmailCaptcha() string {
	rand.Seed(time.Now().UnixNano())
	length := consts.EmailCaptchaLength
	// 定义验证码字符集
	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	verificationCode := make([]byte, length)
	for i := 0; i < length; i++ {
		verificationCode[i] = charset[rand.Intn(len(charset))]
	}

	return string(verificationCode)
}
