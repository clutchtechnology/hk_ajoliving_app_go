package controllers

import (
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController 创建认证控制器
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register 用户注册
func (ctrl *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	user, err := ctrl.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "email already exists" {
			tools.BadRequest(c, "email already exists")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Created(c, user.ToUserResponse())
}

// Login 用户登录
func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	loginResp, err := ctrl.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "invalid email or password" {
			tools.Unauthorized(c, "invalid email or password")
			return
		}
		if err.Error() == "account is not active" {
			tools.Forbidden(c, "account is not active")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, loginResp)
}

// Logout 用户登出
func (ctrl *AuthController) Logout(c *gin.Context) {
	// JWT 是无状态的，客户端只需删除 token 即可
	// 如需实现黑名单功能，可在此处添加 token 到 Redis 黑名单
	tools.Success(c, gin.H{
		"message": "logged out successfully",
	})
}
