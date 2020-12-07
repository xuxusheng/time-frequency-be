package app

import (
	"github.com/gin-gonic/gin"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/pkg/convert"
)

// 从 gin.Context 中获取当前页码，其实就是从 url 的 query 中获取
func GetPn(c *gin.Context) int {
	pn := convert.StrTo(c.Query("pn")).MustInt()
	if pn <= 0 {
		return 1
	}
	return pn
}

func GetPs(c *gin.Context) int {
	ps := convert.StrTo(c.Query("ps")).MustInt()
	if ps <= 0 {
		return global.AppSetting.DefaultPageSize
	}
	return ps
}

func GetPageOffset(pn, ps int) int {
	result := 0
	if pn > 0 {
		result = (pn - 1) * ps
	}
	return result
}
