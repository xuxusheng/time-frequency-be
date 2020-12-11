package app

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"strings"
)

type ValidError struct {
	Key     string
	Message string
}

func (v *ValidError) Error() string {
	return v.Message
}

type ValidErrors []*ValidError

func (v ValidErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

func BindAndValid(c *gin.Context, v interface{}) (bool, ValidErrors) {
	var errs ValidErrors
	err := c.ShouldBind(v)
	if err != nil {
		// 从 ctx 中取出翻译器
		trans, _ := c.Value("trans").(ut.Translator)

		// 将 ShouldBind 返回的 err 推断为 ValidationErrors，gin 默认使用的是 validator 这个库来做校验
		verrs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 如果没法推断为 validator.ValidationErrors 类型，就直接返回好了，此时 errs 是个空
			return false, errs
		}

		// verrs.Translate(trans) 返回的 ValidationErrorsTranslations 是一个 map[string]string 结构
		for key, value := range verrs.Translate(trans) {
			// 拿到 key、value 后，转化为自己定义的 ValidError 类型
			errs = append(errs, &ValidError{
				Key:     key,
				Message: value, // 错误信息
			})
		}
		return false, errs
	}
	return true, nil
}
