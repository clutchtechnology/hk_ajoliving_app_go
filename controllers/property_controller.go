package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// PropertyController 房产控制器
type PropertyController struct {
	propertyService *services.PropertyService
}

// NewPropertyController 创建房产控制器
func NewPropertyController(propertyService *services.PropertyService) *PropertyController {
	return &PropertyController{
		propertyService: propertyService,
	}
}

// ListProperties 获取房产列表
func (ctrl *PropertyController) ListProperties(c *gin.Context) {
	var req models.ListPropertiesRequest
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

	properties, err := ctrl.propertyService.ListProperties(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, properties)
}

// GetProperty 获取房产详情
func (ctrl *PropertyController) GetProperty(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid property id")
		return
	}

	property, err := ctrl.propertyService.GetProperty(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "property not found" {
			tools.NotFound(c, "property not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, property)
}

// CreateProperty 创建房产
func (ctrl *PropertyController) CreateProperty(c *gin.Context) {
	// 从中间件获取用户信息
	userID, _ := c.Get("user_id")
	userType, _ := c.Get("user_type")

	var req models.CreatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	property, err := ctrl.propertyService.CreateProperty(
		c.Request.Context(),
		userID.(uint),
		userType.(string),
		&req,
	)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Created(c, property)
}

// UpdateProperty 更新房产
func (ctrl *PropertyController) UpdateProperty(c *gin.Context) {
	// 从中间件获取用户ID
	userID, _ := c.Get("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid property id")
		return
	}

	var req models.UpdatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	property, err := ctrl.propertyService.UpdateProperty(
		c.Request.Context(),
		uint(id),
		userID.(uint),
		&req,
	)
	if err != nil {
		if err.Error() == "property not found" {
			tools.NotFound(c, "property not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, property)
}

// DeleteProperty 删除房产
func (ctrl *PropertyController) DeleteProperty(c *gin.Context) {
	// 从中间件获取用户ID
	userID, _ := c.Get("user_id")

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid property id")
		return
	}

	err = ctrl.propertyService.DeleteProperty(c.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		if err.Error() == "property not found" {
			tools.NotFound(c, "property not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, gin.H{"message": "property deleted successfully"})
}

// GetSimilarProperties 获取相似房源
func (ctrl *PropertyController) GetSimilarProperties(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid property id")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	similar, err := ctrl.propertyService.GetSimilarProperties(c.Request.Context(), uint(id), limit)
	if err != nil {
		if err.Error() == "property not found" {
			tools.NotFound(c, "property not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, similar)
}

// GetFeaturedProperties 获取精选房源
func (ctrl *PropertyController) GetFeaturedProperties(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	featured, err := ctrl.propertyService.GetFeaturedProperties(c.Request.Context(), limit)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, featured)
}

// GetHotProperties 获取热门房源
func (ctrl *PropertyController) GetHotProperties(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	hot, err := ctrl.propertyService.GetHotProperties(c.Request.Context(), limit)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, hot)
}

// ========== 买房模块 ==========

// ListBuyProperties 买房房源列表
func (ctrl *PropertyController) ListBuyProperties(c *gin.Context) {
	var req models.ListPropertiesRequest
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

	// 强制设置为出售类型
	saleType := "sale"
	req.ListingType = &saleType

	properties, err := ctrl.propertyService.ListProperties(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, properties)
}

// ListNewProperties 新房列表（一手房）
func (ctrl *PropertyController) ListNewProperties(c *gin.Context) {
	var req models.ListPropertiesRequest
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

	// 强制设置为出售类型 + 新房
	saleType := "sale"
	req.ListingType = &saleType
	
	// 新房筛选：可以通过 estate_no 不为空来判断（楼盘编号）
	// 或者按创建时间排序，获取最新发布的
	// 这里简化处理，后续可根据实际需求优化
	req.SortBy = "created_at_desc"

	properties, err := ctrl.propertyService.ListProperties(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, properties)
}

// ListSecondhandProperties 二手房列表
func (ctrl *PropertyController) ListSecondhandProperties(c *gin.Context) {
	var req models.ListPropertiesRequest
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

	// 强制设置为出售类型
	saleType := "sale"
	req.ListingType = &saleType

	properties, err := ctrl.propertyService.ListProperties(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, properties)
}

// ========== 租房模块 ==========

// ListRentProperties 租房房源列表
func (ctrl *PropertyController) ListRentProperties(c *gin.Context) {
	var req models.ListPropertiesRequest
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

	// 强制设置为租赁类型
	rentType := "rent"
	req.ListingType = &rentType

	properties, err := ctrl.propertyService.ListProperties(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, properties)
}

// ListShortTermRent 短租房源列表
func (ctrl *PropertyController) ListShortTermRent(c *gin.Context) {
	var req models.ListPropertiesRequest
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

	// 强制设置为租赁类型
	rentType := "rent"
	req.ListingType = &rentType

	// 短租筛选：价格相对较低（按月租）
	// 这里简化处理，实际可以添加一个 rent_term 字段区分长短租
	properties, err := ctrl.propertyService.ListProperties(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, properties)
}

// ListLongTermRent 长租房源列表
func (ctrl *PropertyController) ListLongTermRent(c *gin.Context) {
	var req models.ListPropertiesRequest
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

	// 强制设置为租赁类型
	rentType := "rent"
	req.ListingType = &rentType

	// 长租筛选：一般为标准月租
	properties, err := ctrl.propertyService.ListProperties(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, properties)
}
