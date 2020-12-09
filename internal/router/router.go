package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xuxusheng/time-frequency-be/internal/middleware"
	v1 "github.com/xuxusheng/time-frequency-be/internal/router/api/v1"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Translations())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ping": "pong",
		})
	})

	// 业务接口
	apiv1 := r.Group("/api/v1")
	{
		user := v1.NewUser()
		// 添加用户
		apiv1.POST("/users", user.Create)
		// 删除用户
		apiv1.DELETE("/users/:id", user.Delete)
		// 更新用户信息（通用信息字段）
		apiv1.PUT("/users/:id", user.Update)
		// 获取用户列表
		apiv1.GET("/users", user.List)
		// 获取单个用户
		apiv1.GET("/users/:id", user.Get)
	}

	return r
}
