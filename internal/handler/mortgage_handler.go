package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/service"
)

// MortgageHandler 按揭处理器
type MortgageHandler struct {
	*BaseHandler
	service service.MortgageService
}

// NewMortgageHandler 创建按揭处理器
func NewMortgageHandler(service service.MortgageService) *MortgageHandler {
	return &MortgageHandler{
		BaseHandler: NewBaseHandler(),
		service:     service,
	}
}

// CalculateMortgage 计算按揭
// @Summary 计算按揭月供
// @Description 根据物业价格、首付、利率和还款期计算月供及还款计划
// @Tags 按揭
// @Accept json
// @Produce json
// @Param request body request.CalculateMortgageRequest true "计算参数"
// @Success 200 {object} response.Response{data=response.MortgageCalculationResponse}
// @Failure 400 {object} response.Response
// @Router /api/v1/mortgage/calculate [post]
func (h *MortgageHandler) CalculateMortgage(c *gin.Context) {
	var req request.CalculateMortgageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	
	result, err := h.service.CalculateMortgage(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, result)
}

// GetMortgageRates 获取按揭利率列表
// @Summary 获取按揭利率列表
// @Description 获取所有有效的按揭利率
// @Tags 按揭
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]response.MortgageRateResponse}
// @Failure 500 {object} response.Response
// @Router /api/v1/mortgage/rates [get]
func (h *MortgageHandler) GetMortgageRates(c *gin.Context) {
	rates, err := h.service.GetMortgageRates(c.Request.Context())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, rates)
}

// GetBankMortgageRate 获取指定银行的按揭利率
// @Summary 获取指定银行的按揭利率
// @Description 根据银行ID获取该银行的所有按揭利率方案
// @Tags 按揭
// @Accept json
// @Produce json
// @Param bank_id path int true "银行ID"
// @Success 200 {object} response.Response{data=[]response.MortgageRateResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/mortgage/rates/bank/{bank_id} [get]
func (h *MortgageHandler) GetBankMortgageRate(c *gin.Context) {
	bankID, err := strconv.ParseUint(c.Param("bank_id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid bank id")
		return
	}
	
	rates, err := h.service.GetBankMortgageRate(c.Request.Context(), uint(bankID))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "bank not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, rates)
}

// CompareMortgageRates 比较按揭利率
// @Summary 比较按揭利率
// @Description 比较不同银行的按揭利率和月供金额
// @Tags 按揭
// @Accept json
// @Produce json
// @Param request body request.CompareMortgageRatesRequest true "比较参数"
// @Success 200 {object} response.Response{data=response.MortgageRateComparisonResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/mortgage/rates/compare [post]
func (h *MortgageHandler) CompareMortgageRates(c *gin.Context) {
	var req request.CompareMortgageRatesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	
	result, err := h.service.CompareMortgageRates(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, result)
}

// ApplyMortgage 申请按揭
// @Summary 申请按揭
// @Description 提交按揭申请
// @Tags 按揭
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body request.ApplyMortgageRequest true "申请信息"
// @Success 200 {object} response.Response{data=response.MortgageApplicationResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/mortgage/apply [post]
func (h *MortgageHandler) ApplyMortgage(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "user not authenticated")
		return
	}
	
	var req request.ApplyMortgageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	
	application, err := h.service.ApplyMortgage(c.Request.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "bank not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, application)
}

// GetMortgageApplications 获取按揭申请列表
// @Summary 获取按揭申请列表
// @Description 获取当前用户的按揭申请列表
// @Tags 按揭
// @Accept json
// @Produce json
// @Security Bearer
// @Param status query string false "状态筛选" Enums(pending, under_review, approved, rejected, withdrawn)
// @Param bank_id query int false "银行ID筛选"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(created_at)
// @Param sort_order query string false "排序方向" Enums(asc, desc) default(desc)
// @Success 200 {object} response.PaginatedResponse{data=[]response.MortgageApplicationResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/mortgage/applications [get]
func (h *MortgageHandler) GetMortgageApplications(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "user not authenticated")
		return
	}
	
	var filter request.ListMortgageApplicationsRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	
	applications, total, err := h.service.GetMortgageApplications(c.Request.Context(), userID, &filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	
	response.SuccessWithPagination(c, applications, &response.Pagination{
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Total:     total,
		TotalPage: int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)),
	})
}

// GetMortgageApplication 获取按揭申请详情
// @Summary 获取按揭申请详情
// @Description 根据申请ID获取详细信息
// @Tags 按揭
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "申请ID"
// @Success 200 {object} response.Response{data=response.MortgageApplicationResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/mortgage/applications/{id} [get]
func (h *MortgageHandler) GetMortgageApplication(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "user not authenticated")
		return
	}
	
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid application id")
		return
	}
	
	application, err := h.service.GetMortgageApplication(c.Request.Context(), userID, uint(id))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "application not found")
			return
		}
		if errors.Is(err, errors.ErrForbidden) {
			response.Forbidden(c, "you don't have permission to access this application")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, application)
}
