package errcode

var (
	Success = NewError(0, "成功")

	BadRequest     = NewError(10000400, "输入参数错误，请检查")
	Unauthorized   = NewError(10000401, "验证失败")
	NotFound       = NewError(10000404, "找不到资源")
	TooManyRequest = NewError(10000429, "请求过多")

	InternalServerError = NewError(10000500, "服务器内部错误")
	ServiceUnavailable  = NewError(10000503, "服务暂时不可用")
)
