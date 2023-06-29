package handler

import (
	"bytes"
	"easy-drive/consts"
	"easy-drive/pkg/ctl"
	"easy-drive/pkg/e"
	"easy-drive/service"
	"easy-drive/types"
	"github.com/dchest/captcha"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func CheckCode() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 解析url参数
		kind, err := strconv.Atoi(ctx.Query("type"))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		// 生成图片验证码，并直接返回
		Captcha(ctx, kind, 4)
	}
}

func SendEmailCode() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 请求模版
		var req types.EmailServiceReq

		// 绑定数据
		err := ctx.ShouldBind(&req)
		if err == nil {
			// 参数校验
			if !ValidateParams(ctx, &req) {
				return
			}

			// 定义返回
			var resp interface{}

			// 校验验证码
			if CaptchaVerify(ctx, 1, req.CheckCode) {
				// 发送邮件服务
				resp, err = service.SendEmail(ctx.Request.Context(), &req)
				if err != nil {
					ctx.JSON(http.StatusInternalServerError, err)
					return
				}
			} else {
				// 图片验证码不正确
				code := e.VerificationCodeError
				resp = ctl.RespError(code)
			}

			// 正常响应
			ctx.JSON(http.StatusOK, resp)
			return
		}

		// 参数绑定失败
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
	}
}

// Captcha 生成验证码
func Captcha(c *gin.Context, kind int, length ...int) {
	l := captcha.DefaultLen
	w, h := 80, 36
	if len(length) == 1 {
		l = length[0]
	}
	if len(length) == 2 {
		w = length[1]
	}
	if len(length) == 3 {
		h = length[2]
	}
	captchaId := captcha.NewLen(l)
	session := sessions.Default(c)
	if kind == 0 {
		session.Set(consts.VerifyCodeKey, captchaId)
	} else if kind == 1 {
		session.Set(consts.VerifyEmailCodeKey, captchaId)
	}
	_ = session.Save()
	_ = Serve(c.Writer, c.Request, captchaId, ".png", "zh", false, w, h)
}

// Serve 返回图片
func Serve(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		_ = captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		_ = captcha.WriteImage(&content, id, width, height)
	default:
		return captcha.ErrNotFound
	}

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}

// CaptchaVerify 验证逻辑
func CaptchaVerify(c *gin.Context, kind int, code string) bool {
	session := sessions.Default(c)
	key := ""
	if kind == 0 {
		key = consts.VerifyCodeKey
	} else if kind == 1 {
		key = consts.VerifyEmailCodeKey
	}
	if captchaId := session.Get(key); captchaId != nil {
		session.Delete(key)
		_ = session.Save()
		if captcha.VerifyString(captchaId.(string), code) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
