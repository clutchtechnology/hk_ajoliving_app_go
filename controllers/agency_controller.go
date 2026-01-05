package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/service"
)

// AgencyHandler 代理公司处理器
//
// AgencyHandler Methods:
// 0. NewAgencyHandler(service service.AgencyService) -> 注入 AgencyService
// 1. ListAgencies(c *gin.Context) -> 获取代理公司列表
// 2. GetAgency(c *gin.Context) -> 获取代理公司详情
// 3. GetAgencyProperties(c *gin.Context) -> 获取代理公司房源列表
// 4. ContactAgency(c *gin.Context) -> 联系代理公司
// 5. SearchAgencies(c *gin.Context) -> 搜索代理公司
type AgencyHandler struct {
	*BaseHandler
	service service.AgencyService
}

// 0. NewAgencyHandler 创建代理公司处理器
func NewAgencyHandler(service service.AgencyService) *AgencyHandler {
	return &AgencyHandler{
		BaseHandler: NewBaseHandler(),
		service:     service,
	}
}

// 1. ListAgencies 获取代理公司列表
// @Summary 获取代理公司列表
// @Description 获取代理公司列表，支持地区、验证状态、评分等筛选
// @Tags 代理公司
// @Accept json
// @Produce json
// @Param district_id query int false "服务地区ID"
// @Param is_verified query bool false "是否已验证"
// @Param min_rating query number false "最低评分"
// @Param keyword query string false "关键词搜索（公司名称、简介）"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(rating) Enums(rating, agent_count, created_at)
// @Param sort_order query string false "排序方向" Enums(asc, desc) default(desc)
// @Success 200 {object} response.PaginatedResponse{data=[]response.AgencyListItemResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/agencies [get]
func (h *AgencyHandler) ListAgencies(c *gin.Context) {
	var filter request.ListAgenciesRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	
	agencies, total, err := h.service.ListAgencies(c.Request.Context(), &filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	
	response.SuccessWithPagination(c, agencies, &response.Pagination{
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Total:     total,
		TotalPage: int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)),
	})
}

// 2. GetAgency 获取代理公司详情
// @Summary 获取代理公司详情
// @Description 根据代理公司ID获取详细信息，包括优秀代理人、房源数量等
// @Tags 代理公司
// @Accept json
// @Produce json
// @Param id path int true "代理公司ID"
// @Success 200 {object} response.Response{data=response.AgencyResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/agencies/{id} [get]
func (h *AgencyHandler) GetAgency(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid agency id")
		return
	}
	
	agency, err := h.service.GetAgency(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "agency not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, agency)
}

// 3. GetAgencyProperties 获取代理公司房源列表
// @Summary 获取代理公司房源列表
// @Description 获取指定代理公司旗下所有代理人的房源
// @Tags 代理公司
// @Accept json
// @Produce json
// @Param id path int true "代理公司ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.PaginatedResponse{data=[]response.PropertyListItemResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/agencies/{id}/properties [get]
func (h *AgencyHandler) GetAgencyProperties(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid agency id")
		return
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	properties, total, err := h.service.GetAgencyProperties(c.Request.Context(), uint(id), page, pageSize)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "agency not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	
	response.SuccessWithPagination(c, properties, &response.Pagination{
		Page:      page,
		PageSize:  pageSize,
		Total:     total,
		TotalPage: int((total + int64(pageSize) - 1) / int64(pageSize)),
	})
}

// 4. ContactAgency 联系代理公司
// @Summary 联系代理公司
// @Description 提交联系代理公司的请求，可以包含相关房产信息
// @Tags 代理公司
// @Accept json
// @Produce json
// @Param id path int true "代理公司ID"
// @Param request body request.ContactAgencyRequest true "联系请求"
// @Success 200 {object} response.Response{data=response.ContactAgencyResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/agencies/{id}/contact [post]
func (h *AgencyHandler) ContactAgency(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid agency id")
		return
	}
	
	var req request.ContactAgencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	
	result, err := h.service.ContactAgency(c.Request.Context(), uint(id), &req)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "agency not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, result)
}

// 5. SearchAgencies 搜索代理公司
// @Summary 搜索代理公司
// @Description 根据关键词搜索代理公司（公司名称、牌照号、简介等）
// @Tags 代理公司
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.PaginatedResponse{data=[]response.AgencyListItemResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/agencies/search [get]
func (h *AgencyHandler) SearchAgencies(c *gin.Context) {
	var filter request.SearchAgenciesRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	
	agencies, total, err := h.service.SearchAgencies(c.Request.Context(), &filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	
	response.SuccessWithPagination(c, agencies, &response.Pagination{
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Total:     total,
		TotalPage: int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)),
	})
}
