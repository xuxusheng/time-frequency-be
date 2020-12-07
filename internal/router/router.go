package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/xuxusheng/time-frequency-be/internal/router/api/v1"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"ping": "pong",
		})
	})

	// 业务接口
	apiv1 := r.Group("/api/v1")
	{
		user := v1.NewUser()
		apiv1.POST("/users", user.Create)
		apiv1.GET("/users", user.List)
	}

	return r
}
