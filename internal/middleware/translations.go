package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/locales/zh_Hant_TW"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/xuxusheng/time-frequency-be/global"
)

// 国际化中间件
func Translations() gin.HandlerFunc {
	return func(c *gin.Context) {
		uni := ut.New(en.New(), zh.New(), zh_Hant_TW.New())
		locale := c.GetHeader("locale")
		trans, ok := uni.GetTranslator(locale)
		if !ok {
			trans, _ = uni.GetTranslator("zh")
		}

		v, ok := binding.Validator.Engine().(*validator.Validate)
		if ok {
			switch locale {
			case "zh":
				global.Logger.Debug(c, 1)
				_ = zhTranslations.RegisterDefaultTranslations(v, trans)
			case "en":
				global.Logger.Debug(c, 2)
				_ = enTranslations.RegisterDefaultTranslations(v, trans)
			default:
				global.Logger.Debug(c, 3)
				_ = zhTranslations.RegisterDefaultTranslations(v, trans)
			}
			c.Set("trans", trans)
		}
		c.Next()
	}
}
