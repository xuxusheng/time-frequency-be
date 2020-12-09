package errcode

var (
	Success        = NewError(0, "成功")
	ServerError    = NewError(10000000, "服务器内部错误")
	InvalidParams  = NewError(10000001, "输入参数错误")
	NotFound       = NewError(10000002, "找不到资源")
	TooManyRequest = NewError(10000003, "请求过多")
)
