package controllers

// UserHandler Methods:
// 0. NewUserHandler(userService *service.UserService) -> 注入 UserService
// 1. GetCurrentUser(c *gin.Context) -> 获取当前用户信息
// 2. UpdateCurrentUser(c *gin.Context) -> 更新当前用户信息
// 3. ChangePassword(c *gin.Context) -> 修改密码
// 4. GetMyListings(c *gin.Context) -> 获取我的发布
// 5. UpdateSettings(c *gin.Context) -> 更新设置

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
)

// UserHandlerInterface 用户处理器接口
type UserHandlerInterface interface {
	GetCurrentUser(c *gin.Context)    // 1. 获取当前用户信息
	UpdateCurrentUser(c *gin.Context) // 2. 更新当前用户信息
	ChangePassword(c *gin.Context)    // 3. 修改密码
	GetMyListings(c *gin.Context)     // 4. 获取我的发布
	UpdateSettings(c *gin.Context)    // 5. 更新设置
}

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
}

// 0. NewUserHandler 注入 UserService
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// 1. GetCurrentUser 获取当前用户信息
// GetCurrentUser godoc
// @Summary      获取当前用户信息
// @Description  获取当前登录用户的详细信息
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.Response{data=models.UserResponse}
// @Failure      401  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/v1/users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		models.Unauthorized(c, "Unauthorized")
		return
	}

	result, err := h.userService.GetCurrentUser(c.Request.Context(), userID.(uint))
	if err != nil {
		if err == service.ErrUserNotFound {
			models.NotFound(c, "User not found")
			return
		}
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, result)
}

// 2. UpdateCurrentUser 更新当前用户信息
// UpdateCurrentUser godoc
// @Summary      更新当前用户信息
// @Description  更新当前登录用户的个人信息
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      models.UpdateUserRequest  true  "用户信息"
// @Success      200   {object}  models.Response{data=models.UserResponse}
// @Failure      400   {object}  models.Response
// @Failure      401   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/users/me [put]
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		models.Unauthorized(c, "Unauthorized")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	result, err := h.userService.UpdateCurrentUser(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		if err == service.ErrUserNotFound {
			models.NotFound(c, "User not found")
			return
		}
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, result)
}

// 3. ChangePassword 修改密码
// ChangePassword godoc
// @Summary      修改密码
// @Description  修改当前登录用户的密码
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      models.ChangePasswordRequest  true  "密码信息"
// @Success      200   {object}  models.Response
// @Failure      400   {object}  models.Response
// @Failure      401   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/users/me/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		models.Unauthorized(c, "Unauthorized")
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	err := h.userService.ChangePassword(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		if err == service.ErrUserNotFound {
			models.NotFound(c, "User not found")
			return
		}
		if err == service.ErrOldPasswordInvalid {
			models.BadRequest(c, "Old password is incorrect")
			return
		}
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, gin.H{
		"message": "Password changed successfully",
	})
}

// 4. GetMyListings 获取我的发布
// GetMyListings godoc
// @Summary      获取我的发布
// @Description  获取当前登录用户发布的房源列表
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page       query     int  false  "页码"     default(1)
// @Param        page_size  query     int  false  "每页数量" default(20)
// @Success      200        {object}  models.Response{data=models.MyListingsResponse}
// @Failure      401        {object}  models.Response
// @Failure      500        {object}  models.Response
// @Router       /api/v1/users/me/listings [get]
func (h *UserHandler) GetMyListings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		models.Unauthorized(c, "Unauthorized")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := h.userService.GetMyListings(c.Request.Context(), userID.(uint), page, pageSize)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.SuccessWithPagination(c, result.Properties, models.NewPagination(page, pageSize, result.Total))
}

// 5. UpdateSettings 更新设置
// UpdateSettings godoc
// @Summary      更新设置
// @Description  更新当前登录用户的设置
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      models.UpdateSettingsRequest  true  "设置信息"
// @Success      200   {object}  models.Response{data=models.UserSettingsResponse}
// @Failure      400   {object}  models.Response
// @Failure      401   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/users/me/settings [put]
func (h *UserHandler) UpdateSettings(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		models.Unauthorized(c, "Unauthorized")
		return
	}

	var req models.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	result, err := h.userService.UpdateSettings(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, result)
}
