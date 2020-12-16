package middleware

import (
	"bytes"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/pkg/logger"
	"time"
)

// gin 的上下文无法获取 body 内容，在 gin.ResponseWriter 上再包一层，用来存 body 内容
type AccessLogWriter struct {
	context.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	// 写 response 内容的时候，先往 body 里面写，方便在中间件中取出来用
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func AccessLog() iris.Handler {
	return func(c iris.Context) {
		bodyWriter := &AccessLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.ResponseWriter(),
		}

		c.ResetResponseWriter(bodyWriter)

		beginTime := time.Now().Unix()
		c.Next()
		endTime := time.Now().Unix()

		fields := logger.Fields{
			"url":      c.Request().RequestURI,
			"request":  c.FormValues(),
			"response": bodyWriter.body.String(),
		}

		s := "【访问日志】 method: %s, status_code: %d, " + "begin_time: %d, end_time: %d"
		global.Logger.WithFields(fields).Infof(
			c,
			s,
			c.Method(),
			bodyWriter.StatusCode(),
			beginTime,
			endTime,
		)
	}

}
