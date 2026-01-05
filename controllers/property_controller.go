package controllers

// PropertyHandler Methods:
// 0. NewPropertyHandler(propertyService *service.PropertyService) -> 注入 PropertyService
// 1. ListProperties(c *gin.Context) -> 房产列表
// 2. GetProperty(c *gin.Context) -> 房产详情
// 3. CreateProperty(c *gin.Context) -> 创建房产
// 4. UpdateProperty(c *gin.Context) -> 更新房产
// 5. DeleteProperty(c *gin.Context) -> 删除房产
// 6. GetSimilarProperties(c *gin.Context) -> 相似房源
// 7. GetFeaturedProperties(c *gin.Context) -> 精选房源
// 8. GetHotProperties(c *gin.Context) -> 热门房源

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
)

// PropertyHandlerInterface 房产处理器接口
type PropertyHandlerInterface interface {
	ListProperties(c *gin.Context)       // 1. 房产列表
	GetProperty(c *gin.Context)          // 2. 房产详情
	CreateProperty(c *gin.Context)       // 3. 创建房产
	UpdateProperty(c *gin.Context)       // 4. 更新房产
	DeleteProperty(c *gin.Context)       // 5. 删除房产
	GetSimilarProperties(c *gin.Context) // 6. 相似房源
	GetFeaturedProperties(c *gin.Context) // 7. 精选房源
	GetHotProperties(c *gin.Context)     // 8. 热门房源
}

// PropertyHandler 房产处理器
type PropertyHandler struct {
	propertyService *service.PropertyService
}

// 0. NewPropertyHandler 注入 PropertyService
func NewPropertyHandler(propertyService *service.PropertyService) *PropertyHandler {
	return &PropertyHandler{
		propertyService: propertyService,
	}
}

// 1. ListProperties 房产列表
// ListProperties godoc
// @Summary      房产列表
// @Description  获取房产列表，支持多种筛选条件
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        district_id    query     int     false  "地区ID"
// @Param        building_name  query     string  false  "楼盘名称"
// @Param        min_price      query     number  false  "最低价格"
// @Param        max_price      query     number  false  "最高价格"
// @Param        min_area       query     number  false  "最小面积"
// @Param        max_area       query     number  false  "最大面积"
// @Param        bedrooms       query     int     false  "卧室数"
// @Param        property_type  query     string  false  "物业类型"
// @Param        listing_type   query     string  false  "房源类型(rent/sale)"
// @Param        status         query     string  false  "状态"
// @Param        school_net     query     string  false  "校网"
// @Param        sort_by        query     string  false  "排序字段"   default(created_at)
// @Param        sort_order     query     string  false  "排序方向"   default(desc)
// @Param        page           query     int     false  "页码"       default(1)
// @Param        page_size      query     int     false  "每页数量"   default(20)
// @Success      200            {object}  models.PaginatedResponse{data=[]models.PropertyListItemResponse}
// @Failure      400            {object}  models.Response
// @Failure      500            {object}  models.Response
// @Router       /api/v1/properties [get]
func (h *PropertyHandler) ListProperties(c *gin.Context) {
	var req models.ListPropertiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	properties, total, err := h.propertyService.ListProperties(c.Request.Context(), &req)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.SuccessWithPagination(c, properties, models.NewPagination(req.Page, req.PageSize, total))
}

// 2. GetProperty 房产详情
// GetProperty godoc
// @Summary      房产详情
// @Description  获取房产的详细信息
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "房产ID"
// @Success      200  {object}  models.Response{data=models.PropertyResponse}
// @Failure      400  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/v1/properties/{id} [get]
func (h *PropertyHandler) GetProperty(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		models.BadRequest(c, "Invalid property ID")
		return
	}

	property, err := h.propertyService.GetProperty(c.Request.Context(), uint(id))
	if err != nil {
		if err == service.ErrPropertyNotFound {
			models.NotFound(c, "Property not found")
			return
		}
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, property)
}

// 3. CreateProperty 创建房产
// CreateProperty godoc
// @Summary      创建房产
// @Description  创建新的房产信息
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      models.CreatePropertyRequest  true  "房产信息"
// @Success      201   {object}  models.Response{data=models.CreatePropertyResponse}
// @Failure      400   {object}  models.Response
// @Failure      401   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/properties [post]
func (h *PropertyHandler) CreateProperty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		models.Unauthorized(c, "Unauthorized")
		return
	}

	var req models.CreatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	result, err := h.propertyService.CreateProperty(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.Created(c, result)
}

