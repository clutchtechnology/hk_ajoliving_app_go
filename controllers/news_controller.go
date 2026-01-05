package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
)

// NewsHandler 新闻处理器
type NewsHandler struct {
	*BaseHandler
	service services.NewsService
}

// NewNewsHandler 创建新闻处理器
func NewNewsHandler(service services.NewsService) *NewsHandler {
	return &NewsHandler{
		BaseHandler: NewBaseHandler(),
		service:     service,
	}
}

// ListNews 获取新闻列表
// @Summary 获取新闻列表
// @Description 获取新闻列表，支持分类、状态、标签等筛选
// @Tags 新闻
// @Accept json
// @Produce json
// @Param category_id query int false "分类ID"
// @Param status query string false "状态" Enums(draft, published, archived)
// @Param is_featured query bool false "是否精选"
// @Param is_hot query bool false "是否热门"
// @Param is_top query bool false "是否置顶"
// @Param keyword query string false "关键词搜索"
// @Param tag query string false "标签筛选"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(published_at)
// @Param sort_order query string false "排序方向" Enums(asc, desc) default(desc)
// @Success 200 {object} models.PaginatedResponse{data=[]models.NewsListItemResponse}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/news [get]
func (h *NewsHandler) ListNews(c *gin.Context) {
	var filter models.ListNewsRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}
	
	news, total, err := h.service.ListNews(c.Request.Context(), &filter)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}
	
	tools.SuccessWithPagination(c, news, &tools.Pagination{
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Total:     total,
		TotalPage: int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)),
	})
}

// GetNewsCategories 获取新闻分类列表
// @Summary 获取新闻分类列表
// @Description 获取所有启用的新闻分类
// @Tags 新闻
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=[]models.NewsCategoryResponse}
// @Failure 500 {object} models.Response
// @Router /api/v1/news/categories [get]
func (h *NewsHandler) GetNewsCategories(c *gin.Context) {
	categories, err := h.service.GetNewsCategories(c.Request.Context())
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}
	
	tools.Success(c, categories)
}

// GetNews 获取新闻详情
// @Summary 获取新闻详情
// @Description 根据新闻ID获取详细信息，自动增加浏览量
// @Tags 新闻
// @Accept json
// @Produce json
// @Param id path int true "新闻ID"
// @Success 200 {object} models.Response{data=models.NewsResponse}
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/news/{id} [get]
func (h *NewsHandler) GetNews(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid news id")
		return
	}
	
	news, err := h.service.GetNews(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			tools.NotFound(c, "news not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}
	
	tools.Success(c, news)
}

// GetHotNews 获取热门新闻
// @Summary 获取热门新闻
// @Description 获取浏览量最高的热门新闻（最多10条）
// @Tags 新闻
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=[]models.NewsListItemResponse}
// @Failure 500 {object} models.Response
// @Router /api/v1/news/hot [get]
func (h *NewsHandler) GetHotNews(c *gin.Context) {
	news, err := h.service.GetHotNews(c.Request.Context())
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}
	
	tools.Success(c, news)
}

// GetFeaturedNews 获取精选新闻
// @Summary 获取精选新闻
// @Description 获取编辑精选的新闻（最多10条）
// @Tags 新闻
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=[]models.NewsListItemResponse}
// @Failure 500 {object} models.Response
// @Router /api/v1/news/featured [get]
func (h *NewsHandler) GetFeaturedNews(c *gin.Context) {
	news, err := h.service.GetFeaturedNews(c.Request.Context())
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}
	
	tools.Success(c, news)
}

// GetLatestNews 获取最新新闻
// @Summary 获取最新新闻
// @Description 获取最新发布的新闻（最多10条）
// @Tags 新闻
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=[]models.NewsListItemResponse}
// @Failure 500 {object} models.Response
// @Router /api/v1/news/latest [get]
func (h *NewsHandler) GetLatestNews(c *gin.Context) {
	news, err := h.service.GetLatestNews(c.Request.Context())
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}
	
	tools.Success(c, news)
}

// GetRelatedNews 获取相关新闻
// @Summary 获取相关新闻
// @Description 根据新闻ID获取同分类的相关新闻（最多5条）
// @Tags 新闻
// @Accept json
// @Produce json
// @Param id path int true "新闻ID"
// @Success 200 {object} models.Response{data=[]models.RelatedNewsResponse}
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/news/{id}/related [get]
func (h *NewsHandler) GetRelatedNews(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid news id")
		return
	}
	
	news, err := h.service.GetRelatedNews(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			tools.NotFound(c, "news not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}
	
	tools.Success(c, news)
}
