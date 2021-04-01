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

	// 用户自身只允许改自己的手机号和邮箱
	Update(c iris.Context)
	UpdatePassword(c iris.Context)

	Login(c iris.Context)
}

type User struct {
	userSvc service.IUser
}

func NewUser(userSvc service.IUser) *User {
	return &User{userSvc: userSvc}
}

// --- R ---
func (u *User) Me(c iris.Context) {
	ctx := c.Request().Context()
	resp := response.New(c)
	claims := jwt.Get(c).(*model.JWTClaims)

	user, err := u.userSvc.Get(ctx, claims.Uid)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}

	resp.Success(user)
}

// --- U ---
func (u User) Update(c iris.Context) {
	p := struct {
		Phone string `json:"phone" validate:"required"`
		Email string `json:"email" validate:"required"`
	}{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)
	claims := jwt.Get(c).(*model.JWTClaims)

	user := model.User{
		Id:    claims.Uid,
		Phone: p.Phone,
		Email: p.Email,
	}
	err := u.userSvc.Update(ctx, &user, []string{"phone", "email"})
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

func (u *User) UpdatePassword(c iris.Context) {
	p := struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required"`
	}{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)
	claims := jwt.Get(c).(*model.JWTClaims)

	user, err := u.userSvc.UpdatePassword(ctx, claims.Uid, p.OldPassword, p.NewPassword)
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

// --- AUTH ---
func (u User) Login(c iris.Context) {
	p := struct {
		Name     string `json:"name" validate:"required"`
		Password string `json:"password" validate:"required"`
	}{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	ctx := c.Request().Context()
	resp := response.New(c)

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
		"user":  user,
	})
}
