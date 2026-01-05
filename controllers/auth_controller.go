package controllers

// AuthHandler Methods:
// 0. NewAuthHandler(authService *services.AuthService) -> 注入 AuthService
// 1. Register(c *gin.Context) -> 用户注册
// 2. Login(c *gin.Context) -> 用户登录
// 3. Logout(c *gin.Context) -> 用户登出
// 4. RefreshToken(c *gin.Context) -> 刷新令牌
// 5. ForgotPassword(c *gin.Context) -> 忘记密码
// 6. ResetPassword(c *gin.Context) -> 重置密码
// 7. VerifyCode(c *gin.Context) -> 验证码验证

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
)

// AuthHandlerInterface 认证处理器接口
type AuthHandlerInterface interface {
	Register(c *gin.Context)       // 1. 用户注册
	Login(c *gin.Context)          // 2. 用户登录
	Logout(c *gin.Context)         // 3. 用户登出
	RefreshToken(c *gin.Context)   // 4. 刷新令牌
	ForgotPassword(c *gin.Context) // 5. 忘记密码
	ResetPassword(c *gin.Context)  // 6. 重置密码
	VerifyCode(c *gin.Context)     // 7. 验证码验证
}

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *services.AuthService
}

// 0. NewAuthHandler 注入 AuthService
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// 1. Register 用户注册
// Register godoc
// @Summary      用户注册
// @Description  创建新用户账号
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RegisterRequest  true  "注册信息"
// @Success      200   {object}  models.Response{data=models.RegisterResponse}
// @Failure      400   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if err == services.ErrUserAlreadyExists {
			tools.BadRequest(c, "User already exists")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 2. Login 用户登录
// Login godoc
// @Summary      用户登录
// @Description  使用邮箱/用户名和密码登录
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.LoginRequest  true  "登录信息"
// @Success      200   {object}  models.Response{data=models.AuthResponse}
// @Failure      400   {object}  models.Response
// @Failure      401   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			tools.Unauthorized(c, "Invalid credentials")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 3. Logout 用户登出
// Logout godoc
// @Summary      用户登出
// @Description  退出当前登录会话
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.Response{data=models.LogoutResponse}
// @Failure      401  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// 从上下文获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "Unauthorized")
		return
	}

	err := h.authService.Logout(c.Request.Context(), userID.(uint))
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, gin.H{
		"message": "Logged out successfully",
	})
}

// 4. RefreshToken 刷新令牌
// RefreshToken godoc
// @Summary      刷新令牌
// @Description  使用刷新令牌获取新的访问令牌
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.RefreshTokenRequest  true  "刷新令牌"
// @Success      200   {object}  models.Response{data=models.AuthResponse}
// @Failure      400   {object}  models.Response
// @Failure      401   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		if err == services.ErrInvalidToken {
			tools.Unauthorized(c, "Invalid refresh token")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 5. ForgotPassword 忘记密码
// ForgotPassword godoc
// @Summary      忘记密码
// @Description  发送密码重置邮件
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.ForgotPasswordRequest  true  "邮箱地址"
// @Success      200   {object}  models.Response{data=models.ForgotPasswordResponse}
// @Failure      400   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.ForgotPassword(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 6. ResetPassword 重置密码
// ResetPassword godoc
// @Summary      重置密码
// @Description  使用重置令牌设置新密码
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.ResetPasswordRequest  true  "重置密码信息"
// @Success      200   {object}  models.Response{data=models.ResetPasswordResponse}
// @Failure      400   {object}  models.Response
// @Failure      401   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		if err == services.ErrInvalidToken {
			tools.Unauthorized(c, "Invalid or expired reset token")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 7. VerifyCode 验证码验证
// VerifyCode godoc
// @Summary      验证码验证
// @Description  验证邮箱验证码
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      models.VerifyCodeRequest  true  "验证码信息"
// @Success      200   {object}  models.Response{data=models.VerifyCodeResponse}
// @Failure      400   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/auth/verify-code [post]
func (h *AuthHandler) VerifyCode(c *gin.Context) {
	var req models.VerifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := h.authService.VerifyCode(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}
