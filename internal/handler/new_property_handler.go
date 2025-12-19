package handler

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	pkgErrors "github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// NewPropertyHandler 新楼盘处理器
// Methods:
// 0. NewNewPropertyHandler(service, logger) -> 注入依赖
// 1. ListNewDevelopments(c *gin.Context) -> 获取新楼盘列表 GET /api/v1/new-properties
// 2. GetNewDevelopment(c *gin.Context) -> 获取新楼盘详情 GET /api/v1/new-properties/:id
// 3. GetDevelopmentUnits(c *gin.Context) -> 获取楼盘户型列表 GET /api/v1/new-properties/:id/units
// 4. GetFeaturedNewDevelopments(c *gin.Context) -> 获取精选新楼盘 GET /api/v1/new-properties/featured
type NewPropertyHandler struct {
	service service.NewPropertyService
	logger  *zap.Logger
}

// NewNewPropertyHandler 创建新楼盘处理器实例
func NewNewPropertyHandler(svc service.NewPropertyService, logger *zap.Logger) *NewPropertyHandler {
	return &NewPropertyHandler{
		service: svc,
		logger:  logger,
	}
}

// ListNewDevelopments 获取新楼盘列表
// @Summary 获取新楼盘列表
// @Description 获取新楼盘列表，支持筛选和分页
// @Tags 新楼盘
// @Accept json
// @Produce json
// @Param district_id query int false "地区ID"
// @Param developer query string false "发展商名称"
// @Param status query string false "状态：upcoming/presale/selling/completed"
// @Param min_price query number false "最低价格"
// @Param max_price query number false "最高价格"
// @Param bedrooms query int false "卧室数"
// @Param school_net query string false "校网"
// @Param is_featured query bool false "是否精选"
// @Param sort_by query string false "排序字段"
// @Param sort_order query string false "排序方向：asc/desc"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=[]response.NewDevelopmentListItemResponse}
// @Router /api/v1/new-properties [get]
func (h *NewPropertyHandler) ListNewDevelopments(c *gin.Context) {
	var req request.ListNewDevelopmentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	items, total, err := h.service.ListNewDevelopments(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to list new developments", zap.Error(err))
		response.InternalError(c, "获取新楼盘列表失败")
		return
	}

	pagination := response.NewPagination(req.Page, req.PageSize, total)
	response.SuccessWithPagination(c, items, pagination)
}

// GetNewDevelopment 获取新楼盘详情
// @Summary 获取新楼盘详情
// @Description 根据ID获取新楼盘详情
// @Tags 新楼盘
// @Accept json
// @Produce json
// @Param id path int true "新楼盘ID"
// @Success 200 {object} response.Response{data=response.NewDevelopmentResponse}
// @Failure 404 {object} response.Response
// @Router /api/v1/new-properties/{id} [get]
func (h *NewPropertyHandler) GetNewDevelopment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的新楼盘ID")
		return
	}

	newDevelopment, err := h.service.GetNewDevelopment(c.Request.Context(), uint(id))
	if err != nil {
		if err == pkgErrors.ErrNotFound {
			response.NotFound(c, "新楼盘不存在")
			return
		}
		h.logger.Error("Failed to get new development", zap.Uint64("id", id), zap.Error(err))
		response.InternalError(c, "获取新楼盘详情失败")
		return
	}

	response.Success(c, newDevelopment)
}

// GetDevelopmentUnits 获取楼盘户型列表
// @Summary 获取楼盘户型列表
// @Description 根据新楼盘ID获取所有户型
// @Tags 新楼盘
// @Accept json
// @Produce json
// @Param id path int true "新楼盘ID"
// @Success 200 {object} response.Response{data=[]response.NewDevelopmentLayoutResponse}
// @Failure 404 {object} response.Response
// @Router /api/v1/new-properties/{id}/units [get]
func (h *NewPropertyHandler) GetDevelopmentUnits(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的新楼盘ID")
		return
	}

	units, err := h.service.GetDevelopmentUnits(c.Request.Context(), uint(id))
	if err != nil {
		if err == pkgErrors.ErrNotFound {
			response.NotFound(c, "新楼盘不存在")
			return
		}
		h.logger.Error("Failed to get development units", zap.Uint64("id", id), zap.Error(err))
		response.InternalError(c, "获取楼盘户型失败")
		return
	}

	response.Success(c, units)
}

// GetFeaturedNewDevelopments 获取精选新楼盘
// @Summary 获取精选新楼盘
// @Description 获取精选新楼盘列表
// @Tags 新楼盘
// @Accept json
// @Produce json
// @Param limit query int false "数量限制" default(10)
// @Success 200 {object} response.Response{data=[]response.NewDevelopmentListItemResponse}
// @Router /api/v1/new-properties/featured [get]
func (h *NewPropertyHandler) GetFeaturedNewDevelopments(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	items, err := h.service.GetFeaturedNewDevelopments(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to get featured new developments", zap.Error(err))
		response.InternalError(c, "获取精选新楼盘失败")
		return
	}

	response.Success(c, items)
}
