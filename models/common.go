package models

// 此文件包含所有 model 共享的常量、类型和辅助函数

// Status 通用状态类型
type Status string

const (
	StatusActive    Status = "active"    // 活跃
	StatusInactive  Status = "inactive"  // 停用
	StatusSuspended Status = "suspended" // 暂停
	StatusPending   Status = "pending"   // 待处理
	StatusApproved  Status = "approved"  // 已批准
	StatusRejected  Status = "rejected"  // 已拒绝
)

// ImageType 图片类型
type ImageType string

const (
	ImageTypeCover     ImageType = "cover"     // 封面图
	ImageTypeInterior  ImageType = "interior"  // 室内
	ImageTypeExterior  ImageType = "exterior"  // 外观
	ImageTypeFloorplan ImageType = "floorplan" // 户型图
	ImageTypeLocation  ImageType = "location"  // 位置图
	ImageTypeFacility  ImageType = "facility"  // 设施
	ImageTypeAerial    ImageType = "aerial"    // 航拍
)

// PublisherType 发布者类型
type PublisherType string

const (
	PublisherTypeIndividual PublisherType = "individual" // 个人
	PublisherTypeAgency     PublisherType = "agency"     // 代理公司
)

// Language 语言常量
const (
	LangZhHant = "zh-Hant" // 繁体中文
	LangZhHans = "zh-Hans" // 简体中文
	LangEn     = "en"      // 英文
)

// ============ 统计相关 DTO ============

// GetPropertyStatisticsRequest 获取房产统计请求
type GetPropertyStatisticsRequest struct {
	DistrictID   *uint   `form:"district_id"`
	EstateID     *uint   `form:"estate_id"`
	EstateNo     *string `form:"estate_no"`
	ListingType  *string `form:"listing_type" binding:"omitempty,oneof=sale rent"`
	PropertyType *string `form:"property_type"`
	Status       *string `form:"status"`
	Period       *string `form:"period" binding:"omitempty,oneof=day week month year"`
	StartDate    *string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate      *string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
}

// GetTransactionStatisticsRequest 获取交易统计请求
type GetTransactionStatisticsRequest struct {
	DistrictID  *uint   `form:"district_id"`
	EstateID    *uint   `form:"estate_id"`
	EstateNo    *string `form:"estate_no"`
	ListingType *string `form:"listing_type" binding:"omitempty,oneof=sale rent"`
	Period      *string `form:"period" binding:"omitempty,oneof=day week month year"`
	StartDate   *string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate     *string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Limit       int     `form:"limit,default=10" binding:"min=1,max=100"`
}

// GetUserStatisticsRequest 获取用户统计请求
type GetUserStatisticsRequest struct {
	Role      *string `form:"role"`
	Status    *string `form:"status"`
	Period    *string `form:"period" binding:"omitempty,oneof=day week month year"`
	StartDate *string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   *string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
}

// PaginationRequest 通用分页请求
type PaginationRequest struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`
	SortBy    string `form:"sort_by,default=created_at"`
	SortOrder string `form:"sort_order,default=desc" binding:"omitempty,oneof=asc desc"`
}
