package controllers

import (
	"net/http"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// StatisticsHandler 统计处理器
type StatisticsHandler struct {
	*BaseHandler
	service service.StatisticsService
}

// NewStatisticsHandler 创建统计处理器实例
func NewStatisticsHandler(baseHandler *BaseHandler, service service.StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{
		BaseHandler: baseHandler,
		service:     service,
	}
}

// GetOverviewStatistics 获取总览统计
// @Summary 获取总览统计
// @Description 获取平台总览统计数据，包括房产、用户、代理人、成交等汇总信息
// @Tags 统计分析
// @Accept json
// @Produce json
// @Param period query string false "统计周期" Enums(day, week, month, year) default(month)
// @Success 200 {object} response.Response{data=response.OverviewStatisticsResponse} "统计数据"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/statistics/overview [get]
func (h *StatisticsHandler) GetOverviewStatistics(c *gin.Context) {
	var req request.GetOverviewStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.Logger.Warn("参数验证失败", zap.Error(err))
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置默认周期
	if req.Period == "" {
		req.Period = "month"
	}

	data, err := h.service.GetOverviewStatistics(c.Request.Context(), &req)
	if err != nil {
		h.Logger.Error("获取总览统计失败", zap.Error(err))
		response.InternalError(c, "获取统计数据失败")
		return
	}

	response.Success(c, data)
}

// GetPropertyStatistics 获取房产统计
// @Summary 获取房产统计
// @Description 获取房产统计数据，包括数量分布、价格统计、趋势分析等
// @Tags 统计分析
// @Accept json
// @Produce json
// @Param period query string false "统计周期" Enums(day, week, month, year) default(month)
// @Param district_id query integer false "地区ID"
// @Param estate_id query integer false "屋苑ID"
// @Param listing_type query string false "房源类型" Enums(rent, sale)
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} response.Response{data=response.PropertyStatisticsResponse} "统计数据"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/statistics/properties [get]
func (h *StatisticsHandler) GetPropertyStatistics(c *gin.Context) {
	var req request.GetPropertyStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.Logger.Warn("参数验证失败", zap.Error(err))
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置默认周期
	if req.Period == "" {
		req.Period = "month"
	}

	data, err := h.service.GetPropertyStatistics(c.Request.Context(), &req)
	if err != nil {
		h.Logger.Error("获取房产统计失败", zap.Error(err))
		response.InternalError(c, "获取统计数据失败")
		return
	}

	response.Success(c, data)
}

// GetTransactionStatistics 获取成交统计
// @Summary 获取成交统计
// @Description 获取成交统计数据，包括成交量、成交金额、趋势分析、地区分布等
// @Tags 统计分析
// @Accept json
// @Produce json
// @Param period query string false "统计周期" Enums(day, week, month, year) default(month)
// @Param district_id query integer false "地区ID"
// @Param estate_id query integer false "屋苑ID"
// @Param listing_type query string false "房源类型" Enums(rent, sale)
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} response.Response{data=response.TransactionStatisticsResponse} "统计数据"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/statistics/transactions [get]
func (h *StatisticsHandler) GetTransactionStatistics(c *gin.Context) {
	var req request.GetTransactionStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.Logger.Warn("参数验证失败", zap.Error(err))
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置默认周期
	if req.Period == "" {
		req.Period = "month"
	}

	data, err := h.service.GetTransactionStatistics(c.Request.Context(), &req)
	if err != nil {
		h.Logger.Error("获取成交统计失败", zap.Error(err))
		response.InternalError(c, "获取统计数据失败")
		return
	}

	response.Success(c, data)
}

// GetUserStatistics 获取用户统计
// @Summary 获取用户统计
// @Description 获取用户统计数据，包括用户数量、增长趋势、状态分布、角色分布等
// @Tags 统计分析
// @Accept json
// @Produce json
// @Param period query string false "统计周期" Enums(day, week, month, year) default(month)
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} response.Response{data=response.UserStatisticsResponse} "统计数据"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /api/v1/statistics/users [get]
func (h *StatisticsHandler) GetUserStatistics(c *gin.Context) {
	var req request.GetUserStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		h.Logger.Warn("参数验证失败", zap.Error(err))
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 设置默认周期
	if req.Period == "" {
		req.Period = "month"
	}

	data, err := h.service.GetUserStatistics(c.Request.Context(), &req)
	if err != nil {
		h.Logger.Error("获取用户统计失败", zap.Error(err))
		response.InternalError(c, "获取统计数据失败")
		return
	}

	response.Success(c, data)
}
