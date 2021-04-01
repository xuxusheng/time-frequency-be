package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/dao"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/response"
)

func IsTeacher() iris.Handler {
	return func(c iris.Context) {
		ctx := c.Request().Context()
		resp := response.New(c)

		claims := jwt.Get(c).(*model.JWTClaims)
		// 验证当前用户是否为管理员或老师
		d := dao.NewUser(global.DB)

		user, err := d.Get(ctx, claims.Uid)
		if err != nil {
			resp.Error(cerror.ServerError.WithDebugs(err))
			return
		}
		if !user.IsAdmin && user.Role != "teacher" {
			resp.Error(cerror.Forbidden.WithMsg("非管理员或老师无权访问"))
			return
		}
		c.Next()
	}
}
