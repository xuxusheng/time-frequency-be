package v1

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/service"
	"github.com/xuxusheng/time-frequency-be/internal/utils"
	"github.com/xuxusheng/time-frequency-be/pkg/app"
	"github.com/xuxusheng/time-frequency-be/pkg/errcode"
)

type User struct {
}

// --- C --- //
func NewUser() User {
	return User{}
}

type CreateUserReq struct {
	Name     string `json:"name" validate:"required,alphanum,min=6,max=20"`
	Phone    string `json:"phone" validate:"required,numeric,len=11"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

// 创建用户 godoc
// @summary 创建新用户
// @Description 创建新用户接口，专供管理平台调用
// @Accept json
// @Produce json
// @Tags user
// @Param name body string true "用户名（6 - 20 位数字或字母构成)" minlength(6) maxlength(20)
// @Param phone body string true "手机号（11 位数字）" minlength(11) maxlength(11)
// @Param password body string true "密码（6 - 20 位字符）"
// @Success 200 {object} model.Resp{data=model.User}
// @Failure 500 {object} model.ErrResp "内部错误"
// @Router /api/v1/users [post]
func (u User) Create(c iris.Context) {
	param := CreateUserReq{}
	resp := app.NewResponse(c)

	// 校验参数
	if !app.BindAndValid(c, &param) {
		return
	}

	userSvc := service.NewUserService(c.Request().Context())

	// service 返回的 err，已经是自定义的 errcode.ToError 类型了
	user, err := userSvc.Create(param.Name, param.Phone, param.Password)
	if err != nil {
		global.Logger.Error(c, "userSvc.Create err: %v", err)
		resp.ToError(err)
		return
	}

	resp.ToSuccess(user)
}

// --- R --- //
type UserListReq struct {
	Name  string `form:"name" binding:""`
	Phone string `form:"phone" binding:""`
	Pn    string `form:"pn"`
	Ps    string `form:"ps"`
}

// 用户列表 godoc
// @summary 分页获取用户列表
// @description 通过 name、phone 字段查询匹配的用户，支持模糊查询、分页
// @accept json
// @produce json
// @tags user
// @param name query string false "用户名" default()
// @param phone query string false "手机号" default()
// @param pn query string false "第几页" default(1)
// @param ps query string false "每页记录数量" default(10)
// @success 200 {object} model.Resp{data=model.DWithP{data=[]model.User}}
// @router /api/v1/users [get]
func (u User) List(c iris.Context) {
	resp := app.NewResponse(c)
	param := UserListReq{}

	// 校验参数
	if !app.BindAndValid(c, &param) {
		return
	}

	pn := app.GetPn(c)
	ps := app.GetPs(c)
	userSvc := service.NewUserService(c.Request().Context())

	users, count, err := userSvc.List(param.Name, param.Phone, &model.Page{Pn: pn, Ps: ps})
	if err != nil {
		resp.ToError(errcode.GetUserListFail.WithDetails(err.Error()))
		return
	}
	resp.ToSuccessList(users, count)
}

// 查询单个用户 godoc
// @summary 查询单个用户
// @description 通过 ID 查询单个用户详细信息
// @accept json
// @produce json
// @tags user
// @param id path string false "用户ID"
// @success 200 {object} model.Resp{data=model.User}
// @router /api/v1/users/{id} [get]
func (u User) Get(c iris.Context) {
	resp := app.NewResponse(c)

	id, _ := c.Params().GetUint("id")

	userSvc := service.NewUserService(c.Request().Context())

	user, err := userSvc.Get(id)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			resp.ToError(errcode.NotFound.WithMsg("用户不存在"))
			return
		}
		resp.ToError(errcode.GetUserFail.WithDetails(err.Error()))
		return
	}
	resp.ToSuccess(user)
	return
}

// 获取当前用户信息 godoc
// @summary 获取当前用户信息
// @description 返回当前已登录账号的用户信息
// @produce json
// @tags user
// @success 200 {object} model.Resp{data=model.User}
// @router /api/v1/users/me [get]
func (u User) Me(c iris.Context) {
	resp := app.NewResponse(c)

	// 从 token 中解析 uid 出来
	claims := jwt.Get(c).(*model.JWTClaims)

	userSvc := service.NewUserService(c.Request().Context())

	user, err := userSvc.Get(claims.UID)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			resp.ToError(errcode.NotFound.WithMsg("用户不存在"))
			return
		}
		resp.ToError(errcode.GetUserFail.WithDetails(err.Error()))
		return
	}
	resp.ToSuccess(user)
	return
}

// --- U --- //
type UpdateUserReq struct {
	Name  string `json:"name" validate:"omitempty,alphanum,min=6,max=20"`
	Phone string `json:"phone" validate:"omitempty,numeric,len=11"`
}

// 更新用户信息 godoc
// @summary 更新用户信息
// @description 更新用户名、手机号等，字段不传或为空字符串不修改此字段。
// @description 只针对用户基本信息修改，其他信息例如角色、密码等，通过专门的接口改，便于权限控制。
// @accept json
// @produce json
// @tags user
// @param id path int true "用户ID"
// @param name body string false "用户名"
// @param phone body string false "手机号"
// @success 200 {object} model.Resp{data=model.User}
// @router /api/v1/users/{id} [put]
func (u User) Update(c iris.Context) {
	param := UpdateUserReq{}
	resp := app.NewResponse(c)

	// 因为路由直接写的 id:uint，不可能拿到的不是 uint 类型
	id, _ := c.Params().GetUint("id")

	// 校验参数
	if !app.BindAndValid(c, &param) {
		return
	}

	// 解析 token
	claims := jwt.Get(c).(*model.JWTClaims)
	if claims.UID != id {
		// 判断用户是否存在 admin 角色
		isAdmin := utils.IsAdmin(claims.Roles)
		if !isAdmin {
			resp.ToError(errcode.Forbidden.WithMsg("非管理员不允许修改他人用户信息"))
			return
		}
	}

	userSvc := service.NewUserService(c.Request().Context())

	user, cerr := userSvc.Update(id, param.Name, param.Phone)
	if cerr != nil {
		resp.ToError(cerr)
		return
	}

	resp.ToSuccess(user)
}

type UpdateUserPasswordReq struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=20"`
}