// 4. UpdateProperty 更新房产
// UpdateProperty godoc
// @Summary      更新房产
// @Description  更新房产信息
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id    path      int                            true  "房产ID"
// @Param        body  body      models.UpdatePropertyRequest  true  "房产信息"
// @Success      200   {object}  models.Response{data=models.UpdatePropertyResponse}
// @Failure      400   {object}  models.Response
// @Failure      401   {object}  models.Response
// @Failure      403   {object}  models.Response
// @Failure      404   {object}  models.Response
// @Failure      500   {object}  models.Response
// @Router       /api/v1/properties/{id} [put]
func (h *PropertyHandler) UpdateProperty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		models.Unauthorized(c, "Unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		models.BadRequest(c, "Invalid property ID")
		return
	}

	var req models.UpdatePropertyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	result, err := h.propertyService.UpdateProperty(c.Request.Context(), userID.(uint), uint(id), &req)
	if err != nil {
		if err == service.ErrPropertyNotFound {
			models.NotFound(c, "Property not found")
			return
		}
		if err == service.ErrNotPropertyOwner {
			models.Forbidden(c, "You are not the owner of this property")
			return
		}
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, result)
}

// 5. DeleteProperty 删除房产
// DeleteProperty godoc
// @Summary      删除房产
// @Description  删除指定的房产
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "房产ID"
// @Success      200  {object}  models.Response
// @Failure      400  {object}  models.Response
// @Failure      401  {object}  models.Response
// @Failure      403  {object}  models.Response
// @Failure      404  {object}  models.Response
// @Failure      500  {object}  models.Response
// @Router       /api/v1/properties/{id} [delete]
func (h *PropertyHandler) DeleteProperty(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		models.Unauthorized(c, "Unauthorized")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		models.BadRequest(c, "Invalid property ID")
		return
	}

	err = h.propertyService.DeleteProperty(c.Request.Context(), userID.(uint), uint(id))
	if err != nil {
		if err == service.ErrPropertyNotFound {
			models.NotFound(c, "Property not found")
			return
		}
		if err == service.ErrNotPropertyOwner {
			models.Forbidden(c, "You are not the owner of this property")
			return
		}
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, gin.H{
		"message": "Property deleted successfully",
	})
}

// 6. GetSimilarProperties 相似房源
// GetSimilarProperties godoc
// @Summary      相似房源
// @Description  获取与指定房产相似的房源列表
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        id     path      int  true   "房产ID"
// @Param        limit  query     int  false  "数量限制"  default(6)
// @Success      200    {object}  models.Response{data=[]models.PropertyListItemResponse}
// @Failure      400    {object}  models.Response
// @Failure      404    {object}  models.Response
// @Failure      500    {object}  models.Response
// @Router       /api/v1/properties/{id}/similar [get]
func (h *PropertyHandler) GetSimilarProperties(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		models.BadRequest(c, "Invalid property ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "6"))
	if limit < 1 || limit > 20 {
		limit = 6
	}

	properties, err := h.propertyService.GetSimilarProperties(c.Request.Context(), uint(id), limit)
	if err != nil {
		if err == service.ErrPropertyNotFound {
			models.NotFound(c, "Property not found")
			return
		}
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, properties)
}

// 7. GetFeaturedProperties 精选房源
// GetFeaturedProperties godoc
// @Summary      精选房源
// @Description  获取精选房源列表
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        listing_type  query     string  false  "房源类型(rent/sale)"
// @Param        limit         query     int     false  "数量限制"  default(10)
// @Success      200           {object}  models.Response{data=[]models.PropertyListItemResponse}
// @Failure      500           {object}  models.Response
// @Router       /api/v1/properties/featured [get]
func (h *PropertyHandler) GetFeaturedProperties(c *gin.Context) {
	listingType := c.Query("listing_type")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	properties, err := h.propertyService.GetFeaturedProperties(c.Request.Context(), listingType, limit)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, properties)
}

// 8. GetHotProperties 热门房源
// GetHotProperties godoc
// @Summary      热门房源
// @Description  获取热门房源列表
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        listing_type  query     string  false  "房源类型(rent/sale)"
// @Param        limit         query     int     false  "数量限制"  default(10)
// @Success      200           {object}  models.Response{data=[]models.PropertyListItemResponse}
// @Failure      500           {object}  models.Response
// @Router       /api/v1/properties/hot [get]
func (h *PropertyHandler) GetHotProperties(c *gin.Context) {
	listingType := c.Query("listing_type")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}

	properties, err := h.propertyService.GetHotProperties(c.Request.Context(), listingType, limit)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.Success(c, properties)
}

