package consts

const (
	VerifyCodeKey      = "captcha"
	VerifyEmailCodeKey = "email_captcha"
	UserInfo           = "user_info" // 用户信息缓存key
	JwtSecret          = "JwtSecret" // jwt加密密钥

	EmailCaptchaLength     = 5  // 邮箱验证码的长度
	ShareCodeLength        = 5  // 分享码长度
	ShareIdLength          = 15 // 分享id长度
	UserIdLength           = 15 // 用户id的长度
	FileIdLength           = 15 // 文件id的长度
	RandomNumberLength     = 8  // 随机数字长度
	EmailCaptchaExpiration = 15 // 邮箱验证码有效时间（Min）
	TokenExpiration        = 24 // token有效时间（hour）
	DownloadExpiration     = 5  // 下载链接有效时间
	PasswordCost           = 12 // 密码加密强度

	UserRegisterTiTleSettingId   = 4
	UserRegisterContentSettingId = 1                       // 用户注册设置id
	UserRePwdSettingId           = 2                       // 用户更新密码id
	UserInitSpaceId              = 3                       // 用户默认内存总容量
	UserLimitedSpace             = 10 * 1024 * 1024 * 1024 // 单用户内存限制10G
	CompressImageLength          = 100                     // 缩略图高度
)
