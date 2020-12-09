package errcode

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Error struct {
	code    int
	msg     string   // msg 信息存放简短的错误提示，用于前端提示用户。
	details []string // details 存放详细的错误信息，比如 sql 执行抛出的错误等，前端默认不展示给用户。
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码%d已存在，请更换一个", code))
	}
	codes[code] = msg
	return &Error{
		code: code,
		msg:  msg,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"错误码：%d，错误信息：%s, 错误详情：%s",
		e.Code(),
		e.Msg(),
		strings.Join(e.Details(), ";"),
	)
}

func (e *Error) Meta() *gin.H {
	return &gin.H{
		"err_code":    e.Code(),
		"err_msg":     e.Msg(),
		"err_details": e.Details(),
	}
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

func (e *Error) Details() []string {
	if e.details == nil {
		return []string{}
	}
	return e.details
}

// 重写 Error 的 msg 字段
func (e *Error) WithMsg(msgs ...string) *Error {
	newError := *e
	newError.msg = strings.Join(msgs, ";")
	return &newError
}

func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.details = []string{}
	for _, d := range details {
		newError.details = append(newError.details, d)
	}
	return &newError
}

func (e *Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case ServerError.Code():
		return http.StatusInternalServerError
	case NotFound.Code():
		return http.StatusNotFound

		// token 校验相关的错误
	case UnauthorizedTokenEmpty.Code():
		fallthrough
	case UnauthorizedAuthNotExist.Code():
		fallthrough
	case UnauthorizedTokenTimeout.Code():
		fallthrough
	case UnauthorizedTokenGenerate.Code():
		return http.StatusUnauthorized

		// 请求过多
	case TooManyRequest.Code():
		return http.StatusTooManyRequests

		// 输入错误
	case CreateUserFailNameExist.Code():
		fallthrough
	case CreateUserFailPhoneExist.Code():
		fallthrough
	case UpdateUserFailNameExist.Code():
		fallthrough
	case UpdateUserFailPhoneExist.Code():
		fallthrough
	case InvalidParams.Code():
		return http.StatusBadRequest

	}

	// 默认返回 500
	return http.StatusInternalServerError
}
