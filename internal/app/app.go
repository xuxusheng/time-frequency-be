package app

import (
	"github.com/kataras/iris/v12"
	"github.com/xuxusheng/time-frequency-be/global"
	v1 "github.com/xuxusheng/time-frequency-be/internal/controller/v1"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/service"
)

func New() *iris.Application {

	app := iris.New()

	apiV1 := app.Party("/api/v1")

	userSvc := service.NewUser(dao.NewUser(global.DB))
	user := v1.NewUser(userSvc)
	{
		apiV1.Post("/users", user.Create)
	}

	return app
}
