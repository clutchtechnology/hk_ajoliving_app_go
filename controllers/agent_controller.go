package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// AgentController Methods:
// 0. NewAgentController(service *services.AgentService) -> 注入 AgentService
// 1. ListAgents(c *gin.Context) -> 代理人列表
// 2. GetAgent(c *gin.Context) -> 代理人详情
// 3. GetAgentProperties(c *gin.Context) -> 代理人房源列表
// 4. ContactAgent(c *gin.Context) -> 联系代理人

type AgentController struct {
	service *services.AgentService
}

// 0. NewAgentController 构造函数
func NewAgentController(service *services.AgentService) *AgentController {
	return &AgentController{service: service}
}

// 1. ListAgents 代理人列表
// @Summary 代理人列表
// @Tags Agent
// @Produce json
// @Param license_type query string false "牌照类型 (individual, salesperson)"
// @Param agency_id query int false "代理公司ID"
// @Param district_id query int false "服务地区ID"
// @Param status query string false "状态 (active, inactive, suspended)"
// @Param is_verified query bool false "是否已验证"
// @Param specialization query string false "专长领域"
// @Param keyword query string false "关键词搜索"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} tools.Response{data=models.PaginatedAgentsResponse}
// @Router /api/v1/agents [get]
func (ctrl *AgentController) ListAgents(c *gin.Context) {
	var req models.ListAgentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := ctrl.service.ListAgents(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 2. GetAgent 代理人详情
// @Summary 代理人详情
// @Tags Agent
// @Produce json
// @Param id path int true "代理人ID"
// @Success 200 {object} tools.Response{data=models.AgentDetailResponse}
// @Router /api/v1/agents/{id} [get]
func (ctrl *AgentController) GetAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid agent id")
		return
	}

	agent, err := ctrl.service.GetAgent(c.Request.Context(), uint(id))
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "agent not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, agent)
}

// 3. GetAgentProperties 代理人房源列表
// @Summary 代理人房源列表
// @Tags Agent
// @Produce json
// @Param id path int true "代理人ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} tools.Response
// @Router /api/v1/agents/{id}/properties [get]
func (ctrl *AgentController) GetAgentProperties(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid agent id")
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

	result, err := ctrl.service.GetAgentProperties(c.Request.Context(), uint(id), page, pageSize)
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "agent not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 4. ContactAgent 联系代理人
// @Summary 联系代理人
// @Tags Agent
// @Accept json
// @Produce json
// @Param id path int true "代理人ID"
// @Param body body models.ContactAgentRequest true "联系信息"
// @Success 200 {object} tools.Response
// @Router /api/v1/agents/{id}/contact [post]
func (ctrl *AgentController) ContactAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid agent id")
		return
	}

	var req models.ContactAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	// 尝试获取用户ID（如果已登录）
	var userID *uint
	if uid, exists := c.Get("user_id"); exists {
		if id, ok := uid.(uint); ok {
			userID = &id
		}
	}

	err = ctrl.service.ContactAgent(c.Request.Context(), uint(id), userID, &req)
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "agent not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, gin.H{"message": "contact request sent successfully"})
}
