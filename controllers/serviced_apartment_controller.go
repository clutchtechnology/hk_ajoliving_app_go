package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/gin-gonic/gin"
)

// ServicedApartmentHandler Methods:
// 0. NewServicedApartmentHandler(service service.ServicedApartmentService) -> 注入 ServicedApartmentService
// 1. ListServicedApartments(c *gin.Context) -> 获取服务式住宅列表
// 2. GetServicedApartment(c *gin.Context) -> 获取单个服务式住宅详情
// 3. GetApartmentUnits(c *gin.Context) -> 获取服务式住宅单元列表
// 4. GetFeaturedApartments(c *gin.Context) -> 获取精选服务式住宅
// 5. CreateServicedApartment(c *gin.Context) -> 创建服务式住宅（需要认证）
// 6. UpdateServicedApartment(c *gin.Context) -> 更新服务式住宅（需要认证）
// 7. DeleteServicedApartment(c *gin.Context) -> 删除服务式住宅（需要认证）

type ServicedApartmentHandler struct {
	service service.ServicedApartmentService
}

// 0. NewServicedApartmentHandler -> 注入 ServicedApartmentService
func NewServicedApartmentHandler(service service.ServicedApartmentService) *ServicedApartmentHandler {
	return &ServicedApartmentHandler{service: service}
}

// 1. ListServicedApartments -> 获取服务式住宅列表
// @Summary 获取服务式住宅列表
// @Description 根据筛选条件获取服务式住宅列表
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Param district_id query uint false "地区ID"
// @Param min_price query number false "最低价格"
// @Param max_price query number false "最高价格"
// @Param bedrooms query int false "卧室数"
// @Param min_stay_days query int false "最短入住天数"
// @Param status query string false "状态(active/inactive/closed)"
// @Param is_featured query bool false "是否精选"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(created_at)
// @Param sort_order query string false "排序顺序(asc/desc)" default(desc)
// @Success 200 {object} models.Response{data=[]models.ServicedApartmentListItemResponse}
// @Router /serviced-apartments [get]
func (h *ServicedApartmentHandler) ListServicedApartments(c *gin.Context) {
	var req models.ListServicedApartmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	apartments, total, err := h.service.ListServicedApartments(c.Request.Context(), &req)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.SuccessWithPagination(c, apartments, req.Page, req.PageSize, total)
}

// 2. GetServicedApartment -> 获取单个服务式住宅详情
// @Summary 获取服务式住宅详情
// @Description 根据ID获取服务式住宅详细信息
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Param id path int true "服务式住宅ID"
// @Success 200 {object} models.Response{data=models.ServicedApartmentResponse}
// @Failure 404 {object} models.Response
// @Router /serviced-apartments/{id} [get]
func (h *ServicedApartmentHandler) GetServicedApartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		models.BadRequest(c, "invalid apartment id")
		return
	}

	apartment, err := h.service.GetServicedApartment(c.Request.Context(), uint(id))
	if err != nil {
		models.NotFound(c, "serviced apartment not found")
		return
	}

	models.Success(c, apartment)
}

// 3. GetApartmentUnits -> 获取服务式住宅单元列表
// @Summary 获取服务式住宅单元列表
// @Description 获取指定服务式住宅的所有可用单元
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Param id path int true "服务式住宅ID"
// @Success 200 {object} models.Response{data=[]models.ServicedApartmentUnitResponse}
// @Failure 404 {object} models.Response
// @Router /serviced-apartments/{id}/units [get]
func (h *ServicedApartmentHandler) GetApartmentUnits(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		models.BadRequest(c, "invalid apartment id")
		return
	}

	units, err := h.service.GetApartmentUnits(c.Request.Context(), uint(id))
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, units)
}

// 4. GetFeaturedApartments -> 获取精选服务式住宅
// @Summary 获取精选服务式住宅
// @Description 获取平台推荐的精选服务式住宅列表
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} models.Response{data=[]models.ServicedApartmentListItemResponse}
// @Router /serviced-apartments/featured [get]
func (h *ServicedApartmentHandler) GetFeaturedApartments(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	apartments, err := h.service.GetFeaturedApartments(c.Request.Context(), limit)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, apartments)
}

// 5. CreateServicedApartment -> 创建服务式住宅（需要认证）
// @Summary 创建服务式住宅
// @Description 创建新的服务式住宅（需要管理员或公司账号权限）
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.CreateServicedApartmentRequest true "服务式住宅信息"
// @Success 201 {object} models.Response{data=models.ServicedApartmentResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /serviced-apartments [post]
func (h *ServicedApartmentHandler) CreateServicedApartment(c *gin.Context) {
	var req models.CreateServicedApartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	// 获取当前用户ID
	userID := c.GetUint("user_id")

	apartment, err := h.service.CreateServicedApartment(c.Request.Context(), userID, &req)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.Created(c, apartment)
}

// 6. UpdateServicedApartment -> 更新服务式住宅（需要认证）
// @Summary 更新服务式住宅
// @Description 更新服务式住宅信息（需要管理员或所有者权限）
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "服务式住宅ID"
// @Param body body models.UpdateServicedApartmentRequest true "更新信息"
// @Success 200 {object} models.Response{data=models.ServicedApartmentResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /serviced-apartments/{id} [put]
func (h *ServicedApartmentHandler) UpdateServicedApartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		models.BadRequest(c, "invalid apartment id")
		return
	}

	var req models.UpdateServicedApartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	apartment, err := h.service.UpdateServicedApartment(c.Request.Context(), uint(id), &req)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, apartment)
}

// 7. DeleteServicedApartment -> 删除服务式住宅（需要认证）
// @Summary 删除服务式住宅
// @Description 删除服务式住宅（需要管理员或所有者权限）
// @Tags 服务式住宅
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "服务式住宅ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /serviced-apartments/{id} [delete]
func (h *ServicedApartmentHandler) DeleteServicedApartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		models.BadRequest(c, "invalid apartment id")
		return
	}

	if err := h.service.DeleteServicedApartment(c.Request.Context(), uint(id)); err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, gin.H{"message": "serviced apartment deleted successfully"})
}