// 9. ListBuyProperties 买房房源列表
// ListBuyProperties godoc
// @Summary      买房房源列表
// @Description  获取所有出售的房产列表
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        district_id    query     int     false  "地区ID"
// @Param        building_name  query     string  false  "楼盘名称"
// @Param        min_price      query     number  false  "最低价格"
// @Param        max_price      query     number  false  "最高价格"
// @Param        min_area       query     number  false  "最小面积"
// @Param        max_area       query     number  false  "最大面积"
// @Param        bedrooms       query     int     false  "卧室数"
// @Param        property_type  query     string  false  "物业类型"
// @Param        school_net     query     string  false  "校网"
// @Param        is_new         query     bool    false  "是否新房"
// @Param        sort_by        query     string  false  "排序字段"   default(created_at)
// @Param        sort_order     query     string  false  "排序方向"   default(desc)
// @Param        page           query     int     false  "页码"       default(1)
// @Param        page_size      query     int     false  "每页数量"   default(20)
// @Success      200            {object}  models.PaginatedResponse{data=[]models.PropertyListItemResponse}
// @Router       /api/v1/properties/buy [get]
func (h *PropertyHandler) ListBuyProperties(c *gin.Context) {
	var req models.ListBuyPropertiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	// 设置默认值
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	// 转换为通用列表请求
	listReq := models.ListPropertiesRequest{
		DistrictID:   req.DistrictID,
		BuildingName: req.BuildingName,
		MinPrice:     req.MinPrice,
		MaxPrice:     req.MaxPrice,
		MinArea:      req.MinArea,
		MaxArea:      req.MaxArea,
		Bedrooms:     req.Bedrooms,
		PropertyType: req.PropertyType,
		SchoolNet:    req.SchoolNet,
		ListingType:  "sale",
		SortBy:       req.SortBy,
		SortOrder:    req.SortOrder,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}

	properties, total, err := h.propertyService.ListProperties(c.Request.Context(), &listReq)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.SuccessWithPagination(c, properties, models.NewPagination(req.Page, req.PageSize, total))
}

// 10. ListNewProperties 新房列表
// ListNewProperties godoc
// @Summary      新房列表
// @Description  获取新房（未入伙）房产列表
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        district_id    query     int     false  "地区ID"
// @Param        min_price      query     number  false  "最低价格"
// @Param        max_price      query     number  false  "最高价格"
// @Param        bedrooms       query     int     false  "卧室数"
// @Param        page           query     int     false  "页码"       default(1)
// @Param        page_size      query     int     false  "每页数量"   default(20)
// @Success      200            {object}  models.PaginatedResponse{data=[]models.PropertyListItemResponse}
// @Router       /api/v1/properties/buy/new [get]
func (h *PropertyHandler) ListNewProperties(c *gin.Context) {
	var req models.ListBuyPropertiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	isNew := true
	listReq := models.ListPropertiesRequest{
		DistrictID:  req.DistrictID,
		MinPrice:    req.MinPrice,
		MaxPrice:    req.MaxPrice,
		Bedrooms:    req.Bedrooms,
		ListingType: "sale",
		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}
	_ = isNew // TODO: 需要在 Property 表添加 is_new 字段或通过其他方式筛选新房

	properties, total, err := h.propertyService.ListProperties(c.Request.Context(), &listReq)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.SuccessWithPagination(c, properties, models.NewPagination(req.Page, req.PageSize, total))
}

// 11. ListSecondhandProperties 二手房列表
// ListSecondhandProperties godoc
// @Summary      二手房列表
// @Description  获取二手房（已入伙）房产列表
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        district_id    query     int     false  "地区ID"
// @Param        min_price      query     number  false  "最低价格"
// @Param        max_price      query     number  false  "最高价格"
// @Param        bedrooms       query     int     false  "卧室数"
// @Param        page           query     int     false  "页码"       default(1)
// @Param        page_size      query     int     false  "每页数量"   default(20)
// @Success      200            {object}  models.PaginatedResponse{data=[]models.PropertyListItemResponse}
// @Router       /api/v1/properties/buy/secondhand [get]
func (h *PropertyHandler) ListSecondhandProperties(c *gin.Context) {
	var req models.ListBuyPropertiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	listReq := models.ListPropertiesRequest{
		DistrictID:  req.DistrictID,
		MinPrice:    req.MinPrice,
		MaxPrice:    req.MaxPrice,
		Bedrooms:    req.Bedrooms,
		ListingType: "sale",
		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}
	// TODO: 需要在 Property 表添加 is_new 字段，这里筛选 is_new = false 的房源

	properties, total, err := h.propertyService.ListProperties(c.Request.Context(), &listReq)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.SuccessWithPagination(c, properties, models.NewPagination(req.Page, req.PageSize, total))
}

