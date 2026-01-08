package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// SchoolController Methods:
// 0. NewSchoolController(service *services.SchoolService) -> 注入 SchoolService
// 1. ListSchools(c *gin.Context) -> 学校列表
// 2. GetSchool(c *gin.Context) -> 学校详情
// 3. GetSchoolNet(c *gin.Context) -> 获取学校所属校网
// 4. SearchSchools(c *gin.Context) -> 搜索学校

type SchoolController struct {
	service *services.SchoolService
}

// 0. NewSchoolController 构造函数
func NewSchoolController(service *services.SchoolService) *SchoolController {
	return &SchoolController{service: service}
}

// 1. ListSchools 学校列表
// @Summary 学校列表
// @Tags School
// @Produce json
// @Param type query string false "学校类型 (primary, secondary)"
// @Param category query string false "学校类别 (government, aided, direct_subsidy, private, international)"
// @Param gender query string false "性别 (coed, boys, girls)"
// @Param school_net_id query int false "校网ID"
// @Param district_id query int false "地区ID"
// @Param keyword query string false "关键词搜索"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} tools.Response{data=models.PaginatedSchoolsResponse}
// @Router /api/v1/schools [get]
func (ctrl *SchoolController) ListSchools(c *gin.Context) {
	var req models.ListSchoolsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := ctrl.service.ListSchools(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 2. GetSchool 学校详情
// @Summary 学校详情
// @Tags School
// @Produce json
// @Param id path int true "学校ID"
// @Success 200 {object} tools.Response{data=models.SchoolDetailResponse}
// @Router /api/v1/schools/{id} [get]
func (ctrl *SchoolController) GetSchool(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid school id")
		return
	}

	school, err := ctrl.service.GetSchool(c.Request.Context(), uint(id))
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "school not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, school)
}

// 3. GetSchoolNet 获取学校所属校网
// @Summary 获取学校所属校网
// @Tags School
// @Produce json
// @Param id path int true "学校ID"
// @Success 200 {object} tools.Response{data=models.SchoolNetResponse}
// @Router /api/v1/schools/{id}/school-net [get]
func (ctrl *SchoolController) GetSchoolNet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid school id")
		return
	}

	schoolNet, err := ctrl.service.GetSchoolNet(c.Request.Context(), uint(id))
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "school not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, schoolNet)
}

// 4. SearchSchools 搜索学校
// @Summary 搜索学校
// @Tags School
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} tools.Response{data=models.PaginatedSchoolsResponse}
// @Router /api/v1/schools/search [get]
func (ctrl *SchoolController) SearchSchools(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		tools.BadRequest(c, "keyword is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	result, err := ctrl.service.SearchSchools(c.Request.Context(), keyword, page, pageSize)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}
