package v1

import (
	"errors"
	"github.com/go-pg/pg/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/response"
	"github.com/xuxusheng/time-frequency-be/internal/service"
	"github.com/xuxusheng/time-frequency-be/internal/utils"
)

// 普通用户具有权限的接口
type IUser interface {
	Me(c iris.Context)
	IsNameExist(c iris.Context)
	IsPhoneExist(c iris.Context)
	IsEmailExist(c iris.Context)

	// 用户自身只允许改自己的手机号和邮箱
	Update(c iris.Context)

	Login(c iris.Context)
}

type User struct {
	userSvc service.IUser
}

func NewUser(userSvc service.IUser) *User {
	return &User{userSvc: userSvc}
}

func (u *User) Me(c iris.Context) {
	ctx := c.Request().Context()
	resp := response.New(c)
	claims := jwt.Get(c).(*model.JWTClaims)

	user, err := u.userSvc.Get(ctx, claims.Uid)
	if err != nil {
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

func (u User) IsNameExist(c iris.Context) {
	panic("implement me")
}

func (u User) IsPhoneExist(c iris.Context) {
	panic("implement me")
}

func (u User) IsEmailExist(c iris.Context) {
	panic("implement me")
}

// --- U ---

type UserUpdateReq struct {
	Id    int    `json:"id" validate:"required,min=1"`
	Phone string `json:"phone" validate:"required"`
	Email string `json:"email" validate:"required"`
}

func (u User) Update(c iris.Context) {
	ctx := c.Request().Context()
	resp := response.New(c)

	p := UserUpdateReq{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	user, err := u.userSvc.UpdatePhoneAndEmail(ctx, p.Id, p.Phone, p.Email)
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
	resp.Success(iris.Map{
		"id":         user.Id,
		"name":       user.Name,
		"phone":      user.Phone,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// --- AUTH ---
type LoginReq struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u User) Login(c iris.Context) {
	ctx := c.Request().Context()
	resp := response.New(c)

	p := LoginReq{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	svc := u.userSvc

	user, err := svc.GetByName(ctx, p.Name)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			resp.Error(cerror.BadRequest.WithMsg("用户名或密码错误"))
			return
		}
		resp.Error(cerror.BadRequest.WithDebugs(err))
		return
	}

	// 用户存在，对比密码
	err = utils.ComparePwd(user.Password, p.Password)
	if err != nil {
		resp.Error(cerror.BadRequest.WithMsg("用户名或密码错误"))
		return
	}

	// 密码正确，生成 token 并返回
	token, err := jwt.Sign(
		jwt.HS256,
		[]byte(global.Setting.JWT.Secret),
		model.JWTClaims{Uid: user.Id},
		jwt.MaxAge(global.Setting.JWT.Expire),
	)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	resp.Success(iris.Map{
		"token": string(token),
		"user": iris.Map{
			"id":   user.Id,
			"name": user.Name,
		},
	})
}
