package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// ValuationController Methods:
// 0. NewValuationController(service *services.ValuationService) -> 注入 ValuationService
// 1. ListValuations(c *gin.Context) -> 获取屋苑估价列表
// 2. GetEstateValuation(c *gin.Context) -> 获取指定屋苑估价参考
// 3. SearchValuations(c *gin.Context) -> 搜索屋苑估价
// 4. GetDistrictValuations(c *gin.Context) -> 获取地区屋苑估价列表

type ValuationController struct {
	valuationService *services.ValuationService
}

// 0. NewValuationController -> 注入 ValuationService
func NewValuationController(valuationService *services.ValuationService) *ValuationController {
	return &ValuationController{
		valuationService: valuationService,
	}
}

// 1. ListValuations -> 获取屋苑估价列表
func (ctrl *ValuationController) ListValuations(c *gin.Context) {
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

	response, err := ctrl.valuationService.ListValuations(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, response)
}

// 2. GetEstateValuation -> 获取指定屋苑估价参考
func (ctrl *ValuationController) GetEstateValuation(c *gin.Context) {
	estateID, err := strconv.ParseUint(c.Param("estateId"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid estate id")
		return
	}

	valuation, err := ctrl.valuationService.GetEstateValuation(c.Request.Context(), uint(estateID))
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "estate not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, valuation)
}

// 3. SearchValuations -> 搜索屋苑估价
func (ctrl *ValuationController) SearchValuations(c *gin.Context) {
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

	response, err := ctrl.valuationService.SearchValuations(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, response)
}

// 4. GetDistrictValuations -> 获取地区屋苑估价列表
func (ctrl *ValuationController) GetDistrictValuations(c *gin.Context) {
	districtID, err := strconv.ParseUint(c.Param("districtId"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid district id")
		return
	}

	summary, err := ctrl.valuationService.GetDistrictValuations(c.Request.Context(), uint(districtID))
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "district not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, summary)
}
