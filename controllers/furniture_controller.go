package controllers

// FurnitureHandler Methods:
// 0. NewFurnitureHandler(furnitureService *services.FurnitureService) -> 注入 FurnitureService
// 1. ListFurniture(c *gin.Context) -> 家具列表
// 2. GetFurnitureCategories(c *gin.Context) -> 家具分类
// 3. GetFurniture(c *gin.Context) -> 家具详情
// 4. CreateFurniture(c *gin.Context) -> 发布家具
// 5. UpdateFurniture(c *gin.Context) -> 更新家具
// 6. DeleteFurniture(c *gin.Context) -> 删除家具
// 7. GetFurnitureImages(c *gin.Context) -> 家具图片
// 8. UpdateFurnitureStatus(c *gin.Context) -> 更新家具状态
// 9. GetFeaturedFurniture(c *gin.Context) -> 精选家具

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
)

// FurnitureHandlerInterface 家具处理器接口
type FurnitureHandlerInterface interface {
	ListFurniture(c *gin.Context)           // 1. 家具列表
	GetFurnitureCategories(c *gin.Context)  // 2. 家具分类
	GetFurniture(c *gin.Context)            // 3. 家具详情
	CreateFurniture(c *gin.Context)         // 4. 发布家具
	UpdateFurniture(c *gin.Context)         // 5. 更新家具
	DeleteFurniture(c *gin.Context)         // 6. 删除家具
	GetFurnitureImages(c *gin.Context)      // 7. 家具图片
	UpdateFurnitureStatus(c *gin.Context)   // 8. 更新家具状态
	GetFeaturedFurniture(c *gin.Context)    // 9. 精选家具
}

// FurnitureHandler 家具处理器
type FurnitureHandler struct {
	furnitureService *services.FurnitureService
}

// 0. NewFurnitureHandler 注入 FurnitureService
func NewFurnitureHandler(furnitureService *services.FurnitureService) *FurnitureHandler {
	return &FurnitureHandler{
		furnitureService: furnitureService,
	}
}

// 1. ListFurniture 家具列表
// ListFurniture godoc
// @Summary      家具列表
// @Description  获取家具列表，支持多种筛选条件
// @Tags         Furniture
// @Accept       json
// @Produce      json
// @Param        category_id           query     int     false  "分类ID"
// @Param        min_price             query     number  false  "最低价格"
// @Param        max_price             query     number  false  "最高价格"
// @Param        condition             query     string  false  "新旧程度(new/like_new/good/fair/poor)"
// @Param        brand                 query     string  false  "品牌"
// @Param        delivery_district_id  query     int     false  "交收地区ID"
// @Param        delivery_method       query     string  false  "交收方法(self_pickup/delivery/negotiable)"
// @Param        status                query     string  false  "状态(available/reserved/sold)"
// @Param        keyword               query     string  false  "关键字搜索"
// @Param        sort_by               query     string  false  "排序字段"   default(created_at)
// @Param        sort_order            query     string  false  "排序方向"   default(desc)
// @Param        page                  query     int     false  "页码"       default(1)
// @Param        page_size             query     int     false  "每页数量"   default(20)
// @Success      200                   {object}  models.PaginatedResponse{data=[]models.FurnitureListItemResponse}
// @Failure      400                   {object}  models.Response
// @Failure      500                   {object}  models.Response
// @Router       /api/v1/furniture [get]
func (h *FurnitureHandler) ListFurniture(c *gin.Context) {
	var req models.ListFurnitureRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	furniture, total, err := h.furnitureService.ListFurniture(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.SuccessWithPagination(c, furniture, models.NewPagination(req.Page, req.PageSize, total))
}

// 2. GetFurnitureCategories 家具分类
// GetFurnitureCategories godoc
// @Summary      家具分类
// @Description  获取家具分类列表（包含子分类）
// @Tags         Furniture
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Response{data=[]models.FurnitureCategoryResponse}
// @Failure      500  {object}  models.Response
// @Router       /api/v1/furniture/categories [get]
func (h *FurnitureHandler) GetFurnitureCategories(c *gin.Context) {
	categories, err := h.furnitureService.GetFurnitureCategories(c.Request.Context())
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, categories)
}

