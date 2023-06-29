package e

const (
	Success       = 200
	Error         = 500
	InvalidParams = 400

	EmailAlreadyExistsError = 10001
	UserAlreadyExistsError  = 10002
	UserSessionExpiration   = 10003
	UserStoreSpaceError     = 10004
	UserNotRegisterError    = 10005
	UserPasswordError       = 10006
	UseAccountDisable       = 10007

	VerificationCodeError    = 20001
	EmailCodeError           = 20002
	ParameterValidationError = 20003
	StoreInSessionError      = 20004
	UpdateAvatarError        = 20006
	ErrorJsonType            = 20009
)
