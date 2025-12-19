package handler

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	pkgErrors "github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/service"
	"github.com/gin-gonic/gin"
)

// FacilityHandler Methods:
// 0. NewFacilityHandler(service *service.FacilityService) -> 注入 FacilityService
// 1. ListFacilities(c *gin.Context) -> 获取设施列表
// 2. GetFacility(c *gin.Context) -> 获取单个设施详情
// 3. CreateFacility(c *gin.Context) -> 创建设施
// 4. UpdateFacility(c *gin.Context) -> 更新设施信息
// 5. DeleteFacility(c *gin.Context) -> 删除设施

// FacilityHandler 设施处理器
type FacilityHandler struct {
	service *service.FacilityService
}

// 0. NewFacilityHandler -> 注入 FacilityService
func NewFacilityHandler(service *service.FacilityService) *FacilityHandler {
	return &FacilityHandler{service: service}
}

// 1. ListFacilities -> 获取设施列表
// ListFacilities godoc
// @Summary      获取设施列表
// @Description  获取设施列表，支持按分类和关键词筛选
// @Tags         Facilities
// @Accept       json
// @Produce      json
// @Param        category    query     string  false  "设施分类 (building, unit)"
// @Param        keyword     query     string  false  "搜索关键词"
// @Param        page        query     int     false  "页码" default(1)
// @Param        page_size   query     int     false  "每页数量" default(50)
// @Param        sort_by     query     string  false  "排序字段" default(sort_order)
// @Param        sort_order  query     string  false  "排序方式 (asc, desc)" default(asc)
// @Success      200  {object}  response.Response{data=response.FacilityListResponse}
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/v1/facilities [get]
func (h *FacilityHandler) ListFacilities(c *gin.Context) {
	var req request.ListFacilitiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.ListFacilities(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// 2. GetFacility -> 获取单个设施详情
// GetFacility godoc
// @Summary      获取设施详情
// @Description  根据ID获取单个设施的详细信息
// @Tags         Facilities
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "设施ID"
// @Success      200  {object}  response.Response{data=response.FacilityResponse}
// @Failure      400  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/v1/facilities/{id} [get]
func (h *FacilityHandler) GetFacility(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid facility id")
		return
	}

	result, err := h.service.GetFacility(c.Request.Context(), uint(id))
	if err != nil {
		if err == pkgErrors.ErrNotFound {
			response.NotFound(c, "facility not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// 3. CreateFacility -> 创建设施
// CreateFacility godoc
// @Summary      创建设施
// @Description  创建新的设施（需要管理员权限）
// @Tags         Facilities
// @Accept       json
// @Produce      json
// @Param        body  body      request.CreateFacilityRequest  true  "设施信息"
// @Success      201   {object}  response.Response{data=response.FacilityResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /api/v1/facilities [post]
// @Security     BearerAuth
func (h *FacilityHandler) CreateFacility(c *gin.Context) {
	var req request.CreateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.CreateFacility(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, result)
}

// 4. UpdateFacility -> 更新设施信息
// UpdateFacility godoc
// @Summary      更新设施信息
// @Description  更新指定设施的信息（需要管理员权限）
// @Tags         Facilities
// @Accept       json
// @Produce      json
// @Param        id    path      int                            true   "设施ID"
// @Param        body  body      request.UpdateFacilityRequest  true   "设施信息"
// @Success      200   {object}  response.Response{data=response.FacilityResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /api/v1/facilities/{id} [put]
// @Security     BearerAuth
func (h *FacilityHandler) UpdateFacility(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid facility id")
		return
	}

	var req request.UpdateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.UpdateFacility(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err == pkgErrors.ErrNotFound {
			response.NotFound(c, "facility not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// 5. DeleteFacility -> 删除设施
// DeleteFacility godoc
// @Summary      删除设施
// @Description  删除指定的设施（需要管理员权限）
// @Tags         Facilities
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "设施ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/v1/facilities/{id} [delete]
// @Security     BearerAuth
func (h *FacilityHandler) DeleteFacility(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid facility id")
		return
	}

	if err := h.service.DeleteFacility(c.Request.Context(), uint(id)); err != nil {
		if err == pkgErrors.ErrNotFound {
			response.NotFound(c, "facility not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "facility deleted successfully"})
}
