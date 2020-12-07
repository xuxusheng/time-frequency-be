package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/service"
	"github.com/xuxusheng/time-frequency-be/pkg/app"
	"github.com/xuxusheng/time-frequency-be/pkg/errcode"
)

type User struct {
}

func NewUser() User {
	return User{}
}

func (u User) List(c *gin.Context) {
	param := service.UserListReq{}

	resp := app.NewResponse(c)

	// 这里的 errs 类型不是 errcode 包中定义的 Error，而是 app 包中定义的 ValidError
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		// 参数校验失败
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		resp.ToErrorResponse(
			errcode.InvalidParams.WithDetails(
				errs.Errors()...,
			),
		)
		return
	}

	svc := service.New(c.Request.Context())
	pageInfo := app.PageInfo{
		Pn: app.GetPn(c),
		Ps: app.GetPs(c),
	}
	total, err := svc.CountUser(&service.CountUserReq{
		Name:  param.Name,
		Phone: param.Phone,
	})
	if err != nil {
		global.Logger.Errorf(c, "svc.CountTag err: %v", err)
		resp.ToErrorResponse(
			errcode.CountUserFail.WithDetails(err.Error()),
		)
		return
	}

	users, err := svc.GetUserList(&param, &pageInfo)
	if err != nil {
		global.Logger.Errorf(c, "svc.GetUserList err: %v", err)
		resp.ToErrorResponse(
			errcode.GetUserListFail.WithDetails(err.Error()),
		)
		return
	}

	resp.ToResponseList(users, total)
	return
}

func (u User) Create(c *gin.Context) {
	param := service.CreateUserReq{}
	resp := app.NewResponse(c)

	// 校验参数
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		// 参数校验失败
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		resp.ToErrorResponse(
			errcode.InvalidParams.WithDetails(
				errs.Errors()...,
			),
		)
		return
	}

	svc := service.New(c.Request.Context())

	// todo 判断用户名和手机号是否重复
	isExit, err := svc.IsUserExist(param.Name, "")
	if err != nil {
		// 查询用户名是否被占用失败
		global.Logger.Error(c, "svc.IsUserExist err: %v", err)
		resp.ToErrorResponse(
			errcode.CreateUserFail.WithDetails(
				err.Error(),
			),
		)
		return
	}
	if isExit {
		resp.ToErrorResponse(
			errcode.CreateUserFail.WithDetails(
				"用户名已被占用",
			),
		)
		return
	}

	isExit, err = svc.IsUserExist("", param.Phone)
	if err != nil {
		// 查询用户名是否被占用失败
		global.Logger.Error(c, "svc.IsUserExist err: %v", err)
		resp.ToErrorResponse(
			errcode.CreateUserFail.WithDetails(
				err.Error(),
			),
		)
		return
	}
	if isExit {
		resp.ToErrorResponse(
			errcode.CreateUserFail.WithDetails(
				"手机号已被占用",
			),
		)
		return
	}

	err = svc.CreateUser(&param)
	if err != nil {
		global.Logger.Error(c, "svc.CreateUser err: %v", err)
		resp.ToErrorResponse(
			errcode.CreateUserFail.WithDetails(
				err.Error(),
			),
		)
		return
	}

	resp.ToResponse(nil)
}
