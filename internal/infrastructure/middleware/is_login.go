package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/jwt"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/internal/model"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/cerror"
	"github.com/xuxusheng/time-frequency-be/internal/pkg/response"
)

func IsLogin() iris.Handler {
	verifier := jwt.NewVerifier(jwt.HS256, global.Setting.JWT.Secret)

	verifier.ErrorHandler = func(c iris.Context, err error) {
		// 判断错误类型
		if err != nil {
			resp := response.New(c)
			switch err {
			case jwt.ErrMissing:
				resp.Error(cerror.TokenEmpty)
			case jwt.ErrExpired:
				resp.Error(cerror.TokenExpired)
			default:
				resp.Error(cerror.TokenInvalid)
			}
		}

	}

	return verifier.Verify(func() interface{} {
		return new(model.JWTClaims)
	})
}
