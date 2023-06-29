package service

import (
	"context"
	"easy-drive/conf"
	"easy-drive/consts"
	"easy-drive/pkg/ctl"
	"easy-drive/pkg/e"
	"easy-drive/pkg/utils"
	"easy-drive/repositry/cache"
	"easy-drive/repositry/dao"
	"easy-drive/types"
	"gopkg.in/mail.v2"
	"strconv"
	"strings"
	"time"
)

func SendEmail(ctx context.Context, req *types.EmailServiceReq) (resp interface{}, err error) {
	code := e.Success

	// 检查邮箱状态
	userDao := dao.NewUserDao(ctx)
	flag, err := userDao.IsEmailExists(req.Email)

	if err != nil {
		utils.LogrusObj.Error("检查邮箱是否存在出错:", err)
		return ctl.RespError(), nil
	}

	// 生成验证码
	captcha := utils.GenerateEmailCaptcha()

	// 定义邮件内容
	var mailText string

	// 发送邮件类型
	switch req.Type {
	case 0: // 用户注册逻辑
		if flag {
			// 邮箱已注册
			code = e.UserAlreadyExistsError
			return ctl.RespError(code), nil
		}

		// 注册邮件内容
		mailText = strings.Join([]string{"您正在进行Easy-Driver注册，验证码：", captcha, "\n注意：验证码有效期仅有", strconv.Itoa(consts.EmailCaptchaExpiration), "分钟"}, "")
	case 1: // 修改密码逻辑
		if !flag {
			// 邮箱尚未注册
			code = e.UserNotRegisterError
			return ctl.RespError(code), nil
		}
		// 忘记密码
		mailText = strings.Join([]string{"您正在进行Easy-Driver密码找回，验证码：", captcha, "\n注意：验证码有效期仅有", strconv.Itoa(consts.EmailCaptchaExpiration), "分钟"}, "")
	default:
		utils.LogrusObj.Info("未知邮件发送类型:", req.Type)
	}

	// 生成缓存的key，注意type 0注册 1修改密码
	key := cache.VerificationCodeCacheKey(req.Type, req.Email)

	// 验证码缓存存入redis
	expiration := consts.EmailCaptchaExpiration * time.Minute
	rdb := cache.RedisClient
	err = rdb.Set(key, captcha, expiration).Err()
	if err != nil {
		utils.LogrusObj.Error("验证码存入redis发生错误:", err)
		return nil, err
	}

	// 发送邮件
	eConfig := conf.Conf.Email
	m := mail.NewMessage()
	m.SetHeader("From", eConfig.SmtpEmail)
	m.SetHeader("To", req.Email)
	m.SetHeader("Subject", "Easy-Driver")
	m.SetBody("text/plain", mailText)
	d := mail.NewDialer(eConfig.SmtpHost, 465, eConfig.SmtpEmail, eConfig.SmtpPassword)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err := d.DialAndSend(m); err != nil {
		utils.LogrusObj.Error("发送邮箱验证码失败:", err)
		return nil, err
	}

	return ctl.RespSuccess(), nil
}
