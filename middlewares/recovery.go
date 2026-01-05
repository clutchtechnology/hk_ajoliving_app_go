package middlewares

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
)

// Recovery 异常恢复中间件
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, tools.Response{
					Code:    tools.CodeInternalError,
					Message: "Internal server error",
				})
			}
		}()

		c.Next()
	}
}
