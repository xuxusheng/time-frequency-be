package v1

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/response"
	"github.com/xuxusheng/time-frequency-be/internal/service"
	"github.com/xuxusheng/time-frequency-be/internal/utils"
)

// 管理员才能调用的接口
type IAdmin interface {
	CreateUser(c iris.Context) // 创建新用户

	GetUser(c iris.Context)  // 查询单个用户详细信息
	ListUser(c iris.Context) // 查询所有用户

	UpdateUser(c iris.Context)  // 修改用户信息
	ToggleAdmin(c iris.Context) // 切换用户是否为管理员

	DeleteUser(c iris.Context) // 删除用户
}

type Admin struct {
	userSvc service.IUser
}

func NewAdmin(userSvc service.IUser) *Admin {
	return &Admin{userSvc: userSvc}
}

// 创建新用户 godoc
// @summary 创建新用户
// @description 管理员创建新用户账号，可以指定用户角色和是否为管理员账户
// @accept json
// @produce json
// @tags admin
// @param name body string true "用户名，建议使用姓名拼音"
// @param nick_name body string true "用户昵称，请使用真实姓名"
// @param phone body string true "手机号"
// @param email body string true "邮箱"
// @param password body string true "密码"
// @param role body string true "用户角色，student 学生或 teacher 老师"  Enums(student, teacher)
// @param is_admin body bool true "是否是管理员"
// @success 200 {object} swagger.Resp{data=model.User}
// @router /api/v1/admin/create-user [post]
func (a Admin) CreateUser(c iris.Context) {
	p := struct {
		Name     string `json:"name" validate:"required"`
		NickName string `json:"nick_name" validate:"required"`
		Phone    string `json:"phone" validate:"required"`
		Email    string `json:"email" validated:"required"`
		Password string `json:"password" validated:"required"`
		Role     string `json:"role" validated:"required,oneof=student teacher"`
		IsAdmin  bool   `json:"is_admin"`
	}{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)
	claims := jwt.Get(c).(*model.JWTClaims)

	user := model.User{
		CreatedById: claims.Uid,
		Name:        p.Name,
		NickName:    p.NickName,
		Phone:       p.Phone,
		Email:       p.Email,
		Password:    p.Password,
		Role:        p.Role,
		IsAdmin:     p.IsAdmin,
	}
	err := a.userSvc.Create(ctx, &user)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	resp.Success(user)
}

// 查询单个用户详情 godoc
// @summary 查询单个用户详情
// @description 查询单个用户详情（供管理员使用）
func (a Admin) GetUser(c iris.Context) {
	p := struct {
		Id int `json:"id" validate:"required"`
	}{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)

	user, err := a.userSvc.Get(ctx, p.Id)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	resp.Success(user)
}

// 查询多个用户 godoc
// @summary 查询多个用户
// @description 管理员查询所有的用户账号
// @accept json
// @produce json
// @tags admin
// @param query body string true "模糊匹配用户名、昵称、手机号和邮箱"
// @param role body string true "通过角色筛选老师或者学生"  Enums(student, teacher)
// @param is_admin body int true "筛选是否是管理员，-1 不限、0 否、1 是" Enums(-1, 0, 1)
// @param pn body int true "pn"
// @param ps body int true "ps"
// @success 200 {object} swagger.Resp{data=swagger.DWithP{data=model.User}}
// @router /api/v1/admin/list-user [post]
func (a Admin) ListUser(c iris.Context) {
	p := struct {
		Query string `json:"query" validated:"required"`
		Role  string `json:"role" validated:"required"`
		Pn    int    `json:"pn" validated:"required"`
		Ps    int    `json:"ps" validated:"required"`
	}{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)
	page := model.NewPage(p.Pn, p.Ps)

	users, count, err := a.userSvc.ListAndCount(ctx, page, p.Query, p.Role)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	page.WithTotal(count)
	resp.SuccessList(users, page)
}

// 修改用户信息 godoc
// @summary 修改用户信息
// @description 修改某个用户的信息（此接口不允许修改用户是否为管理员）
// @accept json
// @produce json
// @tags admin
// @param id body int true "用户ID"
// @param name body string true "用户名"
// @param nick_name body string true "用户昵称"
// @param phone body string true "手机号"
// @param email body string true "邮箱"
// @param role body string true "用户角色" Enums(student, teacher)
// @param password body string false "用户密码，留空则不修改"
// @success 200 {object} swagger.Resp{data=model.User}
// @router /api/v1/admin/update-user [post]
func (a Admin) UpdateUser(c iris.Context) {
	p := struct {
		Id       int    `json:"id" validated:"required"`
		Name     string `json:"name" validate:"required"`
		NickName string `json:"nick_name" validate:"required"`
		Phone    string `json:"phone" validate:"required"`
		Email    string `json:"email" validated:"required"`
		Role     string `json:"role" validated:"required,oneof=student teacher"`
		Password string `json:"password"`
	}{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)

	user := model.User{
		Id:       p.Id,
		Name:     p.Name,
		NickName: p.NickName,
		Phone:    p.Phone,
		Email:    p.Email,
		Password: p.Password,
		Role:     p.Role,
	}
	columns := []string{"name", "nick_name", "phone", "email", "role"}
	// 如果参数中没有 password 或者 password 为空的话，就不修改用户密码
	if user.Password != "" {
		columns = append(columns, "password")
	}
	err := a.userSvc.Update(ctx, &user, columns)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	resp.Success(user)
}

// 修改账号管理员权限 godoc
// @summary 修改账号管理员权限
// @description 修改某个用户是否是管理员，不允许取消自己的管理员权限
// @accept json
// @produce json
// @tags admin
// @param id body int true "用户ID"
// @param is_admin body bool true "是否管理员"
// @success 200 {object} swagger.Resp
// @router /api/v1/admin/toggle-admin [post]
func (a Admin) ToggleAdmin(c iris.Context) {
	p := struct {
		Id      int  `json:"id" validated:"required"`
		IsAdmin bool `json:"is_admin" validated:"required"`
	}{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)

	claims := jwt.Get(c).(*model.JWTClaims)

	if claims.Uid == p.Id && !p.IsAdmin {
		resp.Error(cerror.BadRequest.WithMsg("无法取消自己的管理员权限"))
		return
	}

	user := model.User{
		Id:      p.Id,
		IsAdmin: p.IsAdmin,
	}
	err := a.userSvc.Update(ctx, &user, []string{"is_admin"})
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	resp.Success()
}

// 删除账号 godoc
// @summary 删除账号
// @description 管理员删除某个账号
// @accept json
// @produce json
// @tags admin
// @param id body int true "用户ID"
// @success 200 {object} swagger.Resp
// @router /api/v1/admin/delete-user [post]
func (a Admin) DeleteUser(c iris.Context) {
	p := struct {
		Id int `json:"id" validated:"required"`
	}{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)

	claims := jwt.Get(c).(*model.JWTClaims)

	if claims.Uid == p.Id {
		resp.Error(cerror.BadRequest.WithMsg("无法删除自己的账号"))
		return
	}

	err := a.userSvc.Delete(ctx, p.Id)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	resp.Success()
}
