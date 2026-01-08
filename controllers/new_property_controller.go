package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// NewDevelopmentController 新盘控制器
// Methods:
// 1. ListNewDevelopments(c *gin.Context) -> 获取新盘列表
// 2. GetNewDevelopment(c *gin.Context) -> 获取新盘详情
// 3. GetDevelopmentLayouts(c *gin.Context) -> 获取新盘户型列表
type NewDevelopmentController struct {
	service *services.NewDevelopmentService
}

// NewNewDevelopmentController 创建新盘控制器
func NewNewDevelopmentController(service *services.NewDevelopmentService) *NewDevelopmentController {
	return &NewDevelopmentController{service: service}
}

// ListNewDevelopments 获取新盘列表
// @Summary 获取新盘列表
// @Tags 新盘
// @Accept json
// @Produce json
// @Param district_id query int false "地区ID"
// @Param status query string false "状态: upcoming, presale, selling, sold_out, completed"
// @Param developer query string false "发展商名称（模糊搜索）"
// @Param primary_school_net query string false "小学校网"
// @Param secondary_school_net query string false "中学校网"
// @Param is_featured query bool false "是否精选"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} tools.Response{data=models.PaginatedNewPropertiesResponse}
// @Router /api/v1/new-properties [get]
func (ctrl *NewDevelopmentController) ListNewDevelopments(c *gin.Context) {
	var filter models.ListNewPropertiesRequest
	if err := c.ShouldBindQuery(&filter); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := ctrl.service.ListNewProperties(c.Request.Context(), &filter)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// GetNewDevelopment 获取新盘详情
// @Summary 获取新盘详情
// @Tags 新盘
// @Accept json
// @Produce json
// @Param id path int true "新盘ID"
// @Success 200 {object} tools.Response{data=models.NewPropertyDetailResponse}
// @Failure 404 {object} tools.Response
// @Router /api/v1/new-properties/{id} [get]
func (ctrl *NewDevelopmentController) GetNewDevelopment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid new property id")
		return
	}

	newProperty, err := ctrl.service.GetNewProperty(c.Request.Context(), uint(id))
	if err != nil {
		tools.NotFound(c, err.Error())
		return
	}

	tools.Success(c, newProperty)
}

// GetDevelopmentLayouts 获取新盘户型列表
// @Summary 获取新盘户型列表
// @Tags 新盘
// @Accept json
// @Produce json
// @Param id path int true "新盘ID"
// @Success 200 {object} tools.Response{data=[]models.NewPropertyLayout}
// @Failure 404 {object} tools.Response
// @Router /api/v1/new-properties/{id}/layouts [get]
func (ctrl *NewDevelopmentController) GetDevelopmentLayouts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid new property id")
		return
	}

	layouts, err := ctrl.service.GetNewPropertyLayouts(c.Request.Context(), uint(id))
	if err != nil {
		tools.NotFound(c, err.Error())
		return
	}

	tools.Success(c, layouts)
}
