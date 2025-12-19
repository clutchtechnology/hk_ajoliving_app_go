package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/service"
)

// SchoolHandler 校网和学校处理器
type SchoolHandler struct {
	*BaseHandler
	service service.SchoolService
}

// NewSchoolHandler 创建校网和学校处理器
func NewSchoolHandler(service service.SchoolService) *SchoolHandler {
	return &SchoolHandler{
		BaseHandler: NewBaseHandler(),
		service:     service,
	}
}

// ListSchoolNets 获取校网列表
// @Summary 获取校网列表
// @Description 获取校网列表，支持地区和级别筛选
// @Tags 校网
// @Accept json
// @Produce json
// @Param district_id query int false "地区ID"
// @Param level query string false "级别" Enums(primary, secondary)
// @Param keyword query string false "关键词搜索"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(net_code)
// @Param sort_order query string false "排序方向" Enums(asc, desc) default(asc)
// @Success 200 {object} response.PaginatedResponse{data=[]response.SchoolNetListItemResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/school-nets [get]
func (h *SchoolHandler) ListSchoolNets(c *gin.Context) {
	var filter request.ListSchoolNetsRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	
	schoolNets, total, err := h.service.ListSchoolNets(c.Request.Context(), &filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	
	response.SuccessWithPagination(c, schoolNets, &response.Pagination{
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Total:     total,
		TotalPage: int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)),
	})
}

// GetSchoolNet 获取校网详情
// @Summary 获取校网详情
// @Description 根据校网ID获取详细信息
// @Tags 校网
// @Accept json
// @Produce json
// @Param id path int true "校网ID"
// @Success 200 {object} response.Response{data=response.SchoolNetResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/school-nets/{id} [get]
func (h *SchoolHandler) GetSchoolNet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid school net id")
		return
	}
	
	schoolNet, err := h.service.GetSchoolNet(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "school net not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, schoolNet)
}

// GetSchoolsInNet 获取校网内学校
// @Summary 获取校网内学校
// @Description 根据校网ID获取该校网内的所有学校
// @Tags 校网
// @Accept json
// @Produce json
// @Param id path int true "校网ID"
// @Success 200 {object} response.Response{data=[]response.SchoolListItemResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/school-nets/{id}/schools [get]
func (h *SchoolHandler) GetSchoolsInNet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid school net id")
		return
	}
	
	schools, err := h.service.GetSchoolsInNet(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "school net not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, schools)
}

// GetPropertiesInNet 获取校网内房源
// @Summary 获取校网内房源
// @Description 根据校网ID获取该校网所在地区的房源
// @Tags 校网
// @Accept json
// @Produce json
// @Param id path int true "校网ID"
// @Success 200 {object} response.Response{data=[]response.PropertyListItemResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/school-nets/{id}/properties [get]
func (h *SchoolHandler) GetPropertiesInNet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid school net id")
		return
	}
	
	properties, err := h.service.GetPropertiesInNet(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "school net not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, properties)
}

// GetEstatesInNet 获取校网内屋苑
// @Summary 获取校网内屋苑
// @Description 根据校网ID获取该校网所在地区的屋苑
// @Tags 校网
// @Accept json
// @Produce json
// @Param id path int true "校网ID"
// @Success 200 {object} response.Response{data=[]response.EstateListItemResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/school-nets/{id}/estates [get]
func (h *SchoolHandler) GetEstatesInNet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid school net id")
		return
	}
	
	estates, err := h.service.GetEstatesInNet(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "school net not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, estates)
}

// SearchSchoolNets 搜索校网
// @Summary 搜索校网
// @Description 根据关键词搜索校网
// @Tags 校网
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.PaginatedResponse{data=[]response.SchoolNetListItemResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/school-nets/search [get]
func (h *SchoolHandler) SearchSchoolNets(c *gin.Context) {
	var filter request.SearchSchoolNetsRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	
	schoolNets, total, err := h.service.SearchSchoolNets(c.Request.Context(), &filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	
	response.SuccessWithPagination(c, schoolNets, &response.Pagination{
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Total:     total,
		TotalPage: int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)),
	})
}

// ListSchools 获取学校列表
// @Summary 获取学校列表
// @Description 获取学校列表，支持多种筛选条件
// @Tags 学校
// @Accept json
// @Produce json
// @Param school_net_id query int false "校网ID"
// @Param district_id query int false "地区ID"
// @Param category query string false "类别" Enums(government, aided, direct_subsidy, private, international)
// @Param level query string false "级别" Enums(kindergarten, primary, secondary)
// @Param gender query string false "性别" Enums(co-ed, boys, girls)
// @Param keyword query string false "关键词搜索"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param sort_by query string false "排序字段" default(name_zh_hant)
// @Param sort_order query string false "排序方向" Enums(asc, desc) default(asc)
// @Success 200 {object} response.PaginatedResponse{data=[]response.SchoolListItemResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/schools [get]
func (h *SchoolHandler) ListSchools(c *gin.Context) {
	var filter request.ListSchoolsRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	
	schools, total, err := h.service.ListSchools(c.Request.Context(), &filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	
	response.SuccessWithPagination(c, schools, &response.Pagination{
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Total:     total,
		TotalPage: int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)),
	})
}

// GetSchoolNetBySchoolID 获取学校所属校网
// @Summary 获取学校所属校网
// @Description 根据学校ID获取该学校所属的校网信息
// @Tags 学校
// @Accept json
// @Produce json
// @Param id path int true "学校ID"
// @Success 200 {object} response.Response{data=response.SchoolNetResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/schools/{id}/school-net [get]
func (h *SchoolHandler) GetSchoolNetBySchoolID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid school id")
		return
	}
	
	schoolNet, err := h.service.GetSchoolNetBySchoolID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			response.NotFound(c, "school or school net not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}
	
	response.Success(c, schoolNet)
}

// SearchSchools 搜索学校
// @Summary 搜索学校
// @Description 根据关键词搜索学校
// @Tags 学校
// @Accept json
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.PaginatedResponse{data=[]response.SchoolListItemResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/schools/search [get]
func (h *SchoolHandler) SearchSchools(c *gin.Context) {
	var filter request.SearchSchoolsRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	
	schools, total, err := h.service.SearchSchools(c.Request.Context(), &filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	
	response.SuccessWithPagination(c, schools, &response.Pagination{
		Page:      filter.Page,
		PageSize:  filter.PageSize,
		Total:     total,
		TotalPage: int((total + int64(filter.PageSize) - 1) / int64(filter.PageSize)),
	})
}
