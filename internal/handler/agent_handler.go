package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/service"
)

// AgentHandler 代理人处理器
type AgentHandler struct {
	*BaseHandler
	service service.AgentService
}

// NewAgentHandler 创建代理人处理器
func NewAgentHandler(service service.AgentService) *AgentHandler {
	return &AgentHandler{
		BaseHandler: NewBaseHandler(),
		service:     service,
	}
}

// ListAgents 获取代理人列表
// @Summary 获取代理人列表
// @Description 获取代理人列表，支持代理公司、地区、评分等筛选
// @Tags 代理人
// @Accept json
// @Produce json
// @Param agency_id query int false "代理公司ID"
// @Param district_id query int false "服务地区ID"
// @Param status query string false "状态" Enums(active, inactive, suspended)
// @Param is_verified query bool false "是否已验证"
// @Param specialization query string false "专业领域"
// @Param min_rating query number false "最低评分"
// @Param keyword query string false "关键词搜索"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(rating)
// @Param sort_order query string false "排序方向" Enums(asc, desc) default(desc)
// @Success 200 {object} response.PaginatedResponse{data=[]response.AgentListItemResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/agents [get]
func (h *AgentHandler) ListAgents(c *gin.Context) {
	var filter request.ListAgentsRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	
	agents, total, err := h.service.ListAgents(c.Request.Context(), &filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	
	response.SuccessWithPagination(c, agents, &response.Pagination{
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Total:     total,
		TotalPage: int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)),
	})
}

// GetAgent 获取代理人详情
// @Summary 获取代理人详情
// @Description 根据代理人ID获取详细信息
// @Tags 代理人
// @Accept json
// @Produce json
// @Param id path int true "代理人ID"
// @Success 200 {object} response.Response{data=response.AgentResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/agents/{id} [get]
func (h *AgentHandler) GetAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid agent id")
		return
	}
	
	agent, err := h.service.GetAgent(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "agent not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, agent)
}

// GetAgentProperties 获取代理人房源列表
// @Summary 获取代理人房源列表
// @Description 获取指定代理人的所有房源
// @Tags 代理人
// @Accept json
// @Produce json
// @Param id path int true "代理人ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.PaginatedResponse{data=[]response.PropertyListItemResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/agents/{id}/properties [get]
func (h *AgentHandler) GetAgentProperties(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid agent id")
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
	
	properties, total, err := h.service.GetAgentProperties(c.Request.Context(), uint(id), page, pageSize)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "agent not found")
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

// ContactAgent 联系代理人
// @Summary 联系代理人
// @Description 提交联系代理人的请求
// @Tags 代理人
// @Accept json
// @Produce json
// @Param id path int true "代理人ID"
// @Param request body request.ContactAgentRequest true "联系信息"
// @Success 200 {object} response.Response{data=response.AgentContactResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/agents/{id}/contact [post]
func (h *AgentHandler) ContactAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid agent id")
		return
	}
	
	var req request.ContactAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	
	// 获取用户ID（如果已登录）
	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = &id
		}
	}
	
	contactResp, err := h.service.ContactAgent(c.Request.Context(), uint(id), userID, &req)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "agent not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, contactResp)
}
