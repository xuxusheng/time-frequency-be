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
	teacher := v1.NewTeacher(userSvc)
	admin := v1.NewAdmin(userSvc)

	// 登录
	apiV1.Post("/login", user.Login)
	// 校验登录状态中间件
	apiV1.Use(middleware.IsLogin())

	// 普通用户就可以调用的用户相关接口
	{
		apiV1.Post("/user/me", user.Me)
		apiV1.Post("/user/update", user.Update)
		apiV1.Post("/user/update-password", user.UpdatePassword)
	}

	// 老师才允许调用的接口
	{
		teacherApi := apiV1.Party("/teacher")
		// 检验是否为老师身份中间件
		teacherApi.Use(middleware.IsTeacher())
		teacherApi.Post("/create-student", teacher.CreateStudent)
		teacherApi.Post("/list-student", teacher.ListStudent)
		teacherApi.Post("/delete-student", teacher.DeleteStudent)
	}

	// 管理员才允许调用的接口
	{
		adminApi := apiV1.Party("/admin")
		adminApi.Use(middleware.IsAdmin())
		adminApi.Post("/list-user", admin.ListUser)
		adminApi.Post("/toggle-admin", admin.ToggleAdmin)
	}

	return app
}
