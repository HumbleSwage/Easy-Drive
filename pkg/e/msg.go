package e

var MsgFlags = map[int]string{
	Success:       "成功",
	Error:         "失败",
	InvalidParams: "请求参数失败",

	EmailAlreadyExistsError: "邮箱已被注册",
	UserNotRegisterError:    "邮箱尚未被注册",
	UserAlreadyExistsError:  "该用户名已存在",
	UserPasswordError:       "用户名或密码错误",
	UseAccountDisable:       "用户账户被禁用",

	VerificationCodeError:    "验证码不正确",
	EmailCodeError:           "邮箱验证码不正确",
	ParameterValidationError: "参数验证未通过",
	UserSessionExpiration:    "登录超时，请重新登录",
	UpdateAvatarError:        "上传图片失败",
	ErrorJsonType:            "Json类型不匹配",

	UserStoreSpaceError: "用户存储空间不足",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if !ok {
		return MsgFlags[Error]
	}
	return msg
}
