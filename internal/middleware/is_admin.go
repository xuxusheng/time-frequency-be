package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/utils"
	"github.com/xuxusheng/time-frequency-be/pkg/app"
	"github.com/xuxusheng/time-frequency-be/pkg/errcode"
)

func IsAdmin() iris.Handler {
	return func(c iris.Context) {
		// 判断是否是 admin
		claims := jwt.Get(c).(*model.JWTClaims)
		isAdmin := utils.IsAdmin(claims.Roles)
		if !isAdmin {
			resp := app.NewResponse(c)
			resp.ToError(errcode.Forbidden.WithMsg("非管理员禁止访问此接口"))
			return
		}
		c.Next()
	}
}
