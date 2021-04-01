package v1

import (
	"github.com/kataras/iris/v12"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/response"
	"github.com/xuxusheng/time-frequency-be/internal/service"
	"github.com/xuxusheng/time-frequency-be/internal/utils"
)

// 管理员才能调用的接口
type IAdmin interface {
	ListUser(c iris.Context) // 查询所有用户

	ToggleAdmin(c iris.Context) // 切换用户是否为管理员
}

type Admin struct {
	userSvc service.IUser
}

func NewAdmin(userSvc service.IUser) *Admin {
	return &Admin{userSvc: userSvc}
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
	panic("implement me")
}
