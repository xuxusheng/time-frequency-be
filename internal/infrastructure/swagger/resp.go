package swagger

type Resp struct {
	ErrCode    int         `json:"err_code" example:"0"` // 错误码
	ErrMsg     string      `json:"err_msg" example:""`   // 错误信息
	ErrDetails []string    `json:"err_details"`          // 错误详细信息
	ErrDebugs  []string    `json:"err_debugs"`           // 错误调式信息
	Data       interface{} `json:"data"`
}

type DWithP struct {
	Pn    int         `json:"pn" example:"1"`      // 当前页码
	Ps    int         `json:"ps" example:"10"`     // 每页显示多少条记录
	Total int         `json:"total" example:"199"` // 总共多少条记录
	Data  interface{} `json:"data"`
}
