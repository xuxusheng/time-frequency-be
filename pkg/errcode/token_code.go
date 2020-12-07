package errcode

var (
	UnauthorizedTokenEmpty    = NewError(20000001, "鉴权失败，token 为空")
	UnauthorizedAuthNotExist  = NewError(20000002, "鉴权失败，找不到对应的 AppKey 和 AppSecret")
	UnauthorizedTokenError    = NewError(20000003, "鉴权失败，Token 错误")
	UnauthorizedTokenTimeout  = NewError(20000004, "鉴权失败，Token超时")
	UnauthorizedTokenGenerate = NewError(20000005, "鉴权失败，Token生成失败")
)
