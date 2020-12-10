package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/xuxusheng/time-frequency-be/global"
	"github.com/xuxusheng/time-frequency-be/pkg/logger"
	"time"
)

// gin 的上下文无法获取 body 内容，在 gin.ResponseWriter 上再包一层，用来存 body 内容
type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	// 写 response 内容的时候，先往 body 里面写，方便在中间件中取出来用
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyWriter := &AccessLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = bodyWriter

		beginTime := time.Now().Unix()
		c.Next()
		endTime := time.Now().Unix()

		fields := logger.Fields{
			"url":      c.Request.URL.RequestURI(),
			"request":  c.Request.PostForm.Encode(),
			"response": bodyWriter.body.String(),
		}

		s := "access log: method: %s, status_code: %d, " + "begin_time: %d, end_time: %d"
		global.Logger.WithFields(fields).Infof(c, s, c.Request.Method, bodyWriter.Status(), beginTime, endTime)
	}

}
