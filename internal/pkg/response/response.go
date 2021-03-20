package response

import (
	"github.com/kataras/iris/v12"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
)

type Response struct {
	ctx iris.Context
}

type Map map[string]interface{}

func New(c iris.Context) *Response {
	return &Response{
		ctx: c,
	}
}

func (r *Response) Success(data ...interface{}) {
	d := cerror.Success.ToResponse()
	if len(data) > 0 && data[0] != nil {
		d["data"] = data[0]
	}
	r.ctx.StopWithJSON(cerror.Success.StatusCode(), d)
}

//func (r *Response) SuccessList {
//
//}

func (r *Response) Error(err cerror.IError) {
	r.ctx.StopWithJSON(err.StatusCode(), err.ToResponse())
}
