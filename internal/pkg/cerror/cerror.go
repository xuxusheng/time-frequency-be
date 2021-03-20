package cerror

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// 自定义 Code Error 类型接口
type IError interface {
	Error() string     // 返回格式化后的错误信息字符串，用于实现 go 原生的 Error 接口
	Code() int         // 返回错误码
	Msg() string       // 返回错误信息
	Details() []string // 返回错误详情
	Debugs() []string  // 返回错误 Debug 信息
	StatusCode() int   // 返回 HTTP 状态码

	WithMsg(msg string) *Error           // 设置错误信息
	WithDetail(details ...string) *Error // 设置错误详情
	WithDebugs(debugs ...error) *Error   // 设置错误 Debug 信息

	ToResponse() map[string]interface{} // 将错误类型整理成接口返回需要的形式

	clone() *Error // 拷贝一个新的 Error 对象
}

// 自定义 Code Error 类型
// 封装一个带 【错误码】、【错误概述】、【错误详情】、【错误Debug信息】的错误类型，供接口返回使用
type Error struct {
	code       int
	msg        string   // 错误概述
	details    []string // 错误详情
	debugs     []string // 错误 Debug 信息
	statusCode int      // HTTP 状态码
}

// 用来存储所有已定义的错误，防止 errCode 重复
var codes = map[int]string{}

func New(code int, msg string, statusCode ...int) *Error {
	if _, ok := codes[code]; ok {
		// 如果错误码重复定义，直接报错并退出程序
		log.Fatalf("错误码【%d】已存在，请更换一个", code)
	}
	codes[code] = msg

	// go 不支持函数可选参数，以此种方式实现
	var sc int
	if len(statusCode) > 0 {
		sc = statusCode[0]
	}
	// 如果未定义或定义的 http status code 不符合规范，默认使用 500
	if sc < 200 || sc >= 600 {
		sc = http.StatusInternalServerError
	}
	return &Error{
		code:       code,
		msg:        msg,
		statusCode: sc,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf(
		"错误码：%d，错误信息：%s，错误详情：%s，错误Debug信息：%s",
		e.Code(),
		e.Msg(),
		strings.Join(e.Details(), ";"),
		strings.Join(e.Debugs(), ";"),
	)
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Details() []string {
	if e.details == nil {
		return []string{}
	}
	return e.details
}

func (e *Error) Debugs() []string {
	if e.debugs == nil {
		return []string{}
	}
	return e.debugs
}

func (e *Error) StatusCode() int {
	if e.statusCode != 0 {
		return e.statusCode
	}
	return http.StatusInternalServerError
}

func (e *Error) WithMsg(msg string) *Error {
	n := e.clone()
	n.msg = msg
	return n
}

func (e *Error) WithDetail(details ...string) *Error {
	n := e.clone()
	n.details = []string{}
	for _, d := range details {
		n.details = append(n.details, d)
	}
	return n
}

func (e *Error) WithDebugs(debugs ...error) *Error {
	n := e.clone()
	n.debugs = []string{}
	for _, d := range debugs {
		n.debugs = append(n.debugs, d.Error())
	}
	return n
}

func (e *Error) ToResponse() map[string]interface{} {
	d := map[string]interface{}{
		"err_code":    e.Code(),
		"err_msg":     e.Msg(),
		"err_details": e.Details(),
		"err_debugs":  e.Debugs(),
		"data":        map[string]interface{}{},
	}
	// todo 根据当前 mode 判断是否需要将 debug 信息加入返回中，或者增加一个 is_debug 的配置
	// todo 如果在这里调用全局配置的话，就会不够解耦了，但是如果通过参数来传递，又会显得比较麻烦了，需要考虑下
	return d
}

func (e *Error) clone() *Error {
	n := *e
	return &n
}
