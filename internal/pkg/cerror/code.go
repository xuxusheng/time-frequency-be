package cerror

import "net/http"

var (
	// --- 通用错误 ---
	Success            = New(0, "", http.StatusOK)
	BadRequest         = New(1000_0400, "输入参数错误，请检查", http.StatusBadRequest)
	Unauthorized       = New(1000_0401, "验证失败", http.StatusUnauthorized)
	Forbidden          = New(1000_0403, "禁止访问", http.StatusForbidden)
	NotFound           = New(1000_0404, "资源不存在", http.StatusNotFound)
	TooManyRequest     = New(1000_0429, "请求过多", http.StatusTooManyRequests)
	ServerError        = New(1000_0500, "服务器内部错误", http.StatusInternalServerError)
	ServiceUnavailable = New(1000_0503, "服务暂时不可用", http.StatusServiceUnavailable)

	// token 或鉴权相关错误
	TokenEmpty   = New(2000_0001, "鉴权失败，token 不存在", http.StatusUnauthorized)
	TokenExpired = New(2000_0002, "鉴权失败，token 已过期", http.StatusUnauthorized)
	TokenInvalid = New(2000_0003, "鉴权失败，token 不合法", http.StatusUnauthorized)

	// --- 其他业务相关错误 ---

	// 用户相关
	Login = New(3000_0001, "用户不存在或密码错误", http.StatusUnauthorized)
)
