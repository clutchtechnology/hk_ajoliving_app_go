package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// ServicedApartmentController 服务式住宅控制器
// Methods:
// 1. ListServicedApartments(c *gin.Context) -> 获取服务式住宅列表
// 2. GetServicedApartment(c *gin.Context) -> 获取服务式住宅详情
// 3. CreateServicedApartment(c *gin.Context) -> 创建服务式住宅（需要认证）
// 4. UpdateServicedApartment(c *gin.Context) -> 更新服务式住宅（需要认证）
// 5. DeleteServicedApartment(c *gin.Context) -> 删除服务式住宅（需要认证）
// 6. GetServicedApartmentUnits(c *gin.Context) -> 获取房型列表
// 7. GetServicedApartmentImages(c *gin.Context) -> 获取图片列表
type ServicedApartmentController struct {
	service *services.ServicedApartmentService
}

// NewServicedApartmentController 创建服务式住宅控制器
func NewServicedApartmentController(service *services.ServicedApartmentService) *ServicedApartmentController {
	return &ServicedApartmentController{service: service}
}

// ListServicedApartments 获取服务式住宅列表
// @Summary 获取服务式住宅列表
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Param district_id query int false "地区ID"
// @Param status query string false "状态: active, closed"
// @Param min_rating query number false "最低评分（0-5）"
// @Param is_featured query bool false "是否精选"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} tools.Response{data=models.PaginatedServicedApartmentsResponse}
// @Router /api/v1/serviced-apartments [get]
func (ctrl *ServicedApartmentController) ListServicedApartments(c *gin.Context) {
	var filter models.ListServicedApartmentsRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := ctrl.service.ListServicedApartments(c.Request.Context(), &filter)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// GetServicedApartment 获取服务式住宅详情
// @Summary 获取服务式住宅详情
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Param id path int true "服务式住宅ID"
// @Success 200 {object} tools.Response{data=models.ServicedApartmentDetailResponse}
// @Failure 404 {object} tools.Response
// @Router /api/v1/serviced-apartments/{id} [get]
func (ctrl *ServicedApartmentController) GetServicedApartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid serviced apartment id")
		return
	}

	apartment, err := ctrl.service.GetServicedApartment(c.Request.Context(), uint(id))
	if err != nil {
		tools.NotFound(c, err.Error())
		return
	}

	tools.Success(c, apartment)
}

// CreateServicedApartment 创建服务式住宅（需要认证）
// @Summary 创建服务式住宅
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body models.CreateServicedApartmentRequest true "创建请求"
// @Success 201 {object} tools.Response{data=models.ServicedApartmentDetailResponse}
// @Failure 400 {object} tools.Response
// @Failure 401 {object} tools.Response
// @Router /api/v1/serviced-apartments [post]
func (ctrl *ServicedApartmentController) CreateServicedApartment(c *gin.Context) {
	// 获取当前用户ID（由JWT中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "user not authenticated")
		return
	}

	var req models.CreateServicedApartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	apartment, err := ctrl.service.CreateServicedApartment(c.Request.Context(), &req, userID.(uint))
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Created(c, apartment)
}

// UpdateServicedApartment 更新服务式住宅（需要认证）
// @Summary 更新服务式住宅
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务式住宅ID"
// @Param request body models.UpdateServicedApartmentRequest true "更新请求"
// @Success 200 {object} tools.Response{data=models.ServicedApartmentDetailResponse}
// @Failure 400 {object} tools.Response
// @Failure 401 {object} tools.Response
// @Failure 403 {object} tools.Response
// @Failure 404 {object} tools.Response
// @Router /api/v1/serviced-apartments/{id} [put]
func (ctrl *ServicedApartmentController) UpdateServicedApartment(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid serviced apartment id")
		return
	}

	var req models.UpdateServicedApartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	apartment, err := ctrl.service.UpdateServicedApartment(c.Request.Context(), uint(id), &req, userID.(uint))
	if err != nil {
		if err.Error() == "permission denied" {
			tools.Forbidden(c, "you don't have permission to update this apartment")
			return
		}
		tools.NotFound(c, err.Error())
		return
	}

	tools.Success(c, apartment)
}

// DeleteServicedApartment 删除服务式住宅（需要认证）
// @Summary 删除服务式住宅
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "服务式住宅ID"
// @Success 200 {object} tools.Response
// @Failure 400 {object} tools.Response
// @Failure 401 {object} tools.Response
// @Failure 403 {object} tools.Response
// @Failure 404 {object} tools.Response
// @Router /api/v1/serviced-apartments/{id} [delete]
func (ctrl *ServicedApartmentController) DeleteServicedApartment(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "user not authenticated")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid serviced apartment id")
		return
	}

	err = ctrl.service.DeleteServicedApartment(c.Request.Context(), uint(id), userID.(uint))
	if err != nil {
		if err.Error() == "permission denied" {
			tools.Forbidden(c, "you don't have permission to delete this apartment")
			return
		}
		tools.NotFound(c, err.Error())
		return
	}

	tools.Success(c, gin.H{"message": "serviced apartment deleted successfully"})
}

// GetServicedApartmentUnits 获取房型列表
// @Summary 获取服务式住宅的房型列表
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Param id path int true "服务式住宅ID"
// @Success 200 {object} tools.Response{data=[]models.ServicedApartmentUnit}
// @Failure 404 {object} tools.Response
// @Router /api/v1/serviced-apartments/{id}/units [get]
func (ctrl *ServicedApartmentController) GetServicedApartmentUnits(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid serviced apartment id")
		return
	}

	units, err := ctrl.service.GetServicedApartmentUnits(c.Request.Context(), uint(id))
	if err != nil {
		tools.NotFound(c, err.Error())
		return
	}

	tools.Success(c, units)
}

// GetServicedApartmentImages 获取图片列表
// @Summary 获取服务式住宅的图片列表
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Param id path int true "服务式住宅ID"
// @Success 200 {object} tools.Response{data=[]models.ServicedApartmentImage}
// @Failure 404 {object} tools.Response
// @Router /api/v1/serviced-apartments/{id}/images [get]
func (ctrl *ServicedApartmentController) GetServicedApartmentImages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid serviced apartment id")
		return
	}

	images, err := ctrl.service.GetServicedApartmentImages(c.Request.Context(), uint(id))
	if err != nil {
		tools.NotFound(c, err.Error())
		return
	}

	tools.Success(c, images)
}
