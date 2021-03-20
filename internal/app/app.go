package app

import (
	"github.com/kataras/iris/v12"
	v1 "github.com/xuxusheng/time-frequency-be/internal/controller/v1"
)

func New() *iris.Application {

	app := iris.New()

	apiV1 := app.Party("/api/v1")
	user := v1.NewUser()
	{
		apiV1.Post("/users", user.Create)
	}

	return app
}
