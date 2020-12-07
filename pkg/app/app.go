package app

import (
	"github.com/gin-gonic/gin"
	"github.com/xuxusheng/time-frequency-be/pkg/errcode"
	"net/http"
)

type PageInfo struct {
	Pn    int `json:"pn"`
	Ps    int `json:"ps"`
	Total int `json:"total"`
}

type Response struct {
	Ctx *gin.Context
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{
		Ctx: ctx,
	}
}

func (r *Response) ToResponse(data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	r.Ctx.JSON(http.StatusOK, gin.H{
		"meta": errcode.Success.Meta(),
		"data": data,
	})
}

func (r *Response) ToResponseList(list interface{}, total int64) {
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

func (r *Response) ToErrorResponse(err *errcode.Error) {
	r.Ctx.JSON(err.StatusCode(), gin.H{
		"meta": err.Meta(),
		"data": gin.H{},
	})
}
