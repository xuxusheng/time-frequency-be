package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/pkg/app"
	"github.com/xuxusheng/time-frequency-be/pkg/errcode"
)

func IsLogin() iris.Handler {

	verifier := jwt.NewVerifier(jwt.HS256, global.JWTSetting.Secret)

	verifier.ErrorHandler = func(c iris.Context, err error) {
		// 判断错误类型
		if err != nil {
			resp := app.NewResponse(c)
			switch err {
			case jwt.ErrMissing:
				resp.ToError(errcode.UnauthorizedTokenEmpty)
			case jwt.ErrExpired:
				resp.ToError(errcode.UnauthorizedTokenExpired)
			default:
				resp.ToError(errcode.UnauthorizedTokenError)
			}
		}

	}

	return verifier.Verify(func() interface{} {
		return new(model.JWTClaims)
	})

}