// 修改用户密码 godoc
// @summary 修改用户密码
// @description 更新用户密码
// @accept json
// @produce json
// @tags user
// @param id path int true "用户ID"
// @param old_password body string true "原密码"
// @param new_password body string true "新密码"
// @success 200 {object} model.Resp
// @router /api/v1/users/{id}/password [put]
func (u User) UpdatePassword(c iris.Context) {
	param := UpdateUserPasswordReq{}
	resp := app.NewResponse(c)

	id, _ := c.Params().GetUint("id")

	// 校验参数
	if !app.BindAndValid(c, &param) {
		return
	}

	// 解析 token
	claims := jwt.Get(c).(*model.JWTClaims)
	if claims.UID != id {
		// 判断用户是否存在 admin 角色
		isAdmin := utils.IsAdmin(claims.Roles)
		if !isAdmin {
			resp.ToError(errcode.Forbidden.WithMsg("非管理员不允许修改他人密码"))
			return
		}
	}

	// 修改密码
	userSvc := service.NewUserService(c.Request().Context())
	cerr := userSvc.UpdatePassword(id, param.OldPassword, param.NewPassword)
	if cerr != nil {
		resp.ToError(cerr)
		return
	}
	resp.ToSuccess(nil)
}

type UpdateUserRoleReq struct {
	Role model.Role `json:"role" validate:"required,oneof=admin member"`
}

// 修改用户角色 godoc
// @summary 修改用户角色（管理员）
// @description 只允许管理员调用，修改用户角色信息
// @accept json
// @produce json
// @tags 管理员
// @param id path int true "用户ID"
// @param role body string true "角色，member 或 admin" Enums(admin, member)
// @success 200 {object} model.Resp
// @router /api/v1/users/{id}/role [put]
func (u User) UpdateRole(c iris.Context) {
	param := UpdateUserRoleReq{}
	resp := app.NewResponse(c)

	id, _ := c.Params().GetUint("id")

	// 校验参数
	if !app.BindAndValid(c, &param) {
		return
	}

	userSvc := service.NewUserService(c.Request().Context())
	err := userSvc.UpdateRole(id, param.Role)
	if err != nil {
		resp.ToError(errcode.InternalServerError.WithDetails(err.Error()))
		return
	}
	resp.ToSuccess(nil)
}

// --- D --- //

// 删除用户 godoc
// @Summary 删除用户
// @Description 根据 ID 删除用户
// @Accept json
// @Produce json
// @tags user
// @param id path int true "用户ID"
// @success 200 {object} model.Resp
// @router /api/v1/users/{id} [delete]
func (u User) Delete(c iris.Context) {
	resp := app.NewResponse(c)

	id, _ := c.Params().GetUint("id")

	userSvc := service.NewUserService(c.Request().Context())

	// 执行删除
	if err := userSvc.Delete(id); err != nil {
		resp.ToError(err)
		return
	}
	resp.ToSuccess(nil)
	return
}
