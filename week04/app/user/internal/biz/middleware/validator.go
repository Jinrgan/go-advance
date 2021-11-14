package middleware

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func RegisterValidation(trans ut.Translator) error {
	if vld, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := vld.RegisterValidation("mobile", ValidateMobile)
		if err != nil {
			zap.L().Error("cannot 01register validation", zap.Error(err))
			return fmt.Errorf("cannot 01register validation: %v", err)
		}
		err = vld.RegisterTranslation("mobile", trans, registerFn, translationFn)
	}

	return nil
}

func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	ok, err := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if err != nil {
		zap.L().Error("cannot match mobile", zap.Error(err))
		return false
	}

	return ok
}

func registerFn(ut ut.Translator) error {
	return ut.Add("mobile", "{0} 非法的手机号码!", true) // see universal-translator for details
}

func translationFn(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("mobile", fe.Field())
	return t
}
