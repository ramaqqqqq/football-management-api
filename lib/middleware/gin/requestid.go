package ginmiddleware

import (
	"context"
	"fmt"
	"time"

	"go-test/lib/middleware/request"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-Id")
		if reqID == "" {
			uid, err := uuid.NewRandom()
			if err != nil {
				reqID = fmt.Sprintf("%d", time.Now().UnixNano())
			} else {
				reqID = uid.String()
			}
		}

		c.Set(string("request_id"), reqID)

		ctx := context.WithValue(c.Request.Context(), request.CtxKeyReqId, reqID)
		c.Request = c.Request.WithContext(ctx)

		c.Header("X-Request-Id", reqID)
		c.Next()
	}
}

func GetRequestID(c *gin.Context) string {
	if v, exists := c.Get("request_id"); exists {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
