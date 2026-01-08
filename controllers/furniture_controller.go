package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// FurnitureController Methods:
// 0. NewFurnitureController(service *services.FurnitureService) -> 注入 FurnitureService
// 1. ListFurniture(c *gin.Context) -> 家具列表
// 2. GetFurnitureCategories(c *gin.Context) -> 家具分类
// 3. GetFurniture(c *gin.Context) -> 家具详情
// 4. CreateFurniture(c *gin.Context) -> 发布家具（需认证）
// 5. UpdateFurniture(c *gin.Context) -> 更新家具（需认证）
// 6. DeleteFurniture(c *gin.Context) -> 删除家具（需认证）
// 7. GetFurnitureImages(c *gin.Context) -> 家具图片
// 8. UpdateFurnitureStatus(c *gin.Context) -> 更新家具状态（需认证）
// 9. GetFeaturedFurniture(c *gin.Context) -> 精选家具

type FurnitureController struct {
	furnitureService *services.FurnitureService
}

// 0. NewFurnitureController -> 注入 FurnitureService
func NewFurnitureController(furnitureService *services.FurnitureService) *FurnitureController {
	return &FurnitureController{
		furnitureService: furnitureService,
	}
}

// 1. ListFurniture -> 家具列表
func (ctrl *FurnitureController) ListFurniture(c *gin.Context) {
	var req models.ListFurnitureRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	response, err := ctrl.furnitureService.ListFurniture(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, response)
}

// 2. GetFurnitureCategories -> 家具分类
func (ctrl *FurnitureController) GetFurnitureCategories(c *gin.Context) {
	categories, err := ctrl.furnitureService.GetFurnitureCategories(c.Request.Context())
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, categories)
}

// 3. GetFurniture -> 家具详情
func (ctrl *FurnitureController) GetFurniture(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid furniture id")
		return
	}

	furniture, err := ctrl.furnitureService.GetFurniture(c.Request.Context(), uint(id))
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "furniture not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, furniture)
}

// 4. CreateFurniture -> 发布家具（需认证）
func (ctrl *FurnitureController) CreateFurniture(c *gin.Context) {
	var req models.CreateFurnitureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	// 从上下文获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "user not authenticated")
		return
	}

	// 获取用户类型（默认为 individual）
	userType := "individual"
	if ut, exists := c.Get("user_type"); exists {
		if utStr, ok := ut.(string); ok {
			userType = utStr
		}
	}

	furniture, err := ctrl.furnitureService.CreateFurniture(c.Request.Context(), userID.(uint), userType, &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Created(c, furniture)
}

// 5. UpdateFurniture -> 更新家具（需认证）
func (ctrl *FurnitureController) UpdateFurniture(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid furniture id")
		return
	}

	var req models.UpdateFurnitureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "user not authenticated")
		return
	}

	furniture, err := ctrl.furnitureService.UpdateFurniture(c.Request.Context(), uint(id), userID.(uint), &req)
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "furniture not found")
			return
		}
		if err == tools.ErrForbidden {
			tools.Forbidden(c, "you don't have permission to update this furniture")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, furniture)
}

// 6. DeleteFurniture -> 删除家具（需认证）
func (ctrl *FurnitureController) DeleteFurniture(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid furniture id")
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "user not authenticated")
		return
	}

	if err := ctrl.furnitureService.DeleteFurniture(c.Request.Context(), uint(id), userID.(uint)); err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "furniture not found")
			return
		}
		if err == tools.ErrForbidden {
			tools.Forbidden(c, "you don't have permission to delete this furniture")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, gin.H{"message": "furniture deleted successfully"})
}

// 7. GetFurnitureImages -> 家具图片
func (ctrl *FurnitureController) GetFurnitureImages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid furniture id")
		return
	}

	images, err := ctrl.furnitureService.GetFurnitureImages(c.Request.Context(), uint(id))
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "furniture not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, images)
}

// 8. UpdateFurnitureStatus -> 更新家具状态（需认证）
func (ctrl *FurnitureController) UpdateFurnitureStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid furniture id")
		return
	}

	var req models.UpdateFurnitureStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "user not authenticated")
		return
	}

	if err := ctrl.furnitureService.UpdateFurnitureStatus(c.Request.Context(), uint(id), userID.(uint), req.Status); err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "furniture not found")
			return
		}
		if err == tools.ErrForbidden {
			tools.Forbidden(c, "you don't have permission to update this furniture")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, gin.H{"message": "furniture status updated successfully"})
}

// 9. GetFeaturedFurniture -> 精选家具
func (ctrl *FurnitureController) GetFeaturedFurniture(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	furniture, err := ctrl.furnitureService.GetFeaturedFurniture(c.Request.Context(), limit)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, furniture)
}
