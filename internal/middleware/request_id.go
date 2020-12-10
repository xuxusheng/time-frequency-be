package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 给请求加上 request_id
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 先判断下 request Header 中是否已存在 x-request-id 字段
		reqID := c.Request.Header.Get("x-request-id")

		if reqID == "" {
			// 请求中未包含 reqID，自己生成一个
			uuid, _ := uuid.NewRandom()
			reqID = uuid.String()
		}

		c.Set("x-request-id", reqID)
		// 将 reqID 写入 response Header 中
		c.Header("x-request-id", reqID)

		c.Next()
	}
}
