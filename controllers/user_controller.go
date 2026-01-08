package controllers

import (
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *services.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GetCurrentUser 获取当前用户信息
func (ctrl *UserController) GetCurrentUser(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "unauthorized")
		return
	}

	user, err := ctrl.userService.GetUserByID(c.Request.Context(), userID.(uint))
	if err != nil {
		if err.Error() == "user not found" {
			tools.NotFound(c, "user not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, user.ToUserResponse())
}

// UpdateCurrentUser 更新当前用户信息
func (ctrl *UserController) UpdateCurrentUser(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "unauthorized")
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	user, err := ctrl.userService.UpdateUser(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		if err.Error() == "user not found" {
			tools.NotFound(c, "user not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, user.ToUserResponse())
}

// GetMyListings 获取我的发布（房源、家具等）
func (ctrl *UserController) GetMyListings(c *gin.Context) {
	// 从中间件获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "unauthorized")
		return
	}

	// 获取查询参数
	listingType := c.Query("type") // property, furniture, all
	var typeFilter *string
	if listingType == "property" || listingType == "furniture" {
		typeFilter = &listingType
	}

	// 获取用户发布列表
	listings, err := ctrl.userService.GetUserListings(c.Request.Context(), userID.(uint), typeFilter)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, listings)
}
