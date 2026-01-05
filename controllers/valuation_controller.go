package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/gin-gonic/gin"
)

// ValuationHandler Methods:
// 0. NewValuationHandler(service services.ValuationService) -> 注入 ValuationService
// 1. ListValuations(c *gin.Context) -> 获取屋苑估价列表
// 2. GetEstateValuation(c *gin.Context) -> 获取指定屋苑估价参考
// 3. SearchValuations(c *gin.Context) -> 搜索屋苑估价
// 4. GetDistrictValuations(c *gin.Context) -> 获取地区屋苑估价列表

type ValuationHandler struct {
	service services.ValuationService
}

// 0. NewValuationHandler -> 注入 ValuationService
func NewValuationHandler(service services.ValuationService) *ValuationHandler {
	return &ValuationHandler{service: service}
}

// 1. ListValuations -> 获取屋苑估价列表
// @Summary 获取屋苑估价列表
// @Description 获取所有屋苑的估价信息列表（支持筛选）
// @Tags 物业估价
// @Accept json
// @Produce json
// @Param district_id query uint false "地区ID"
// @Param min_price query number false "最低价格"
// @Param max_price query number false "最高价格"
// @Param min_area query number false "最小面积"
// @Param max_area query number false "最大面积"
// @Param school_net query string false "校网号"
// @Param sort_by query string false "排序字段(avg_price/name/completion_year)" default(avg_price)
// @Param sort_order query string false "排序顺序(asc/desc)" default(desc)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} models.Response{data=[]models.ValuationListItemResponse}
// @Router /valuation [get]
func (h *ValuationHandler) ListValuations(c *gin.Context) {
	var req models.ListValuationsRequest
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

	valuations, total, err := h.service.ListValuations(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.SuccessWithPagination(c, valuations, req.Page, req.PageSize, total)
}

// 2. GetEstateValuation -> 获取指定屋苑估价参考
// @Summary 获取指定屋苑估价参考
// @Description 根据屋苑ID获取该屋苑的详细估价信息
// @Tags 物业估价
// @Accept json
// @Produce json
// @Param estateId path int true "屋苑ID"
// @Success 200 {object} models.Response{data=models.ValuationResponse}
// @Failure 404 {object} models.Response
// @Router /valuation/{estateId} [get]
func (h *ValuationHandler) GetEstateValuation(c *gin.Context) {
	estateID, err := strconv.ParseUint(c.Param("estateId"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid estate id")
		return
	}

	valuation, err := h.service.GetEstateValuation(c.Request.Context(), uint(estateID))
	if err != nil {
		tools.NotFound(c, "estate valuation not found")
		return
	}

	tools.Success(c, valuation)
}

// 3. SearchValuations -> 搜索屋苑估价
// @Summary 搜索屋苑估价
// @Description 根据关键词搜索屋苑估价信息
// @Tags 物业估价
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词（屋苑名称或地址）"
// @Param district_id query uint false "地区ID"
// @Param min_price query number false "最低价格"
// @Param max_price query number false "最高价格"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} models.Response{data=[]models.ValuationListItemResponse}
// @Router /valuation/search [get]
func (h *ValuationHandler) SearchValuations(c *gin.Context) {
	var req models.SearchValuationsRequest
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

	valuations, total, err := h.service.SearchValuations(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.SuccessWithPagination(c, valuations, req.Page, req.PageSize, total)
}

// 4. GetDistrictValuations -> 获取地区屋苑估价列表
// @Summary 获取地区屋苑估价列表
// @Description 获取指定地区内所有屋苑的估价信息及地区统计数据
// @Tags 物业估价
// @Accept json
// @Produce json
// @Param districtId path int true "地区ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} models.Response{data=models.DistrictValuationResponse}
// @Failure 404 {object} models.Response
// @Router /valuation/districts/{districtId} [get]
func (h *ValuationHandler) GetDistrictValuations(c *gin.Context) {
	districtID, err := strconv.ParseUint(c.Param("districtId"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid district id")
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	valuation, err := h.service.GetDistrictValuations(c.Request.Context(), uint(districtID), page, pageSize)
	if err != nil {
		tools.NotFound(c, "district valuations not found")
		return
	}

	tools.Success(c, valuation)
}
