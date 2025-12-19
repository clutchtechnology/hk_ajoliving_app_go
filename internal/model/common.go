package model

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
