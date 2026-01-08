package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FacilityController Methods:
// 0. NewFacilityController(service *services.FacilityService) -> 注入 FacilityService
// 1. ListFacilities(c *gin.Context) -> 获取设施列表
// 2. GetFacility(c *gin.Context) -> 获取设施详情
// 3. CreateFacility(c *gin.Context) -> 创建设施（需认证）
// 4. UpdateFacility(c *gin.Context) -> 更新设施（需认证）
// 5. DeleteFacility(c *gin.Context) -> 删除设施（需认证）

type FacilityController struct {
	service *services.FacilityService
}

// 0. NewFacilityController 构造函数
func NewFacilityController(service *services.FacilityService) *FacilityController {
	return &FacilityController{service: service}
}

// 1. ListFacilities 获取设施列表
// GET /api/v1/facilities
func (ctrl *FacilityController) ListFacilities(c *gin.Context) {
	var req models.ListFacilitiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	facilities, err := ctrl.service.ListFacilities(c.Request.Context(), req.Category)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, facilities)
}

// 2. GetFacility 获取设施详情
// GET /api/v1/facilities/:id
func (ctrl *FacilityController) GetFacility(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid facility id")
		return
	}

	facility, err := ctrl.service.GetFacility(c.Request.Context(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == tools.ErrNotFound {
			tools.NotFound(c, "facility not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, facility)
}

// 3. CreateFacility 创建设施（需认证）
// POST /api/v1/facilities
func (ctrl *FacilityController) CreateFacility(c *gin.Context) {
	var req models.CreateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	facility, err := ctrl.service.CreateFacility(c.Request.Context(), &req)
	if err != nil {
		if err.Error() == "facility name already exists" {
			tools.BadRequest(c, err.Error())
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Created(c, facility)
}

// 4. UpdateFacility 更新设施（需认证）
// PUT /api/v1/facilities/:id
func (ctrl *FacilityController) UpdateFacility(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid facility id")
		return
	}

	var req models.UpdateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	facility, err := ctrl.service.UpdateFacility(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == tools.ErrNotFound {
			tools.NotFound(c, "facility not found")
			return
		}
		if err.Error() == "facility name already exists" {
			tools.BadRequest(c, err.Error())
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, facility)
}

// 5. DeleteFacility 删除设施（需认证）
// DELETE /api/v1/facilities/:id
func (ctrl *FacilityController) DeleteFacility(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid facility id")
		return
	}

	if err := ctrl.service.DeleteFacility(c.Request.Context(), uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound || err == tools.ErrNotFound {
			tools.NotFound(c, "facility not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, gin.H{"message": "facility deleted successfully"})
}
