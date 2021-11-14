package captcha

import (
	"go-advance/week04/app/user/api"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type Generator struct {
	base64Captcha.Captcha
}

func NewCaptcha(store base64Captcha.Store) *Generator {
	return &Generator{
		Captcha: *base64Captcha.NewCaptcha(base64Captcha.DefaultDriverDigit, store),
	}
}

func (c *Generator) Gen(ctx *gin.Context) {
	id, b64s, err := c.Generate()
	if err != nil {
		zap.L().Error("cannot generate captcha", zap.Error(err))
		ctx.String(http.StatusInternalServerError, "生成验证码错误")
		return
	}

	ctx.JSON(http.StatusOK, api.Captcha{
		ID:      id,
		PicPath: b64s,
	})
}
