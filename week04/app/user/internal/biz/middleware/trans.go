package middleware

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en2 "github.com/go-playground/validator/v10/translations/en"
	zh2 "github.com/go-playground/validator/v10/translations/zh"
	"go.uber.org/zap"
)

func NewTrans(locale string) (ut.Translator, error) {
	// 修改 gin 框架中的 validator 引擎属性, 实现定制
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, fmt.Errorf("wrong type of validator")
	}
	// 注册一个获取 json 的 tag 的自定义方法
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	zhT := zh.New() // 中文翻译器
	enT := en.New() // 英文翻译器
	// 第一个参数是备用的语言环境，后面的参数是应该支持的语言环境
	uni := ut.New(enT, zhT, enT)
	trans, ok := uni.GetTranslator(locale)
	if !ok {
		return nil, fmt.Errorf("uni.GetTranslator(%s)", locale)
	}

	switch locale {
	case "en":
		zap.S().Error(en2.RegisterDefaultTranslations(v, trans))
	case "zh":
		zap.S().Error(zh2.RegisterDefaultTranslations(v, trans))
	default:
		zap.S().Error(en2.RegisterDefaultTranslations(v, trans))
	}

	return trans, nil
}
