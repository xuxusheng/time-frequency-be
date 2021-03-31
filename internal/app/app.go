package app

import (
	"github.com/kataras/iris/v12"
	"github.com/xuxusheng/time-frequency-be/global"
	v1 "github.com/xuxusheng/time-frequency-be/internal/api/v1"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/middleware"
	"github.com/xuxusheng/time-frequency-be/internal/service"
)

func New() *iris.Application {

	app := iris.New()
	app.Validator = global.Validator

	app.UseRouter(iris.Compression)

	app.Use(middleware.Translations())

	apiV1 := app.Party("/api/v1")

	userSvc := service.NewUser(dao.NewUser(global.DB))
	user := v1.NewUser(userSvc)
	{
		// 登录
		apiV1.Post("/login", user.Login)

		apiV1.Use(middleware.IsLogin())

		// 用户 CRUD
		//apiV1.Post("/user/create", user.Create)
		apiV1.Post("/user/me", user.Me)
		//apiV1.Post("/user/list", user.List)
		apiV1.Post("/user/update", user.Update)
		//apiV1.Post("/user/delete", user.Delete)
	}

	teacher := v1.NewTeacher(userSvc)
	{
		apiV1.Post("/teacher/create-user", teacher.CreateUser)
		apiV1.Post("/teacher/list-user", teacher.ListUser)
		apiV1.Post("/teacher/delete-user", teacher.DeleteUser)
	}

	return app
}
