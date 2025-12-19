package handler

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/service"
	"github.com/gin-gonic/gin"
)

// EstateHandler Methods:
// 0. NewEstateHandler(service service.EstateService) -> 注入 EstateService
// 1. ListEstates(c *gin.Context) -> 获取屋苑列表
// 2. GetEstate(c *gin.Context) -> 获取单个屋苑详情
// 3. GetEstateProperties(c *gin.Context) -> 获取屋苑内的房源列表
// 4. GetEstateStatistics(c *gin.Context) -> 获取屋苑统计数据
// 5. GetFeaturedEstates(c *gin.Context) -> 获取精选屋苑
// 6. CreateEstate(c *gin.Context) -> 创建屋苑（需要认证）
// 7. UpdateEstate(c *gin.Context) -> 更新屋苑（需要认证）
// 8. DeleteEstate(c *gin.Context) -> 删除屋苑（需要认证）

type EstateHandler struct {
	service service.EstateService
}

// 0. NewEstateHandler -> 注入 EstateService
func NewEstateHandler(service service.EstateService) *EstateHandler {
	return &EstateHandler{service: service}
}

// 1. ListEstates -> 获取屋苑列表
// @Summary 获取屋苑列表
// @Description 根据筛选条件获取屋苑列表
// @Tags 屋苑
// @Accept json
// @Produce json
// @Param district_id query uint false "地区ID"
// @Param school_net query string false "校网号"
// @Param min_completion_year query int false "最早落成年份"
// @Param max_completion_year query int false "最晚落成年份"
// @Param min_avg_price query number false "最低平均成交价"
// @Param max_avg_price query number false "最高平均成交价"
// @Param has_listings query bool false "是否有房源"
// @Param has_transactions query bool false "是否有成交记录"
// @Param is_featured query bool false "是否精选"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(created_at)
// @Param sort_order query string false "排序顺序(asc/desc)" default(desc)
// @Success 200 {object} response.Response{data=[]response.EstateListItemResponse}
// @Router /estates [get]
func (h *EstateHandler) ListEstates(c *gin.Context) {
	var req request.ListEstatesRequest
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

	estates, total, err := h.service.ListEstates(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithPagination(c, estates, req.Page, req.PageSize, total)
}

// 2. GetEstate -> 获取单个屋苑详情
// @Summary 获取屋苑详情
// @Description 根据ID获取屋苑详细信息
// @Tags 屋苑
// @Accept json
// @Produce json
// @Param id path int true "屋苑ID"
// @Success 200 {object} response.Response{data=response.EstateResponse}
// @Failure 404 {object} response.Response
// @Router /estates/{id} [get]
func (h *EstateHandler) GetEstate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid estate id")
		return
	}

	estate, err := h.service.GetEstate(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "estate not found")
		return
	}

	response.Success(c, estate)
}

// 3. GetEstateProperties -> 获取屋苑内的房源列表
// @Summary 获取屋苑房源列表
// @Description 获取指定屋苑内的房源列表
// @Tags 屋苑
// @Accept json
// @Produce json
// @Param id path int true "屋苑ID"
// @Param listing_type query string false "房源类型(sale/rent)"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=[]model.Property}
// @Failure 404 {object} response.Response
// @Router /estates/{id}/properties [get]
func (h *EstateHandler) GetEstateProperties(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid estate id")
		return
	}

	listingType := c.Query("listing_type")
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

	properties, total, err := h.service.GetEstateProperties(c.Request.Context(), uint(id), listingType, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithPagination(c, properties, page, pageSize, total)
}

// 4. GetEstateStatistics -> 获取屋苑统计数据
// @Summary 获取屋苑统计数据
// @Description 获取屋苑的成交统计、在售在租数量等数据
// @Tags 屋苑
// @Accept json
// @Produce json
// @Param id path int true "屋苑ID"
// @Success 200 {object} response.Response{data=response.EstateStatisticsResponse}
// @Failure 404 {object} response.Response
// @Router /estates/{id}/statistics [get]
func (h *EstateHandler) GetEstateStatistics(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid estate id")
		return
	}

	statistics, err := h.service.GetEstateStatistics(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "estate not found")
		return
	}

	response.Success(c, statistics)
}

// 5. GetFeaturedEstates -> 获取精选屋苑
// @Summary 获取精选屋苑
// @Description 获取平台推荐的精选屋苑列表
// @Tags 屋苑
// @Accept json
// @Produce json
// @Param limit query int false "返回数量" default(10)
// @Success 200 {object} response.Response{data=[]response.EstateListItemResponse}
// @Router /estates/featured [get]
func (h *EstateHandler) GetFeaturedEstates(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	estates, err := h.service.GetFeaturedEstates(c.Request.Context(), limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, estates)
}

// 6. CreateEstate -> 创建屋苑（需要认证）
// @Summary 创建屋苑
// @Description 创建新的屋苑（需要管理员权限）
// @Tags 屋苑
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body request.CreateEstateRequest true "屋苑信息"
// @Success 201 {object} response.Response{data=response.EstateResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /estates [post]
func (h *EstateHandler) CreateEstate(c *gin.Context) {
	var req request.CreateEstateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	estate, err := h.service.CreateEstate(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, estate)
}

// 7. UpdateEstate -> 更新屋苑（需要认证）
// @Summary 更新屋苑
// @Description 更新屋苑信息（需要管理员权限）
// @Tags 屋苑
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "屋苑ID"
// @Param body body request.UpdateEstateRequest true "更新信息"
// @Success 200 {object} response.Response{data=response.EstateResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /estates/{id} [put]
func (h *EstateHandler) UpdateEstate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid estate id")
		return
	}

	var req request.UpdateEstateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	estate, err := h.service.UpdateEstate(c.Request.Context(), uint(id), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, estate)
}

// 8. DeleteEstate -> 删除屋苑（需要认证）
// @Summary 删除屋苑
// @Description 删除屋苑（需要管理员权限）
// @Tags 屋苑
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "屋苑ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /estates/{id} [delete]
func (h *EstateHandler) DeleteEstate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid estate id")
		return
	}

	if err := h.service.DeleteEstate(c.Request.Context(), uint(id)); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "estate deleted successfully"})
}
