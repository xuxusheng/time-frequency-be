package app

import (
	"github.com/kataras/iris/v12"
	"github.com/xuxusheng/time-frequency-be/global"
)

// 从 gin.Context 中获取当前页码，其实就是从 url 的 query 中获取
func GetPn(c iris.Context) int {
	pn := c.URLParamIntDefault("pn", 1)
	if pn <= 0 {
		return 1
	}
	return pn
}

func GetPs(c iris.Context) int {
	ps := c.URLParamIntDefault("ps", global.AppSetting.DefaultPageSize)
	if ps <= 0 {
		return global.AppSetting.DefaultPageSize
	}
	return ps
}
