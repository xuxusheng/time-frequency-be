package middleware

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/kataras/iris/v12"
	"github.com/xuxusheng/time-frequency-be/global"
)

// 国际化中间件
func Translations() iris.Handler {
	return func(c iris.Context) {
		uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
		locale := c.GetHeader("locale")
		trans, ok := uni.GetTranslator(locale)
		if !ok {
			trans, _ = uni.GetTranslator("zh")
		}

		switch locale {
		case "zh":
			_ = zhTranslations.RegisterDefaultTranslations(global.Validator, trans)
		case "en":
			_ = enTranslations.RegisterDefaultTranslations(global.Validator, trans)
		default:
			_ = zhTranslations.RegisterDefaultTranslations(global.Validator, trans)
		}
		c.Values().Set("trans", trans)
		c.Next()
	}
}
