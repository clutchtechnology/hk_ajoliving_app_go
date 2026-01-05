package middlewares

import (
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
)

// JWTAuth JWT 认证中间件
func JWTAuth(jwtManager *tools.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			tools.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		// 检查 Bearer 前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			tools.Unauthorized(c, "Invalid authorization format")
			c.Abort()
			return
		}

		// 提取 token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			tools.Unauthorized(c, "Missing token")
			c.Abort()
			return
		}

		// 解析 token
		claims, err := jwtManager.ParseToken(tokenString)
		if err != nil {
			if err == tools.ErrExpiredToken {
				tools.Unauthorized(c, "Token expired")
			} else {
				tools.Unauthorized(c, "Invalid token")
			}
			c.Abort()
			return
		}

		// 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("user_type", claims.UserType)

		c.Next()
	}
}

// OptionalJWTAuth 可选的 JWT 认证中间件（不强制要求认证）
func OptionalJWTAuth(jwtManager *tools.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := jwtManager.ParseToken(tokenString)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("email", claims.Email)
				c.Set("user_type", claims.UserType)
			}
		}

		c.Next()
	}
}
