package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/service"
	"github.com/xuxusheng/time-frequency-be/pkg/app"
	"github.com/xuxusheng/time-frequency-be/pkg/convert"
	"github.com/xuxusheng/time-frequency-be/pkg/errcode"
	"gorm.io/gorm"
)

type User struct {
}

func NewUser() User {
	return User{}
}

type CreateUserReq struct {
	Name     string `form:"name" binding:"required,min=4,max=20"`
	Phone    string `form:"phone" binding:"required"`
	Password string `form:"password" binding:"required,min=4"`
}

// todo，加入一个角色字段，添加时允许选择是 admin 还是 member
func (u User) Create(c *gin.Context) {
	param := CreateUserReq{}
	resp := app.NewResponse(c)

	// 校验参数
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		// 参数校验失败
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		resp.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	userSvc := service.NewUserService(c.Request.Context())

	// service 返回的 err，已经是自定义的 errcode.Error 类型了
	_, err := userSvc.Create(param.Name, param.Phone, param.Password)
	if err != nil {
		global.Logger.Error(c, "userSvc.Create err: %v", err)
		resp.ToErrorResponse(err)
		return
	}

	resp.ToResponse(nil)
}

func (u User) Delete(c *gin.Context) {
	resp := app.NewResponse(c)

	// 检查 id 格式
	id, err := convert.StrTo(c.Param("id")).UInt()
	if err != nil {
		resp.ToErrorResponse(errcode.InvalidParams.WithMsg("id 格式错误"))
		return
	}

	userSvc := service.NewUserService(c.Request.Context())

	// 执行删除
	if err := userSvc.Delete(id); err != nil {
		resp.ToErrorResponse(err)
		return
	}
	resp.ToResponse(nil)
	return
}

type UpdateUserReq struct {
	Name  string `form:"name" binding:"required,min=4,max=20"`
	Phone string `form:"phone" binding:"required"`
}

func (u User) Update(c *gin.Context) {
	param := UpdateUserReq{}
	resp := app.NewResponse(c)

	// 检查 id 格式
	id, err := convert.StrTo(c.Param("id")).UInt()
	if err != nil {
		resp.ToErrorResponse(errcode.InvalidParams.WithMsg("id 格式错误"))
		return
	}

	// 校验参数
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		// 参数校验失败
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		resp.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	userSvc := service.NewUserService(c.Request.Context())

	if err := userSvc.Update(id, param.Name, param.Phone); err != nil {
		resp.ToErrorResponse(err)
		return
	}

	resp.ToResponse(nil)
}

type UserListReq struct {
	Name  string `form:"id" binding:""`
	Phone string `form:"phone" binding:""`
	Pn    string `form:"pn"`
	Ps    string `form:"ps"`
}

func (u User) List(c *gin.Context) {
	resp := app.NewResponse(c)
	param := UserListReq{}

	// 这里的 errs 类型不是 errcode 包中定义的 Error，而是 app 包中定义的 ValidError
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		// 参数校验失败
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		resp.ToErrorResponse(
			errcode.InvalidParams.WithDetails(errs.Errors()...),
		)
		return
	}

	pn := app.GetPn(c)
	ps := app.GetPs(c)
	userSvc := service.NewUserService(c.Request.Context())

	users, count, err := userSvc.List(param.Phone, param.Name, pn, ps)
	if err != nil {
		resp.ToErrorResponse(errcode.GetUserListFail.WithDetails(err.Error()))
		return
	}
	resp.ToResponseList(users, count)
}

func (u User) Get(c *gin.Context) {

	resp := app.NewResponse(c)
	idStr := c.Param("id")

	id, err := convert.StrTo(idStr).UInt()
	if err != nil {
		resp.ToErrorResponse(
			errcode.InvalidParams.WithMsg(
				"id 格式错误",
			),
		)
		return
	}

	userSvc := service.NewUserService(c.Request.Context())

	user, err := userSvc.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.ToErrorResponse(errcode.NotFound.WithMsg("用户不存在"))
			return
		}
		resp.ToErrorResponse(errcode.GetUserFail.WithDetails(err.Error()))
		return
	}
	resp.ToResponse(user)
	return
}
