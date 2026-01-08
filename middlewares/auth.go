package middlewares

import (
	"strings"

	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			tools.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		// 检查 Bearer 前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			tools.Unauthorized(c, "invalid token format")
			c.Abort()
			return
		}

		// 提取 token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析 token
		claims, err := tools.ParseToken(tokenString)
		if err != nil {
			tools.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("user_type", claims.UserType)

		c.Next()
	}
}
