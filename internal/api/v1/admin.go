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

	ListUser(c iris.Context) // 查询所有用户

	ToggleAdmin(c iris.Context) // 切换用户是否为管理员
}

type Admin struct {
	userSvc service.IUser
}

func NewAdmin(userSvc service.IUser) *Admin {
	return &Admin{userSvc: userSvc}
}

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

func (a Admin) ListUser(c iris.Context) {
	p := struct {
		Query string `json:"query"`
		Role  string `json:"role"`
		Pn    int    `json:"pn"`
		Ps    int    `json:"ps"`
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

func (a Admin) ToggleAdmin(c iris.Context) {
	p := struct {
		Id      int  `json:"id"`
		IsAdmin bool `json:"is_admin"`
	}{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)

	claims := jwt.Get(c).(*model.JWTClaims)

	if claims.Uid == p.Id && !p.IsAdmin {
		resp.Error(cerror.BadRequest.WithMsg("无法取消自己的管理员权限"))
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
