package controllers

import (
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/gin-gonic/gin"
)

// SearchHandler Methods:
// 0. NewSearchHandler(service *services.SearchService) -> 注入 SearchService
// 1. GlobalSearch(c *gin.Context) -> 全局搜索
// 2. SearchProperties(c *gin.Context) -> 搜索房产
// 3. SearchEstates(c *gin.Context) -> 搜索屋苑
// 4. SearchAgents(c *gin.Context) -> 搜索代理人
// 5. GetSearchSuggestions(c *gin.Context) -> 获取搜索建议
// 6. GetSearchHistory(c *gin.Context) -> 获取搜索历史

// SearchHandler 搜索处理器
type SearchHandler struct {
	service *services.SearchService
}

// 0. NewSearchHandler -> 注入 SearchService
func NewSearchHandler(service *services.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

// 1. GlobalSearch -> 全局搜索
// GlobalSearch godoc
// @Summary      全局搜索
// @Description  在房产、屋苑、代理人、新闻等多个类别中进行全局搜索
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        keyword    query     string  true   "搜索关键词"
// @Param        page       query     int     false  "页码" default(1)
// @Param        page_size  query     int     false  "每页数量" default(20)
// @Success      200  {object}  models.Response{data=models.GlobalSearchResponse}
// @Failure      400  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/v1/search [get]
func (h *SearchHandler) GlobalSearch(c *gin.Context) {
	var req models.GlobalSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.GlobalSearch(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 2. SearchProperties -> 搜索房产
// SearchProperties godoc
// @Summary      搜索房产
// @Description  搜索房产，支持关键词和多种筛选条件
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        keyword        query     string   true   "搜索关键词"
// @Param        listing_type   query     string   false  "房源类型 (rent, sale)"
// @Param        district_id    query     int      false  "地区ID"
// @Param        min_price      query     number   false  "最低价格"
// @Param        max_price      query     number   false  "最高价格"
// @Param        bedrooms       query     int      false  "卧室数量"
// @Param        property_type  query     string   false  "物业类型"
// @Param        page           query     int      false  "页码" default(1)
// @Param        page_size      query     int      false  "每页数量" default(20)
// @Success      200  {object}  models.Response{data=models.SearchPropertiesResponse}
// @Failure      400  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/v1/search/properties [get]
func (h *SearchHandler) SearchProperties(c *gin.Context) {
	var req models.SearchPropertiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.SearchProperties(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 3. SearchEstates -> 搜索屋苑
// SearchEstates godoc
// @Summary      搜索屋苑
// @Description  搜索屋苑，支持关键词和地区筛选
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        keyword      query     string  true   "搜索关键词"
// @Param        district_id  query     int     false  "地区ID"
// @Param        page         query     int     false  "页码" default(1)
// @Param        page_size    query     int     false  "每页数量" default(20)
// @Success      200  {object}  models.Response{data=models.SearchEstatesResponse}
// @Failure      400  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/v1/search/estates [get]
func (h *SearchHandler) SearchEstates(c *gin.Context) {
	var req models.SearchEstatesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.SearchEstates(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 4. SearchAgents -> 搜索代理人
// SearchAgents godoc
// @Summary      搜索代理人
// @Description  搜索代理人，支持关键词和地区筛选
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        keyword      query     string  true   "搜索关键词"
// @Param        district_id  query     int     false  "地区ID"
// @Param        page         query     int     false  "页码" default(1)
// @Param        page_size    query     int     false  "每页数量" default(20)
// @Success      200  {object}  models.Response{data=models.SearchAgentsResponse}
// @Failure      400  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/v1/search/agents [get]
func (h *SearchHandler) SearchAgents(c *gin.Context) {
	var req models.SearchAgentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.SearchAgents(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 5. GetSearchSuggestions -> 获取搜索建议
// GetSearchSuggestions godoc
// @Summary      获取搜索建议
// @Description  根据输入的关键词获取搜索建议
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        keyword  query     string  true   "搜索关键词"
// @Param        type     query     string  false  "搜索类型 (property, estate, agent, agency)"
// @Param        limit    query     int     false  "返回数量" default(10)
// @Success      200  {object}  models.Response{data=models.SearchSuggestionsResponse}
// @Failure      400  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/v1/search/suggestions [get]
func (h *SearchHandler) GetSearchSuggestions(c *gin.Context) {
	var req models.GetSearchSuggestionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.GetSearchSuggestions(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 6. GetSearchHistory -> 获取搜索历史
// GetSearchHistory godoc
// @Summary      获取搜索历史
// @Description  获取用户的搜索历史记录（需要登录）
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        type       query     string  false  "搜索类型"
// @Param        page       query     int     false  "页码" default(1)
// @Param        page_size  query     int     false  "每页数量" default(20)
// @Success      200  {object}  models.Response{data=models.SearchHistoryResponse}
// @Failure      400  {object}  models.Response
// @Failure      401  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/v1/search/history [get]
// @Security     BearerAuth
func (h *SearchHandler) GetSearchHistory(c *gin.Context) {
	var req models.GetSearchHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	// 获取当前用户ID（从JWT中间件设置）
	userID, exists := c.Get("user_id")
	if !exists {
		tools.Unauthorized(c, "user not authenticated")
		return
	}

	uid := userID.(uint)
	result, err := h.service.GetSearchHistory(c.Request.Context(), &uid, &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}
