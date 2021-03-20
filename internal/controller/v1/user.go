package v1

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/infrastructure/database"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/response"
	"github.com/xuxusheng/time-frequency-be/internal/service"
	"time"
)

type IUser interface {
	Create()
}

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (u User) Create(c iris.Context) {

	ctx := c.Request().Context()
	resp := response.New(c)

	s := time.Now().String()
	name := "name"
	phone := s + "phone"
	email := s + "email"
	password := s + "password"

	db, err := database.New(context.Background(), "postgres://postgres:1234@localhost:5432/example2?sslmode=disable")
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	svc := service.NewUser(dao.NewUser(db))

	user, err := svc.Create(ctx, 1, name, phone, email, password)

	if err != nil {
		if cerr, ok := err.(cerror.IError); ok {
			resp.Error(cerr)
			return
		}
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}

	resp.Success(user)
}
