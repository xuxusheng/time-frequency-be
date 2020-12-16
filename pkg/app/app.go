package app

import (
	"github.com/kataras/iris/v12"
	"github.com/xuxusheng/time-frequency-be/pkg/errcode"
)

type Response struct {
	Ctx iris.Context
}

func NewResponse(ctx iris.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}

func (r *Response) ToSuccess(data interface{}) {
	if data == nil {
		data = iris.Map{}
	}

	r.Ctx.JSON(iris.Map{
		"meta": errcode.Success.Meta(),
		"data": data,
	})
}

func (r *Response) ToSuccessList(list interface{}, total int) {

	r.Ctx.JSON(iris.Map{
		"meta": errcode.Success.Meta(),
		"data": iris.Map{
			"data":  list,
			"pn":    GetPn(r.Ctx),
			"ps":    GetPs(r.Ctx),
			"total": total,
		},
	})
}

func (r *Response) ToError(err *errcode.Error) {
	r.Ctx.StatusCode(err.StatusCode())
	r.Ctx.JSON(iris.Map{
		"meta": err.Meta(),
		"data": iris.Map{},
	})
}
