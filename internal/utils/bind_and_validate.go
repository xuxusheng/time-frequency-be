package utils

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/response"
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
	return strings.Join(v.Errors(), ";")
}

func (v ValidErrors) Errors() []string {
	var errs []string
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return errs
}

func BindAndValidate(c iris.Context, v interface{}) bool {

	var errs ValidErrors

	err := c.ReadBody(v)

	// 校验通过，没有错误
	if err == nil {
		return true
	}

	resp := response.New(c)
	// 从 ctx 中取出翻译器
	trans, _ := c.Values().Get("trans").(ut.Translator)
	// 将 ShouldBind 返回的 err 推断为 ValidationErrors，gin 默认使用的是 validator 这个库来做校验
	verrs, ok := err.(validator.ValidationErrors)
	if !ok {
		// todo 如果 body 为空的话，会触发此错误，待优化
		// 如果没法推断为 validator.ValidationErrors 类型，就直接返回好了
		resp.Error(cerror.BadRequest.WithDebugs(err))
		return false
	}

	// verrs.Translate(trans) 返回的 ValidationErrorsTranslations 是一个 map[string]string 结构
	for key, value := range verrs.Translate(trans) {
		// 拿到 key、value 后，转化为自己定义的 ValidError 类型
		errs = append(errs, &ValidError{
			Key:     key,
			Message: value, // 错误信息
		})
	}

	// 直接将错误写入 response，外部调用方就不再重复写了
	resp.Error(cerror.BadRequest.WithDetails(errs.Errors()...))
	return false
}
