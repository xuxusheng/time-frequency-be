package router

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	_ "github.com/xuxusheng/time-frequency-be/docs"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/middleware"
	v1 "github.com/xuxusheng/time-frequency-be/internal/router/api/v1"
)

func NewApp() *iris.Application {

	app := iris.New()
	app.Validator = global.Validator

	app.UseRouter(recover.New())

	// 健康检查
	app.Get("/liveness", func(c iris.Context) {
		c.StopWithStatus(iris.StatusNoContent)
	})
	app.Get("/readiness", func(c iris.Context) {
		if global.PGEngine != nil {
			c.StopWithStatus(iris.StatusNoContent)
			return
		}
		c.StopWithStatus(iris.StatusServiceUnavailable)
	})

	// swagger 文档
	//app.Get("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//if global.ServerSetting.RunMode == gin.DebugMode {
	//app.Use(gin.Logger())
	//}

	//app.Use(middleware.RequestID())
	app.Use(middleware.Translations())
	//app.Use(middleware.AccessLog())

	app.Get("/ping", func(c iris.Context) {
		c.JSON(iris.Map{
			"message": "pong",
		})
	})

	// 业务接口
	apiv1 := app.Party("/api/v1")
	{
		user := v1.NewUser()
		// 添加用户
		apiv1.Post("/users", user.Create)
		// 删除用户
		apiv1.Delete("/users/{id:uint}", user.Delete)
		// 更新用户信息（通用信息字段）
		apiv1.Put("/users/{id:int:uint}", user.Update)
		// 获取用户列表
		apiv1.Get("/users", user.List)
		// 获取单个用户
		apiv1.Get("/users/{id:uint}", user.Get)
	}

	return app
}
