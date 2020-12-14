package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/service"
	"github.com/xuxusheng/time-frequency-be/pkg/app"
	"github.com/xuxusheng/time-frequency-be/pkg/convert"
	"github.com/xuxusheng/time-frequency-be/pkg/errcode"
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
// 创建用户 godoc
// @summary 创建新用户
// @Description 创建新用户接口，专供管理平台调用
// @Accept json
// @Produce json
// @Tags user
// @Param name body string true "用户名（6-20位数字或字母构成)" minlength(6) maxlength(20)
// @Param phone body string true "手机号（十一位数字）" minlength(11) maxlength(11)
// @Param password body string true "密码"
// @ToSuccess 200 {object} model.Resp{data=model.DWithP{data=model.User}}
// @Failure 500 {object} model.ErrResp "内部错误"
// @Router /api/v1/users [post]
func (u User) Create(c *gin.Context) {
	param := CreateUserReq{}
	resp := app.NewResponse(c)

	// 校验参数
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		// 参数校验失败
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		resp.ToError(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	userSvc := service.NewUserService(c.Request.Context())

	// service 返回的 err，已经是自定义的 errcode.ToError 类型了
	_, err := userSvc.Create(param.Name, param.Phone, param.Password)
	if err != nil {
		global.Logger.Error(c, "userSvc.Create err: %v", err)
		resp.ToError(err)
		return
	}

	resp.ToSuccess(nil)
}

// 删除用户 godoc
// @Summary 删除用户
// @Description 根据 ID 删除用户
// @Accept json
// @Produce json
// @tags user
// @param id path int true "用户ID"
// @success 200 {object} model.Resp
// @router /api/v1/users/{id} [delete]
func (u User) Delete(c *gin.Context) {
	resp := app.NewResponse(c)

	// 检查 id 格式
	id, err := convert.StrTo(c.Param("id")).Int()
	if err != nil {
		resp.ToError(errcode.InvalidParams.WithMsg("id 格式错误"))
		return
	}

	userSvc := service.NewUserService(c.Request.Context())

	// 执行删除
	if err := userSvc.Delete(id); err != nil {
		resp.ToError(err)
		return
	}
	resp.ToSuccess(nil)
	return
}

type UpdateUserReq struct {
	Name  string `form:"name" binding:"required,min=4,max=20"`
	Phone string `form:"phone" binding:"required"`
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
// @success 200 {object} model.Resp
// @router /api/v1/users/{id} [put]
func (u User) Update(c *gin.Context) {
	param := UpdateUserReq{}
	resp := app.NewResponse(c)

	// 检查 id 格式
	id, err := convert.StrTo(c.Param("id")).Int()
	if err != nil {
		resp.ToError(errcode.InvalidParams.WithMsg("id 格式错误"))
		return
	}

	// 校验参数
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		// 参数校验失败
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		resp.ToError(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	userSvc := service.NewUserService(c.Request.Context())

	if err := userSvc.Update(id, param.Name, param.Phone); err != nil {
		resp.ToError(err)
		return
	}

	resp.ToSuccess(nil)
}

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
func (u User) List(c *gin.Context) {
	resp := app.NewResponse(c)
	param := UserListReq{}

	// 这里的 errs 类型不是 errcode 包中定义的 ToError，而是 app 包中定义的 ValidError
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		// 参数校验失败
		global.Logger.Errorf(c, "app.BindAndValid errs: %v", errs)
		resp.ToError(
			errcode.InvalidParams.WithDetails(errs.Errors()...),
		)
		return
	}

	pn := app.GetPn(c)
	ps := app.GetPs(c)
	userSvc := service.NewUserService(c.Request.Context())

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
func (u User) Get(c *gin.Context) {

	resp := app.NewResponse(c)
	idStr := c.Param("id")

	id, err := convert.StrTo(idStr).Int()
	if err != nil {
		resp.ToError(
			errcode.InvalidParams.WithMsg(
				"id 格式错误",
			),
		)
		return
	}

	userSvc := service.NewUserService(c.Request.Context())

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
