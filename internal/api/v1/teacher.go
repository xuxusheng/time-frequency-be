package v1

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/response"
	"github.com/xuxusheng/time-frequency-be/internal/service"
	"github.com/xuxusheng/time-frequency-be/internal/utils"
)

// 老师角色才能调用的接口
type ITeacher interface {
	CreateStudent(c iris.Context) // 创建用户

	ListStudent(c iris.Context) // 查询学生列表
	IsNameExist(c iris.Context)
	IsPhoneExist(c iris.Context)
	IsEmailExist(c iris.Context)

	UpdateUser(c iris.Context) // 修改用户信息，老师修改时，允许修改用户名和昵称

	DeleteStudent(c iris.Context) // 删除用户
}

type Teacher struct {
	userSvc service.IUser
}

func NewTeacher(userSvc service.IUser) *Teacher {
	return &Teacher{userSvc: userSvc}
}

// --- C ---
func (t *Teacher) CreateStudent(c iris.Context) {
	p := struct {
		Name     string `json:"name" validate:"required"`
		NickName string `json:"nick_name" validate:"required"`
		Phone    string `json:"phone" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}{}
	ctx := c.Request().Context()
	resp := response.New(c)
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	claims := jwt.Get(c).(*model.JWTClaims)

	user := model.User{
		CreatedById: claims.Uid,
		Name:        p.Name,
		NickName:    p.NickName,
		Phone:       p.Phone,
		Email:       p.Email,
		Role:        model.UserRoleStudent,
		IsAdmin:     false,
		Password:    p.Password,
	}
	err := t.userSvc.Create(ctx, &user)
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

// --- R ---
func (t *Teacher) ListStudent(c iris.Context) {
	p := struct {
		Query string `json:"query"`
		Pn    int    `json:"pn"`
		Ps    int    `json:"ps"`
	}{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)
	page := model.NewPage(p.Pn, p.Ps)

	users, count, err := t.userSvc.ListAndCount(ctx, page, p.Query, "student")
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}

	page.WithTotal(count)
	resp.SuccessList(users, page)
}

func (t *Teacher) IsNameExist(c iris.Context) {
	p := struct {
		Name      string `json:"name" validate:"required, min=1"`
		ExcludeId int    `json:"exclude_id"`
	}{}

	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)

	is, err := t.userSvc.IsNameExist(ctx, p.Name, p.ExcludeId)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	resp.Success(is)
}

func (t *Teacher) IsPhoneExist(c iris.Context) {
	p := struct {
		Phone     string `json:"phone" validate:"required, min=1"`
		ExcludeId int    `json:"exclude_id"`
	}{}

	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)

	is, err := t.userSvc.IsPhoneExist(ctx, p.Phone, p.ExcludeId)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	resp.Success(is)
}

func (t *Teacher) IsEmailExist(c iris.Context) {
	p := struct {
		Email     string `json:"email" validate:"required, min=1"`
		ExcludeId int    `json:"exclude_id"`
	}{}

	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)

	is, err := t.userSvc.IsEmailExist(ctx, p.Email, p.ExcludeId)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	resp.Success(is)
}

// --- U ---
func (t *Teacher) UpdateUser(c iris.Context) {
	p := struct {
		Id       int    `json:"id" validate:"required"`
		NickName string `json:"nick_name" validate:"required"`
		Phone    string `json:"phone" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password"`
	}{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)
	columns := []string{"nick_name", "phone", "email"}

	// Password 字段如果为空的话，就不修改
	if p.Password != "" {
		columns = append(columns, "password")
	}

	user := model.User{
		Id:       p.Id,
		NickName: p.NickName,
		Phone:    p.Phone,
		Email:    p.Email,
		Password: p.Password,
	}
	err := t.userSvc.Update(ctx, &user, columns)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			resp.Error(cerror.NotFound.WithMsg("用户不存在"))
			return
		}
		if cerr, ok := err.(cerror.IError); ok {
			resp.Error(cerr)
			return
		}
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}

	resp.Success(user)
}

// --- D ---
func (t *Teacher) DeleteStudent(c iris.Context) {
	p := struct {
		Id int `json:"id" validate:"required,min=1"`
	}{}
	ctx := c.Request().Context()
	resp := response.New(c)
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	user, err := t.userSvc.Get(ctx, p.Id)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			resp.Error(cerror.NotFound.WithMsg("用户不存在"))
			return
		}
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}

	// 验证被删除的用户是否是学生
	if user.Role != "student" {
		resp.Error(cerror.Forbidden.WithMsg("无权删除非学生账号，请联系管理员"))
		return
	}

	err = t.userSvc.Delete(ctx, p.Id)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	resp.Success()
}
