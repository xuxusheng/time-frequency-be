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

// 老师角色才能调用的接口
type ITeacher interface {
	CreateUser(c iris.Context) // 创建用户

	ListUser(c iris.Context) // 查询用户列表

	UpdateUser(c iris.Context) // 修改用户信息，老师修改时，允许修改用户名和昵称

	DeleteUser(c iris.Context) // 删除用户
}

type Teacher struct {
	userSvc service.IUser
}

func (t *Teacher) CreateUser(c iris.Context) {
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

	user, err := t.userSvc.Create(ctx, claims.Uid, p.Name, p.NickName, p.Phone, p.Email, p.Password)
	if err != nil {
		if cerr, ok := err.(cerror.IError); ok {
			resp.Error(cerr)
			return
		}
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}

	resp.Success(iris.Map{
		"id":         user.Id,
		"name":       user.Name,
		"nick_name":  user.NickName,
		"phone":      user.Phone,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

func (t *Teacher) ListUser(c iris.Context) {
	p := struct {
		Query string `json:"query"`
		Pn    int    `json:"pn"`
		Ps    int    `json:"ps"`
	}{}
	ctx := c.Request().Context()
	resp := response.New(c)
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	page := model.NewPage(p.Pn, p.Ps)

	users, count, err := t.userSvc.ListAndCount(ctx, p.Query, page)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}

	page.WithTotal(count)
	data := []iris.Map{}
	for _, user := range users {
		data = append(data, iris.Map{
			"id":         user.Id,
			"name":       user.Name,
			"nick_name":  user.NickName,
			"phone":      user.Phone,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}
	resp.SuccessList(data, page)
}

func (t *Teacher) UpdateUser(c iris.Context) {
	panic("implement me")
}

func (t *Teacher) DeleteUser(c iris.Context) {
	p := struct {
		Id int `json:"id" validate:"required,min=1"`
	}{}
	ctx := c.Request().Context()
	resp := response.New(c)
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	err := t.userSvc.Delete(ctx, p.Id)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	resp.Success()
}

func NewTeacher(userSvc service.IUser) *Teacher {
	return &Teacher{userSvc: userSvc}
}
