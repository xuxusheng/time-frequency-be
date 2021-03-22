package v1

import (
	"github.com/kataras/iris/v12"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/response"
	"github.com/xuxusheng/time-frequency-be/internal/service"
	"time"
)

type IUser interface {
	Create(c iris.Context)
	Get(c iris.Context)
	Update(c iris.Context)
	Delete(c iris.Context)
	IsNameExist(c iris.Context)
	IsPhoneExist(c iris.Context)
	IsEmailExist(c iris.Context)
}

type User struct {
	userSvc service.IUser
}

func NewUser(userSvc service.IUser) *User {
	return &User{userSvc: userSvc}
}

func (u User) Create(c iris.Context) {

	ctx := c.Request().Context()
	resp := response.New(c)

	s := time.Now().String()
	name := s + "name"
	phone := s + "phone"
	email := s + "email"
	password := s + "password"

	user, err := u.userSvc.Create(ctx, 1, name, phone, email, password)

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
