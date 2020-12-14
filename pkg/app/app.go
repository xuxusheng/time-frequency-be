package app

import (
	"github.com/gin-gonic/gin"
	"github.com/xuxusheng/time-frequency-be/pkg/errcode"
	"net/http"
)

type Response struct {
	Ctx *gin.Context
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}

func (r *Response) ToSuccess(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, gin.H{
		"meta": errcode.Success.Meta(),
		"data": data,
	})
}

func (r *Response) ToSuccessList(list interface{}, total int) {
	r.Ctx.JSON(http.StatusOK, gin.H{
		"meta": errcode.Success.Meta(),
		"data": gin.H{
			"data":  list,
			"pn":    GetPn(r.Ctx),
			"ps":    GetPs(r.Ctx),
			"total": total,
		},
	})
}

func (r *Response) ToError(err *errcode.Error) {
	r.Ctx.JSON(err.StatusCode(), gin.H{
		"meta": err.Meta(),
		"data": gin.H{},
	})
}
