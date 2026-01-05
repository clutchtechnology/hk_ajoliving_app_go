package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
)

// PriceIndexHandler 楼价指数处理器
//
// PriceIndexHandler Methods:
// 0. NewPriceIndexHandler(service service.PriceIndexService) -> 注入 PriceIndexService
// 1. GetPriceIndex(c *gin.Context) -> 获取楼价指数列表
// 2. GetLatestPriceIndex(c *gin.Context) -> 获取最新楼价指数
// 3. GetDistrictPriceIndex(c *gin.Context) -> 获取地区楼价指数
// 4. GetEstatePriceIndex(c *gin.Context) -> 获取屋苑楼价指数
// 5. GetPriceTrends(c *gin.Context) -> 获取价格走势
// 6. ComparePriceIndex(c *gin.Context) -> 对比楼价指数
// 7. ExportPriceData(c *gin.Context) -> 导出价格数据
// 8. GetPriceIndexHistory(c *gin.Context) -> 获取历史楼价指数
// 9. CreatePriceIndex(c *gin.Context) -> 创建楼价指数（需认证）
// 10. UpdatePriceIndex(c *gin.Context) -> 更新楼价指数（需认证）
type PriceIndexHandler struct {
	*BaseHandler
	service service.PriceIndexService
}

// 0. NewPriceIndexHandler 创建楼价指数处理器
func NewPriceIndexHandler(service service.PriceIndexService) *PriceIndexHandler {
	return &PriceIndexHandler{
		BaseHandler: NewBaseHandler(),
		service:     service,
	}
}

// 1. GetPriceIndex 获取楼价指数列表
// @Summary 获取楼价指数列表
// @Description 获取楼价指数列表，支持类型、地区、屋苑、周期等筛选
// @Tags 楼价指数
// @Accept json
// @Produce json
// @Param index_type query string false "指数类型" Enums(overall, district, estate, property_type)
// @Param district_id query int false "地区ID"
// @Param estate_id query int false "屋苑ID"
// @Param property_type query string false "物业类型"
// @Param start_period query string false "开始周期 YYYY-MM"
// @Param end_period query string false "结束周期 YYYY-MM"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} models.PaginatedResponse{data=[]models.PriceIndexListItemResponse}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/price-index [get]
func (h *PriceIndexHandler) GetPriceIndex(c *gin.Context) {
	var filter models.GetPriceIndexRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		models.BadRequest(c, err.Error())
		return
	}
	
	indices, total, err := h.service.GetPriceIndex(c.Request.Context(), &filter)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}
	
	models.SuccessWithPagination(c, indices, &models.Pagination{
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Total:     total,
		TotalPage: int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)),
	})
}

// 2. GetLatestPriceIndex 获取最新楼价指数
// @Summary 获取最新楼价指数
// @Description 获取最新的整体、各地区和各物业类型楼价指数
// @Tags 楼价指数
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.LatestPriceIndexResponse}
// @Failure 500 {object} models.Response
// @Router /api/v1/price-index/latest [get]
func (h *PriceIndexHandler) GetLatestPriceIndex(c *gin.Context) {
	result, err := h.service.GetLatestPriceIndex(c.Request.Context())
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}
	
	models.Success(c, result)
}

// 3. GetDistrictPriceIndex 获取地区楼价指数
// @Summary 获取地区楼价指数
// @Description 获取指定地区的历史楼价指数数据
// @Tags 楼价指数
// @Accept json
// @Produce json
// @Param districtId path int true "地区ID"
// @Param start_period query string false "开始周期 YYYY-MM"
// @Param end_period query string false "结束周期 YYYY-MM"
// @Param limit query int false "返回记录数" default(12)
// @Success 200 {object} models.Response{data=[]models.PriceIndexResponse}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/price-index/districts/{districtId} [get]
func (h *PriceIndexHandler) GetDistrictPriceIndex(c *gin.Context) {
	districtID, err := strconv.ParseUint(c.Param("districtId"), 10, 32)
	if err != nil {
		models.BadRequest(c, "invalid district id")
		return
	}
	
	var filter models.GetDistrictPriceIndexRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		models.BadRequest(c, err.Error())
		return
	}
	
	indices, err := h.service.GetDistrictPriceIndex(c.Request.Context(), uint(districtID), &filter)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}
	
	models.Success(c, indices)
}

// 4. GetEstatePriceIndex 获取屋苑楼价指数
// @Summary 获取屋苑楼价指数
// @Description 获取指定屋苑的历史楼价指数数据
// @Tags 楼价指数
// @Accept json
// @Produce json
// @Param estateId path int true "屋苑ID"
// @Param start_period query string false "开始周期 YYYY-MM"
// @Param end_period query string false "结束周期 YYYY-MM"
// @Param limit query int false "返回记录数" default(12)
// @Success 200 {object} models.Response{data=[]models.PriceIndexResponse}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/price-index/estates/{estateId} [get]
func (h *PriceIndexHandler) GetEstatePriceIndex(c *gin.Context) {
	estateID, err := strconv.ParseUint(c.Param("estateId"), 10, 32)
	if err != nil {
		models.BadRequest(c, "invalid estate id")
		return
	}
	
	var filter models.GetEstatePriceIndexRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		models.BadRequest(c, err.Error())
		return
	}
	
	indices, err := h.service.GetEstatePriceIndex(c.Request.Context(), uint(estateID), &filter)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}
	
	models.Success(c, indices)
}