// 3. GetFurniture 家具详情
// GetFurniture godoc
// @Summary      家具详情
// @Description  获取家具的详细信息
// @Tags         Furniture
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "家具ID"
// @Success      200  {object}  models.Response{data=models.FurnitureResponse}
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/v1/furniture/{id} [get]
func (h *FurnitureHandler) GetFurniture(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "Invalid furniture ID")
		return
	}

	furniture, err := h.furnitureService.GetFurniture(c.Request.Context(), uint(id))
	if err != nil {
		if err == services.ErrFurnitureNotFound {
			tools.NotFound(c, "Furniture not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, furniture)
}

// 4. CreateFurniture 发布家具
// CreateFurniture godoc
// @Summary      发布家具
// @Description  创建新的家具信息（需要认证）
// @Tags         Furniture
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      models.CreateFurnitureRequest  true  "家具信息"
// @Success      201   {object}  models.Response{data=models.CreateFurnitureResponse}
// @Failure      400   {object}  models.Response
// @Failure      401   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/furniture [post]
func (h *FurnitureHandler) CreateFurniture(c *gin.Context) {
	var req models.CreateFurnitureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	// 获取当前用户ID（从JWT中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "User not authenticated")
		return
	}

	furniture, err := h.furnitureService.CreateFurniture(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		if err == services.ErrCategoryNotFound {
			tools.BadRequest(c, "Category not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Created(c, furniture)
}

// 5. UpdateFurniture 更新家具
// UpdateFurniture godoc
// @Summary      更新家具
// @Description  更新家具信息（需要认证且为发布者）
// @Tags         Furniture
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                             true  "家具ID"
// @Param        body  body      models.UpdateFurnitureRequest  true  "家具信息"
// @Success      200   {object}  models.Response{data=models.UpdateFurnitureResponse}
// @Failure      400   {object}  models.Response
// @Failure      401   {object}  models.Response
// @Failure      403   {object}  models.Response
// @Failure      404   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/furniture/{id} [put]
func (h *FurnitureHandler) UpdateFurniture(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "Invalid furniture ID")
		return
	}

	var req models.UpdateFurnitureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "User not authenticated")
		return
	}

	furniture, err := h.furnitureService.UpdateFurniture(c.Request.Context(), userID.(uint), uint(id), &req)
	if err != nil {
		if err == services.ErrFurnitureNotFound {
			tools.NotFound(c, "Furniture not found")
			return
		}
		if err == services.ErrNotFurnitureOwner {
			tools.Forbidden(c, "You are not the owner of this furniture")
			return
		}
		if err == services.ErrCategoryNotFound {
			tools.BadRequest(c, "Category not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, furniture)
}

// 6. DeleteFurniture 删除家具
// DeleteFurniture godoc
// @Summary      删除家具
// @Description  删除家具信息（需要认证且为发布者）
// @Tags         Furniture
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "家具ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      401  {object}  models.Response
// @Failure      403  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/v1/furniture/{id} [delete]
func (h *FurnitureHandler) DeleteFurniture(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "Invalid furniture ID")
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "User not authenticated")
		return
	}

	err = h.furnitureService.DeleteFurniture(c.Request.Context(), userID.(uint), uint(id))
	if err != nil {
		if err == services.ErrFurnitureNotFound {
			tools.NotFound(c, "Furniture not found")
			return
		}
		if err == services.ErrNotFurnitureOwner {
			tools.Forbidden(c, "You are not the owner of this furniture")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, gin.H{"message": "Furniture deleted successfully"})
}

// 7. GetFurnitureImages 家具图片
// GetFurnitureImages godoc
// @Summary      家具图片
// @Description  获取家具的所有图片
// @Tags         Furniture
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "家具ID"
// @Success      200  {object}  models.Response{data=[]models.FurnitureImageResponse}
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/v1/furniture/{id}/images [get]
func (h *FurnitureHandler) GetFurnitureImages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "Invalid furniture ID")
		return
	}

	images, err := h.furnitureService.GetFurnitureImages(c.Request.Context(), uint(id))
	if err != nil {
		if err == services.ErrFurnitureNotFound {
			tools.NotFound(c, "Furniture not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, images)
}

// 8. UpdateFurnitureStatus 更新家具状态
// UpdateFurnitureStatus godoc
// @Summary      更新家具状态
// @Description  更新家具的状态（需要认证且为发布者）
// @Tags         Furniture
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                                     true  "家具ID"
// @Param        body  body      models.UpdateFurnitureStatusRequest  true  "状态信息"
// @Success      200   {object}  models.Response{data=models.UpdateFurnitureStatusResponse}
// @Failure      400   {object}  models.Response
// @Failure      401   {object}  models.Response
// @Failure      403   {object}  models.Response
// @Failure      404   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/furniture/{id}/status [put]
func (h *FurnitureHandler) UpdateFurnitureStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "Invalid furniture ID")
		return
	}

	var req models.UpdateFurnitureStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "User not authenticated")
		return
	}

	result, err := h.furnitureService.UpdateFurnitureStatus(c.Request.Context(), userID.(uint), uint(id), &req)
	if err != nil {
		if err == services.ErrFurnitureNotFound {
			tools.NotFound(c, "Furniture not found")
			return
		}
		if err == services.ErrNotFurnitureOwner {
			tools.Forbidden(c, "You are not the owner of this furniture")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 9. GetFeaturedFurniture 精选家具
// GetFeaturedFurniture godoc
// @Summary      精选家具
// @Description  获取精选家具列表（按浏览量和收藏数排序）
// @Tags         Furniture
// @Accept       json
// @Produce      json
// @Param        limit  query     int  false  "数量限制"  default(10)
// @Success      200    {object}  models.Response{data=[]models.FurnitureListItemResponse}
// @Failure      500    {object}  models.Response
// @Router       /api/v1/furniture/featured [get]
func (h *FurnitureHandler) GetFeaturedFurniture(c *gin.Context) {
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	furniture, err := h.furnitureService.GetFeaturedFurniture(c.Request.Context(), limit)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, furniture)
}
