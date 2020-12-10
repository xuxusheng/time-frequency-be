package model

// Swagger 文档使用的一些 struct
type Meta struct {
	ErrCode    int      `json:"err_code" example:"0"` // 错误码
	ErrMsg     string   `json:"err_msg" example:"成功"` // 错误信息
	ErrDetails []string `json:"err_details"`          // 错误详情
}

type ErrMeta struct {
	ErrCode    int      `json:"err_code" example:"10000000"` // 错误码
	ErrMsg     string   `json:"err_msg" example:"服务器内部错误"`   // 错误信息
	ErrDetails []string `json:"err_details"`                 // 错误详情
}

type Resp struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type ErrResp struct {
	Meta ErrMeta     `json:"meta"`
	Data interface{} `json:"data"`
}

type DWithP struct {
	Pn    int         `json:"pn"`    // 当前页码
	Ps    int         `json:"ps"`    // 每页显示记录数
	Total int         `json:"total"` // 总共多少条记录
	Data  interface{} `json:"data"`
}
