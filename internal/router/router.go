package router

import (
	"github.com/iris-contrib/swagger/v12"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/middleware/requestid"
	swaggerFiles "github.com/swaggo/files"
	_ "github.com/xuxusheng/time-frequency-be/docs"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/middleware"
	v1 "github.com/xuxusheng/time-frequency-be/internal/router/api/v1"
)

func NewApp() *iris.Application {

	app := iris.New()
	app.Validator = global.Validator

	app.UseRouter(recover.New())
	// 对请求进行压缩
	app.Use(iris.Compression)

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
	app.Get("/swagger/{any:path}", swagger.WrapHandler(swaggerFiles.Handler))

	// 在请求 header 中查询或生成 x-request-id 头，并调用 c.SetID
	app.Use(requestid.New())
	app.Use(middleware.Translations())
	app.Use(middleware.AccessLog())

	// 业务接口
	apiV1 := app.Party("/api/v1")
	auth := v1.NewAuth()
	user := v1.NewUser()
	{

		// 登录
		apiV1.Post("/login", auth.Login)

		// ---- 以下接口需要登录才能访问
		apiV1.Use(middleware.JWT())
		// 添加用户
		apiV1.Post("/users", user.Create)
		// 删除用户
		apiV1.Delete("/users/{id:uint}", user.Delete)
		// 更新用户信息（通用信息字段）
		apiV1.Put("/users/{id:uint}", user.Update)
		// 更新密码
		apiV1.Put("/users/{id:uint}/password", user.UpdatePassword)
		// 获取用户列表
		apiV1.Get("/users", user.List)
		// 获取单个用户
		apiV1.Get("/users/{id:uint}", user.Get)
		// 获取当前用户的信息
		apiV1.Get("/users/me", user.Me)
	}

	return app
}
