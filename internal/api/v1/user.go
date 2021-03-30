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

type IUser interface {
	Create(c iris.Context)
	Get(c iris.Context)
	List(c iris.Context)
	Update(c iris.Context)
	Delete(c iris.Context)
	IsNameExist(c iris.Context)
	IsPhoneExist(c iris.Context)
	IsEmailExist(c iris.Context)
	Login(c iris.Context)
}

type User struct {
	userSvc service.IUser
}

func NewUser(userSvc service.IUser) *User {
	return &User{userSvc: userSvc}
}

// --- C ---
type CreateUserReq struct {
	Name     string `json:"name" validate:"required,min=1"`
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (u User) Create(c iris.Context) {
	ctx := c.Request().Context()
	resp := response.New(c)

	p := CreateUserReq{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	// 计算密码 hash
	hash, err := utils.EncodePwd(p.Password)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}

	user, err := u.userSvc.Create(
		ctx,
		1,
		p.Name,
		p.Phone,
		p.Email,
		hash,
	)

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
type GetUserReq struct {
	Id int `json:"id" validate:"required,min=1"`
}

func (u User) Get(c iris.Context) {
	ctx := c.Request().Context()
	resp := response.New(c)

	p := GetUserReq{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	user, err := u.userSvc.Get(ctx, p.Id)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			resp.Error(cerror.NotFound.WithMsg("用户不存在"))
			return
		}

		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}

	data := iris.Map{
		"id":         user.Id,
		"name":       user.Name,
		"phone":      user.Phone,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"created_by": iris.Map{},
	}

	// 第一个用户，可能没有创建人
	if user.CreatedBy != nil {
		data["created_by"] = iris.Map{
			"id":    user.CreatedBy.Id,
			"name":  user.CreatedBy.Name,
			"phone": user.CreatedBy.Phone,
			"email": user.CreatedBy.Email,
		}
	}

	resp.Success(data)
}

type ListUserReq struct {
	Query string `json:"query"`
	Pn    int    `json:"pn"`
	Ps    int    `json:"ps"`
}

func (u User) List(c iris.Context) {
	ctx := c.Request().Context()
	resp := response.New(c)

	p := ListUserReq{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	page := model.NewPage(p.Pn, p.Ps)

	users, count, err := u.userSvc.ListAndCount(ctx, p.Query, page)
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
			"phone":      user.Phone,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		})
	}
	resp.SuccessList(data, page)
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

type UpdateUserReq struct {
	Id    int    `json:"id" validate:"required,min=1"`
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required"`
	Email string `json:"email" validate:"required"`
}

func (u User) Update(c iris.Context) {
	ctx := c.Request().Context()
	resp := response.New(c)

	p := UpdateUserReq{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	user, err := u.userSvc.Update(ctx, p.Id, p.Name, p.Phone, p.Email)
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

// --- D ---
type DeleteUserReq struct {
	Id int `json:"id" validate:"required,min=1"`
}

func (u User) Delete(c iris.Context) {
	ctx := c.Request().Context()
	resp := response.New(c)

	p := DeleteUserReq{}
	if ok := utils.BindAndValidate(c, &p); !ok {
		return
	}

	err := u.userSvc.Delete(ctx, p.Id)
	if err != nil {
		resp.Error(cerror.ServerError.WithDebugs(err))
		return
	}
	resp.Success()
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
