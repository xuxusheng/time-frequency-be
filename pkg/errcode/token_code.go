package errcode

var (
	UnauthorizedUserError = NewError(20000001, "用户不存在或密码错误")

	UnauthorizedTokenEmpty   = NewError(20000002, "鉴权失败，token 不存在")
	UnauthorizedTokenError   = NewError(20000003, "鉴权失败，Token 不合法")
	UnauthorizedTokenExpired = NewError(20000004, "鉴权失败，Token 已过期")
	// todo 这个错误好像有点问题，应该放登录的时候用？
	UnauthorizedTokenGenerate = NewError(20000005, "鉴权失败，Token 生成失败")
)
