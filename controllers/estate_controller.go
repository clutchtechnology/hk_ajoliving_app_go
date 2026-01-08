package controllers

import (
	"strconv"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/gin-gonic/gin"
)

// EstateController Methods:
// 0. NewEstateController(service *services.EstateService) -> 注入 EstateService
// 1. ListEstates(c *gin.Context) -> 获取屋苑列表
// 2. GetEstate(c *gin.Context) -> 获取屋苑详情
// 3. GetEstateProperties(c *gin.Context) -> 获取屋苑内房源列表
// 4. GetEstateImages(c *gin.Context) -> 获取屋苑图片
// 5. GetEstateFacilities(c *gin.Context) -> 获取屋苑设施
// 6. GetEstateTransactions(c *gin.Context) -> 获取屋苑成交记录
// 7. GetEstateStatistics(c *gin.Context) -> 获取屋苑统计数据
// 8. GetFeaturedEstates(c *gin.Context) -> 获取精选屋苑
// 9. CreateEstate(c *gin.Context) -> 创建屋苑（需认证）
// 10. UpdateEstate(c *gin.Context) -> 更新屋苑（需认证）
// 11. DeleteEstate(c *gin.Context) -> 删除屋苑（需认证）

type EstateController struct {
	estateService *services.EstateService
}

// 0. NewEstateController -> 注入 EstateService
func NewEstateController(estateService *services.EstateService) *EstateController {
	return &EstateController{
		estateService: estateService,
	}
}

// 1. ListEstates -> 获取屋苑列表
func (ctrl *EstateController) ListEstates(c *gin.Context) {
	var req models.ListEstatesRequest
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

	response, err := ctrl.estateService.ListEstates(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, response)
}

// 2. GetEstate -> 获取屋苑详情
func (ctrl *EstateController) GetEstate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid estate id")
		return
	}

	estate, err := ctrl.estateService.GetEstate(c.Request.Context(), uint(id))
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "estate not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, estate)
}

// 3. GetEstateProperties -> 获取屋苑内房源列表
func (ctrl *EstateController) GetEstateProperties(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid estate id")
		return
	}

	var req models.GetEstatePropertiesRequest
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

	properties, err := ctrl.estateService.GetEstateProperties(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "estate not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, properties)
}

// 4. GetEstateImages -> 获取屋苑图片
func (ctrl *EstateController) GetEstateImages(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid estate id")
		return
	}

	images, err := ctrl.estateService.GetEstateImages(c.Request.Context(), uint(id))
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "estate not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, images)
}

// 5. GetEstateFacilities -> 获取屋苑设施
func (ctrl *EstateController) GetEstateFacilities(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid estate id")
		return
	}

	facilities, err := ctrl.estateService.GetEstateFacilities(c.Request.Context(), uint(id))
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "estate not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, facilities)
}

// 6. GetEstateTransactions -> 获取屋苑成交记录
func (ctrl *EstateController) GetEstateTransactions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid estate id")
		return
	}

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	transactions, err := ctrl.estateService.GetEstateTransactions(c.Request.Context(), uint(id), page, pageSize)
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "estate not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, transactions)
}

// 7. GetEstateStatistics -> 获取屋苑统计数据
func (ctrl *EstateController) GetEstateStatistics(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid estate id")
		return
	}

	stats, err := ctrl.estateService.GetEstateStatistics(c.Request.Context(), uint(id))
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "estate not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, stats)
}

// 8. GetFeaturedEstates -> 获取精选屋苑
func (ctrl *EstateController) GetFeaturedEstates(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	estates, err := ctrl.estateService.GetFeaturedEstates(c.Request.Context(), limit)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, estates)
}

// 9. CreateEstate -> 创建屋苑（需认证）
func (ctrl *EstateController) CreateEstate(c *gin.Context) {
	var req models.CreateEstateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	estate, err := ctrl.estateService.CreateEstate(c.Request.Context(), &req)
	if err != nil {
		tools.InternalError(c, err.Error())
		return
	}

	tools.Created(c, estate)
}

// 10. UpdateEstate -> 更新屋苑（需认证）
func (ctrl *EstateController) UpdateEstate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid estate id")
		return
	}

	var req models.UpdateEstateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		tools.BadRequest(c, err.Error())
		return
	}

	estate, err := ctrl.estateService.UpdateEstate(c.Request.Context(), uint(id), &req)
	if err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "estate not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, estate)
}

// 11. DeleteEstate -> 删除屋苑（需认证）
func (ctrl *EstateController) DeleteEstate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		tools.BadRequest(c, "invalid estate id")
		return
	}

	if err := ctrl.estateService.DeleteEstate(c.Request.Context(), uint(id)); err != nil {
		if err == tools.ErrNotFound {
			tools.NotFound(c, "estate not found")
			return
		}
		tools.InternalError(c, err.Error())
		return
	}

	tools.Success(c, gin.H{"message": "estate deleted successfully"})
}
