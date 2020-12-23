package v1

import (
	"github.com/kataras/iris/v12"
	"github.com/xuxusheng/time-frequency-be/internal/service"
	"github.com/xuxusheng/time-frequency-be/pkg/app"
)

type Auth struct {
}

func NewAuth() Auth {
	return Auth{}
}

type LoginReq struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// 登录 godoc
// @summary 登录
// @description 使用用户名、密码登录，返回 JWT Token 字符串
// @accept json
// @produce json
// @tags auth
// @param name body string true "用户名"
// @param phone body string true "密码"
// @success 200 {object} model.Resp{data=string}
// @router /api/v1/login [post]
func (a Auth) Login(c iris.Context) {

	param := LoginReq{}
	resp := app.NewResponse(c)

	// 校验参数
	if !app.BindAndValid(c, &param) {
		return
	}

	userSvc := service.NewUserService(c.Request().Context())

	token, cerr := userSvc.Login(param.Name, param.Password)
	if cerr != nil {
		resp.ToError(cerr)
		return
	}
	resp.ToSuccess(token)
}
