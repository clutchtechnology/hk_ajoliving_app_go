package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DistrictController Methods:
// 0. NewDistrictController(service *services.DistrictService) -> 注入 DistrictService
// 1. ListDistricts(c *gin.Context) -> 获取地区列表
// 2. GetDistrict(c *gin.Context) -> 获取地区详情
// 3. GetDistrictProperties(c *gin.Context) -> 获取地区房源
// 4. GetDistrictEstates(c *gin.Context) -> 获取地区屋苑
// 5. GetDistrictStatistics(c *gin.Context) -> 获取地区统计数据

type DistrictController struct {
	service *services.DistrictService
}

// 0. NewDistrictController 构造函数
func NewDistrictController(service *services.DistrictService) *DistrictController {
	return &DistrictController{service: service}
}

// 1. ListDistricts 获取地区列表
// GET /api/v1/districts
func (ctrl *DistrictController) ListDistricts(c *gin.Context) {
	var req models.ListDistrictsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	districts, err := ctrl.service.ListDistricts(c.Request.Context(), req.Region)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, districts)
}

// 2. GetDistrict 获取地区详情
// GET /api/v1/districts/:id
func (ctrl *DistrictController) GetDistrict(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid district id")
		return
	}

	district, err := ctrl.service.GetDistrict(c.Request.Context(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == tools.ErrNotFound {
			tools.NotFound(c, "district not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, district)
}

// 3. GetDistrictProperties 获取地区房源
// GET /api/v1/districts/:id/properties
func (ctrl *DistrictController) GetDistrictProperties(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid district id")
		return
	}

	var req models.GetDistrictPropertiesRequest
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

	properties, total, err := ctrl.service.GetDistrictProperties(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == tools.ErrNotFound {
			tools.NotFound(c, "district not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	// 构建响应
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	response := gin.H{
		"properties":  properties,
		"total":       total,
		"page":        req.Page,
		"page_size":   req.PageSize,
		"total_pages": totalPages,
	}

	tools.Success(c, response)
}

// 4. GetDistrictEstates 获取地区屋苑
// GET /api/v1/districts/:id/estates
func (ctrl *DistrictController) GetDistrictEstates(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid district id")
		return
	}

	var req models.GetDistrictEstatesRequest
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

	estates, total, err := ctrl.service.GetDistrictEstates(c.Request.Context(), uint(id), req.Page, req.PageSize)
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == tools.ErrNotFound {
			tools.NotFound(c, "district not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	// 构建响应
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	response := gin.H{
		"estates":     estates,
		"total":       total,
		"page":        req.Page,
		"page_size":   req.PageSize,
		"total_pages": totalPages,
	}

	tools.Success(c, response)
}

// 5. GetDistrictStatistics 获取地区统计数据
// GET /api/v1/districts/:id/statistics
func (ctrl *DistrictController) GetDistrictStatistics(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid district id")
		return
	}

	statistics, err := ctrl.service.GetDistrictStatistics(c.Request.Context(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == tools.ErrNotFound {
			tools.NotFound(c, "district not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, statistics)
}
