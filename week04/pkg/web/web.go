package web

import (
	"net/http"
	"os"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type ErrWrapper struct {
	Trans ut.Translator
}

type AppHandler func(*gin.Context) error

func (w *ErrWrapper) Wrap(handler AppHandler) func(*gin.Context) {
	return func(ctx *gin.Context) {
		// panic
		defer func() {
			if r := recover(); r != nil {
				zap.L().Error("panic", zap.Any("recover", r))
				ctx.JSON(
					http.StatusInternalServerError,
					http.StatusText(http.StatusInternalServerError))
			}
		}()

		err := handler(ctx)
		if err != nil {
			zap.L().Error("Error occurred handling request", zap.Error(err))
			switch err.(type) {
			case UserError:
				userErr := err.(UserError)
				ctx.String(
					http.StatusBadRequest,
					userErr.Message())
				return
			case validator.ValidationErrors:
				// 获取 validator.ValidationErrors 类型的 errors
				errs := err.(validator.ValidationErrors)
				// validator.ValidationErrors 类型错误则进行翻译
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": removeTopTag(errs.Translate(w.Trans)),
				})
				return
			}

			// GRPC err
			if s, ok := status.FromError(err); ok {
				switch s.Code() {
				case codes.NotFound:
					ctx.JSON(http.StatusBadRequest, "user not exist")
				case codes.Unauthenticated:
					ctx.JSON(http.StatusUnauthorized, "wrong account or password")
				default:
					ctx.JSON(http.StatusInternalServerError, http.StatusText(int(s.Code())))
				}
				return
			}

			// system error
			code := http.StatusOK
			switch {
			case os.IsNotExist(err):
				code = http.StatusNotFound
			case os.IsPermission(err):
				code = http.StatusForbidden
			default:
				code = http.StatusInternalServerError
			}
			ctx.JSON(code, http.StatusText(code))
		}
	}
}

func removeTopTag(errTrans validator.ValidationErrorsTranslations) map[string]string {
	rsp := make(map[string]string)
	for field, err := range errTrans {
		rsp[field[strings.LastIndex(field, ".")+1:]] = err
	}

	return rsp
}

type UserError interface {
	error
	Message() string
}
