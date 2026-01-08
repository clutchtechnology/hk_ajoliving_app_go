package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AgencyController Methods:
// 0. NewAgencyController(service *services.AgencyService) -> 注入 AgencyService
// 1. ListAgencies(c *gin.Context) -> 获取代理公司列表
// 2. GetAgency(c *gin.Context) -> 获取代理公司详情
// 3. GetAgencyProperties(c *gin.Context) -> 获取代理公司房源列表
// 4. ContactAgency(c *gin.Context) -> 联系代理公司
// 5. SearchAgencies(c *gin.Context) -> 搜索代理公司

type AgencyController struct {
	service *services.AgencyService
}

// 0. NewAgencyController 构造函数
func NewAgencyController(service *services.AgencyService) *AgencyController {
	return &AgencyController{service: service}
}

// 1. ListAgencies 获取代理公司列表
// GET /api/v1/agencies
func (ctrl *AgencyController) ListAgencies(c *gin.Context) {
	var req models.ListAgenciesRequest
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

	result, err := ctrl.service.ListAgencies(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 2. GetAgency 获取代理公司详情
// GET /api/v1/agencies/:id
func (ctrl *AgencyController) GetAgency(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid agency id")
		return
	}

	agency, err := ctrl.service.GetAgency(c.Request.Context(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == tools.ErrNotFound {
			tools.NotFound(c, "agency not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, agency)
}

// 3. GetAgencyProperties 获取代理公司房源列表
// GET /api/v1/agencies/:id/properties
func (ctrl *AgencyController) GetAgencyProperties(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid agency id")
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	properties, total, err := ctrl.service.GetAgencyProperties(c.Request.Context(), uint(id), page, pageSize)
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == tools.ErrNotFound {
			tools.NotFound(c, "agency not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	// 构建响应
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response := gin.H{
		"properties":  properties,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	}

	tools.Success(c, response)
}

// 4. ContactAgency 联系代理公司
// POST /api/v1/agencies/:id/contact
func (ctrl *AgencyController) ContactAgency(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid agency id")
		return
	}

	var req models.ContactAgencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	// 获取可选的用户ID（如果已登录）
	var userID *uint
	if userIDValue, exists := c.Get("user_id"); exists {
		if uid, ok := userIDValue.(uint); ok {
			userID = &uid
		}
	}

	result, err := ctrl.service.ContactAgency(c.Request.Context(), uint(id), userID, &req)
	if err != nil {
		if err == gorm.ErrRecordNotFound || err == tools.ErrNotFound {
			tools.NotFound(c, "agency not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Created(c, result)
}

// 5. SearchAgencies 搜索代理公司
// GET /api/v1/agencies/search
func (ctrl *AgencyController) SearchAgencies(c *gin.Context) {
	var req models.SearchAgenciesRequest
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

	result, err := ctrl.service.SearchAgencies(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}
