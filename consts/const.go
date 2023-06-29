package consts

const (
	VerifyCodeKey      = "captcha"
	VerifyEmailCodeKey = "email_captcha"
	UserInfo           = "user_info" // 用户信息缓存key
	JwtSecret          = "JwtSecret" // jwt加密密钥

	EmailCaptchaLength     = 5  // 邮箱验证码的长度
	UserIdLength           = 15 // 用户id的长度
	FileIdLength           = 15 // 文件id的长度
	RandomNumberLength     = 5  // 随机数字长度
	EmailCaptchaExpiration = 15 // 邮箱验证码有效时间（Min）
	TokenExpiration        = 24 // token有效时间（hour）
	PasswordCost           = 12 // 密码加密强度

	UserInitSpace       = 5   // 用户默认内存总容量
	CompressImageLength = 100 // 缩略图高度
)
