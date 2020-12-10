package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/xuxusheng/time-frequency-be/docs"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/middleware"
	v1 "github.com/xuxusheng/time-frequency-be/internal/router/api/v1"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())

	// swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if global.ServerSetting.RunMode == gin.DebugMode {
		r.Use(gin.Logger())
	}

	r.Use(middleware.RequestID())
	r.Use(middleware.Translations())
	r.Use(middleware.AccessLog())

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
