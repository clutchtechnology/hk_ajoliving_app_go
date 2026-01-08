package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// SchoolNetController Methods:
// 0. NewSchoolNetController(service *services.SchoolNetService) -> 注入 SchoolNetService
// 1. ListSchoolNets(c *gin.Context) -> 校网列表
// 2. GetSchoolNet(c *gin.Context) -> 校网详情
// 3. GetSchoolsInNet(c *gin.Context) -> 校网内学校
// 4. GetPropertiesInNet(c *gin.Context) -> 校网内房源
// 5. GetEstatesInNet(c *gin.Context) -> 校网内屋苑
// 6. SearchSchoolNets(c *gin.Context) -> 搜索校网

type SchoolNetController struct {
	service *services.SchoolNetService
}

// 0. NewSchoolNetController 构造函数
func NewSchoolNetController(service *services.SchoolNetService) *SchoolNetController {
	return &SchoolNetController{service: service}
}

// 1. ListSchoolNets 校网列表
// @Summary 校网列表
// @Tags SchoolNet
// @Produce json
// @Param type query string false "校网类型 (primary, secondary)"
// @Param district_id query int false "地区ID"
// @Param keyword query string false "关键词搜索"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} tools.Response{data=models.PaginatedSchoolNetsResponse}
// @Router /api/v1/school-nets [get]
func (ctrl *SchoolNetController) ListSchoolNets(c *gin.Context) {
	var req models.ListSchoolNetsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	result, err := ctrl.service.ListSchoolNets(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 2. GetSchoolNet 校网详情
// @Summary 校网详情
// @Tags SchoolNet
// @Produce json
// @Param id path int true "校网ID"
// @Success 200 {object} tools.Response{data=models.SchoolNetDetailResponse}
// @Router /api/v1/school-nets/{id} [get]
func (ctrl *SchoolNetController) GetSchoolNet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid school net id")
		return
	}

	schoolNet, err := ctrl.service.GetSchoolNet(c.Request.Context(), uint(id))
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "school net not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, schoolNet)
}

// 3. GetSchoolsInNet 校网内学校
// @Summary 校网内学校
// @Tags SchoolNet
// @Produce json
// @Param id path int true "校网ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} tools.Response{data=models.PaginatedSchoolsResponse}
// @Router /api/v1/school-nets/{id}/schools [get]
func (ctrl *SchoolNetController) GetSchoolsInNet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid school net id")
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

	result, err := ctrl.service.GetSchoolsInNet(c.Request.Context(), uint(id), page, pageSize)
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "school net not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 4. GetPropertiesInNet 校网内房源
// @Summary 校网内房源
// @Tags SchoolNet
// @Produce json
// @Param id path int true "校网ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} tools.Response
// @Router /api/v1/school-nets/{id}/properties [get]
func (ctrl *SchoolNetController) GetPropertiesInNet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid school net id")
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

	result, err := ctrl.service.GetPropertiesInNet(c.Request.Context(), uint(id), page, pageSize)
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "school net not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 5. GetEstatesInNet 校网内屋苑
// @Summary 校网内屋苑
// @Tags SchoolNet
// @Produce json
// @Param id path int true "校网ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} tools.Response
// @Router /api/v1/school-nets/{id}/estates [get]
func (ctrl *SchoolNetController) GetEstatesInNet(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid school net id")
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

	result, err := ctrl.service.GetEstatesInNet(c.Request.Context(), uint(id), page, pageSize)
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "school net not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}

// 6. SearchSchoolNets 搜索校网
// @Summary 搜索校网
// @Tags SchoolNet
// @Produce json
// @Param keyword query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} tools.Response{data=models.PaginatedSchoolNetsResponse}
// @Router /api/v1/school-nets/search [get]
func (ctrl *SchoolNetController) SearchSchoolNets(c *gin.Context) {
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

	result, err := ctrl.service.SearchSchoolNets(c.Request.Context(), keyword, page, pageSize)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, result)
}