// 5. GetPriceTrends 获取价格走势
// @Summary 获取价格走势
// @Description 获取指定条件的楼价走势数据，包含统计信息
// @Tags 楼价指数
// @Accept json
// @Produce json
// @Param index_type query string true "指数类型" Enums(overall, district, estate, property_type)
// @Param district_id query int false "地区ID（当index_type=district时必填）"
// @Param estate_id query int false "屋苑ID（当index_type=estate时必填）"
// @Param property_type query string false "物业类型"
// @Param start_period query string true "开始周期 YYYY-MM"
// @Param end_period query string true "结束周期 YYYY-MM"
// @Success 200 {object} models.Response{data=models.PriceTrendResponse}
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/price-index/trends [get]
func (h *PriceIndexHandler) GetPriceTrends(c *gin.Context) {
	var filter models.GetPriceTrendsRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		models.BadRequest(c, err.Error())
		return
	}
	
	trends, err := h.service.GetPriceTrends(c.Request.Context(), &filter)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			models.NotFound(c, "no trend data found")
			return
		}
		models.InternalError(c, err.Error())
		return
	}
	
	models.Success(c, trends)
}

// 6. ComparePriceIndex 对比楼价指数
// @Summary 对比楼价指数
// @Description 对比多个地区、屋苑或物业类型的楼价指数
// @Tags 楼价指数
// @Accept json
// @Produce json
// @Param compare_type query string true "对比类型" Enums(districts, estates, property_types)
// @Param district_ids query []int false "地区ID列表（当compare_type=districts时使用）"
// @Param estate_ids query []int false "屋苑ID列表（当compare_type=estates时使用）"
// @Param property_types query []string false "物业类型列表（当compare_type=property_types时使用）"
// @Param start_period query string true "开始周期 YYYY-MM"
// @Param end_period query string true "结束周期 YYYY-MM"
// @Success 200 {object} models.Response{data=models.ComparePriceIndexResponse}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/price-index/compare [get]
func (h *PriceIndexHandler) ComparePriceIndex(c *gin.Context) {
	var filter models.ComparePriceIndexRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		models.BadRequest(c, err.Error())
		return
	}
	
	result, err := h.service.ComparePriceIndex(c.Request.Context(), &filter)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}
	
	models.Success(c, result)
}

// 7. ExportPriceData 导出价格数据
// @Summary 导出价格数据
// @Description 导出楼价指数数据为CSV、JSON或Excel格式
// @Tags 楼价指数
// @Accept json
// @Produce json
// @Param index_type query string false "指数类型" Enums(overall, district, estate, property_type)
// @Param district_id query int false "地区ID"
// @Param estate_id query int false "屋苑ID"
// @Param property_type query string false "物业类型"
// @Param start_period query string true "开始周期 YYYY-MM"
// @Param end_period query string true "结束周期 YYYY-MM"
// @Param format query string false "导出格式" Enums(csv, json, excel) default(csv)
// @Success 200 {object} models.Response{data=models.ExportPriceDataResponse}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/price-index/export [get]
func (h *PriceIndexHandler) ExportPriceData(c *gin.Context) {
	var filter models.ExportPriceDataRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		models.BadRequest(c, err.Error())
		return
	}
	
	result, err := h.service.ExportPriceData(c.Request.Context(), &filter)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}
	
	models.Success(c, result)
}

// 8. GetPriceIndexHistory 获取历史楼价指数
// @Summary 获取历史楼价指数
// @Description 获取指定年数的历史楼价指数，包含年度统计
// @Tags 楼价指数
// @Accept json
// @Produce json
// @Param index_type query string true "指数类型" Enums(overall, district, estate, property_type)
// @Param district_id query int false "地区ID"
// @Param estate_id query int false "屋苑ID"
// @Param property_type query string false "物业类型"
// @Param years query int false "查询最近几年的数据" default(5)
// @Success 200 {object} models.Response{data=models.PriceIndexHistoryResponse}
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/price-index/history [get]
func (h *PriceIndexHandler) GetPriceIndexHistory(c *gin.Context) {
	var filter models.GetPriceIndexHistoryRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		models.BadRequest(c, err.Error())
		return
	}
	
	history, err := h.service.GetPriceIndexHistory(c.Request.Context(), &filter)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			models.NotFound(c, "no history data found")
			return
		}
		models.InternalError(c, err.Error())
		return
	}
	
	models.Success(c, history)
}

// 9. CreatePriceIndex 创建楼价指数
// @Summary 创建楼价指数
// @Description 创建新的楼价指数记录（需要管理员权限）
// @Tags 楼价指数
// @Accept json
// @Produce json
// @Param request body models.CreatePriceIndexRequest true "创建请求"
// @Success 201 {object} models.Response{data=models.CreatePriceIndexResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/price-index [post]
// @Security BearerAuth
func (h *PriceIndexHandler) CreatePriceIndex(c *gin.Context) {
	var req models.CreatePriceIndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}
	
	result, err := h.service.CreatePriceIndex(c.Request.Context(), &req)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}
	
	models.Created(c, result)
}

// 10. UpdatePriceIndex 更新楼价指数
// @Summary 更新楼价指数
// @Description 更新现有的楼价指数记录（需要管理员权限）
// @Tags 楼价指数
// @Accept json
// @Produce json
// @Param id path int true "楼价指数ID"
// @Param request body models.UpdatePriceIndexRequest true "更新请求"
// @Success 200 {object} models.Response{data=models.UpdatePriceIndexResponse}
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/v1/price-index/{id} [put]
// @Security BearerAuth
func (h *PriceIndexHandler) UpdatePriceIndex(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		models.BadRequest(c, "invalid price index id")
		return
	}
	
	var req models.UpdatePriceIndexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}
	
	result, err := h.service.UpdatePriceIndex(c.Request.Context(), uint(id), &req)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			models.NotFound(c, "price index not found")
			return
		}
		models.InternalError(c, err.Error())
		return
	}
	
	models.Success(c, result)
}
