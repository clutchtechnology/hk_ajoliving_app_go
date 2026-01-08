package controllers

import (
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// SearchController Methods:
// 0. NewSearchController(service *services.SearchService) -> 注入 SearchService
// 1. GlobalSearch(c *gin.Context) -> 全局搜索
// 2. SearchProperties(c *gin.Context) -> 搜索房产
// 3. SearchEstates(c *gin.Context) -> 搜索屋苑
// 4. SearchAgents(c *gin.Context) -> 搜索代理人
// 5. GetSearchSuggestions(c *gin.Context) -> 获取搜索建议
// 6. GetSearchHistory(c *gin.Context) -> 获取搜索历史

type SearchController struct {
	service *services.SearchService
}

// 0. NewSearchController 构造函数
func NewSearchController(service *services.SearchService) *SearchController {
	return &SearchController{service: service}
}

// 1. GlobalSearch 全局搜索
// GET /api/v1/search
func (ctrl *SearchController) GlobalSearch(c *gin.Context) {
	var req models.GlobalSearchRequest
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

	// 获取可选的用户ID
	var userID *uint
	if userIDValue, exists := c.Get("user_id"); exists {
		if uid, ok := userIDValue.(uint); ok {
			userID = &uid
		}
	}

	// 获取请求信息
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	results, err := ctrl.service.GlobalSearch(c.Request.Context(), &req, userID, ipAddress, userAgent)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, results)
}

// 2. SearchProperties 搜索房产
// GET /api/v1/search/properties
func (ctrl *SearchController) SearchProperties(c *gin.Context) {
	var req models.SearchPropertiesRequest
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

	// 获取可选的用户ID
	var userID *uint
	if userIDValue, exists := c.Get("user_id"); exists {
		if uid, ok := userIDValue.(uint); ok {
			userID = &uid
		}
	}

	// 获取请求信息
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	results, total, err := ctrl.service.SearchProperties(c.Request.Context(), &req, userID, ipAddress, userAgent)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	// 构建响应
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	response := gin.H{
		"properties":  results,
		"total":       total,
		"page":        req.Page,
		"page_size":   req.PageSize,
		"total_pages": totalPages,
	}

	tools.Success(c, response)
}

// 3. SearchEstates 搜索屋苑
// GET /api/v1/search/estates
func (ctrl *SearchController) SearchEstates(c *gin.Context) {
	var req models.SearchEstatesRequest
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

	// 获取可选的用户ID
	var userID *uint
	if userIDValue, exists := c.Get("user_id"); exists {
		if uid, ok := userIDValue.(uint); ok {
			userID = &uid
		}
	}

	// 获取请求信息
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	results, total, err := ctrl.service.SearchEstates(c.Request.Context(), &req, userID, ipAddress, userAgent)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	// 构建响应
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	response := gin.H{
		"estates":     results,
		"total":       total,
		"page":        req.Page,
		"page_size":   req.PageSize,
		"total_pages": totalPages,
	}

	tools.Success(c, response)
}

// 4. SearchAgents 搜索代理人
// GET /api/v1/search/agents
func (ctrl *SearchController) SearchAgents(c *gin.Context) {
	var req models.SearchAgentsRequest
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

	// 获取可选的用户ID
	var userID *uint
	if userIDValue, exists := c.Get("user_id"); exists {
		if uid, ok := userIDValue.(uint); ok {
			userID = &uid
		}
	}

	// 获取请求信息
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	results, total, err := ctrl.service.SearchAgents(c.Request.Context(), &req, userID, ipAddress, userAgent)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	// 构建响应
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	response := gin.H{
		"agents":      results,
		"total":       total,
		"page":        req.Page,
		"page_size":   req.PageSize,
		"total_pages": totalPages,
	}

	tools.Success(c, response)
}

// 5. GetSearchSuggestions 获取搜索建议
// GET /api/v1/search/suggestions
func (ctrl *SearchController) GetSearchSuggestions(c *gin.Context) {
	var req models.GetSearchSuggestionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	if req.Limit == 0 {
		req.Limit = 10
	}

	suggestions, err := ctrl.service.GetSearchSuggestions(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, suggestions)
}

// 6. GetSearchHistory 获取搜索历史（需认证）
// GET /api/v1/search/history
func (ctrl *SearchController) GetSearchHistory(c *gin.Context) {
	var req models.GetSearchHistoryRequest
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

	// 获取用户ID（可选）
	var userID *uint
	if userIDValue, exists := c.Get("user_id"); exists {
		if uid, ok := userIDValue.(uint); ok {
			userID = &uid
		}
	}

	result, err := ctrl.service.GetSearchHistory(c.Request.Context(), userID, &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}
