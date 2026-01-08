package controllers

import (
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// StatisticsController 统计控制器
type StatisticsController struct {
	service *services.StatisticsService
}

// NewStatisticsController 创建统计控制器实例
func NewStatisticsController(service *services.StatisticsService) *StatisticsController {
	return &StatisticsController{service: service}
}

// GetOverviewStatistics 获取总览统计
// @Summary 获取总览统计
// @Description 获取平台总览统计数据，包括房产、用户、代理、屋苑等核心指标
// @Tags 统计分析
// @Accept json
// @Produce json
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} tools.Response{data=models.OverviewStatisticsResponse}
// @Router /api/v1/statistics/overview [get]
func (ctrl *StatisticsController) GetOverviewStatistics(c *gin.Context) {
	var req models.GetOverviewStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	stats, err := ctrl.service.GetOverviewStatistics(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, stats)
}

// GetPropertyStatistics 获取房产统计
// @Summary 获取房产统计
// @Description 获取房产详细统计数据，包括数量、价格、面积、房型、地区分布等
// @Tags 统计分析
// @Accept json
// @Produce json
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Param district_id query int false "地区ID"
// @Success 200 {object} tools.Response{data=models.PropertyStatisticsResponse}
// @Router /api/v1/statistics/properties [get]
func (ctrl *StatisticsController) GetPropertyStatistics(c *gin.Context) {
	var req models.GetPropertyStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	stats, err := ctrl.service.GetPropertyStatistics(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, stats)
}

// GetTransactionStatistics 获取成交统计
// @Summary 获取成交统计
// @Description 获取成交统计数据，包括成交数量、金额、地区分布、月度趋势等
// @Tags 统计分析
// @Accept json
// @Produce json
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Param district_id query int false "地区ID"
// @Success 200 {object} tools.Response{data=models.TransactionStatisticsResponse}
// @Router /api/v1/statistics/transactions [get]
func (ctrl *StatisticsController) GetTransactionStatistics(c *gin.Context) {
	var req models.GetTransactionStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	stats, err := ctrl.service.GetTransactionStatistics(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, stats)
}

// GetUserStatistics 获取用户统计
// @Summary 获取用户统计
// @Description 获取用户统计数据，包括用户总数、新增用户、活跃度、代理统计等
// @Tags 统计分析
// @Accept json
// @Produce json
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} tools.Response{data=models.UserStatisticsResponse}
// @Router /api/v1/statistics/users [get]
func (ctrl *StatisticsController) GetUserStatistics(c *gin.Context) {
	var req models.GetUserStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	stats, err := ctrl.service.GetUserStatistics(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, stats)
}