// 12. ListRentProperties 租房房源列表
// ListRentProperties godoc
// @Summary      租房房源列表
// @Description  获取所有出租的房产列表
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        district_id    query     int     false  "地区ID"
// @Param        building_name  query     string  false  "楼盘名称"
// @Param        min_price      query     number  false  "最低月租"
// @Param        max_price      query     number  false  "最高月租"
// @Param        min_area       query     number  false  "最小面积"
// @Param        max_area       query     number  false  "最大面积"
// @Param        bedrooms       query     int     false  "卧室数"
// @Param        property_type  query     string  false  "物业类型"
// @Param        school_net     query     string  false  "校网"
// @Param        rent_type      query     string  false  "租期类型：short/long"
// @Param        sort_by        query     string  false  "排序字段"   default(created_at)
// @Param        sort_order     query     string  false  "排序方向"   default(desc)
// @Param        page           query     int     false  "页码"       default(1)
// @Param        page_size      query     int     false  "每页数量"   default(20)
// @Success      200            {object}  models.PaginatedResponse{data=[]models.PropertyListItemResponse}
// @Router       /api/v1/properties/rent [get]
func (h *PropertyHandler) ListRentProperties(c *gin.Context) {
	var req models.ListRentPropertiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	listReq := models.ListPropertiesRequest{
		DistrictID:   req.DistrictID,
		BuildingName: req.BuildingName,
		MinPrice:     req.MinPrice,
		MaxPrice:     req.MaxPrice,
		MinArea:      req.MinArea,
		MaxArea:      req.MaxArea,
		Bedrooms:     req.Bedrooms,
		PropertyType: req.PropertyType,
		SchoolNet:    req.SchoolNet,
		ListingType:  "rent",
		SortBy:       req.SortBy,
		SortOrder:    req.SortOrder,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}

	properties, total, err := h.propertyService.ListProperties(c.Request.Context(), &listReq)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.SuccessWithPagination(c, properties, models.NewPagination(req.Page, req.PageSize, total))
}

// 13. ListShortTermRent 短租房源
// ListShortTermRent godoc
// @Summary      短租房源
// @Description  获取短租（<6个月）房产列表
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        district_id    query     int     false  "地区ID"
// @Param        min_price      query     number  false  "最低月租"
// @Param        max_price      query     number  false  "最高月租"
// @Param        bedrooms       query     int     false  "卧室数"
// @Param        page           query     int     false  "页码"       default(1)
// @Param        page_size      query     int     false  "每页数量"   default(20)
// @Success      200            {object}  models.PaginatedResponse{data=[]models.PropertyListItemResponse}
// @Router       /api/v1/properties/rent/short-term [get]
func (h *PropertyHandler) ListShortTermRent(c *gin.Context) {
	var req models.ListRentPropertiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	req.RentType = "short"
	listReq := models.ListPropertiesRequest{
		DistrictID:  req.DistrictID,
		MinPrice:    req.MinPrice,
		MaxPrice:    req.MaxPrice,
		Bedrooms:    req.Bedrooms,
		ListingType: "rent",
		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}
	// TODO: 需要在 Property 表添加 rent_type 字段来区分短租/长租

	properties, total, err := h.propertyService.ListProperties(c.Request.Context(), &listReq)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.SuccessWithPagination(c, properties, models.NewPagination(req.Page, req.PageSize, total))
}

// 14. ListLongTermRent 长租房源
// ListLongTermRent godoc
// @Summary      长租房源
// @Description  获取长租（>=6个月）房产列表
// @Tags         Properties
// @Accept       json
// @Produce      json
// @Param        district_id    query     int     false  "地区ID"
// @Param        min_price      query     number  false  "最低月租"
// @Param        max_price      query     number  false  "最高月租"
// @Param        bedrooms       query     int     false  "卧室数"
// @Param        page           query     int     false  "页码"       default(1)
// @Param        page_size      query     int     false  "每页数量"   default(20)
// @Success      200            {object}  models.PaginatedResponse{data=[]models.PropertyListItemResponse}
// @Router       /api/v1/properties/rent/long-term [get]
func (h *PropertyHandler) ListLongTermRent(c *gin.Context) {
	var req models.ListRentPropertiesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		models.BadRequest(c, err.Error())
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}

	req.RentType = "long"
	listReq := models.ListPropertiesRequest{
		DistrictID:  req.DistrictID,
		MinPrice:    req.MinPrice,
		MaxPrice:    req.MaxPrice,
		Bedrooms:    req.Bedrooms,
		ListingType: "rent",
		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}
	// TODO: 需要在 Property 表添加 rent_type 字段来区分短租/长租

	properties, total, err := h.propertyService.ListProperties(c.Request.Context(), &listReq)
	if err != nil {
		models.InternalError(c, err.Error())
		return
	}

	models.SuccessWithPagination(c, properties, models.NewPagination(req.Page, req.PageSize, total))
}
