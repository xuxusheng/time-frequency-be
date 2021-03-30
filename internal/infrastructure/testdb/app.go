package testdb

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
)

func NewApp() *iris.Application {
	app := iris.New()
	app.Validator = validator.New()
	return app
}
