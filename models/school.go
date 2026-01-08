package models

import (
	"time"

	"gorm.io/gorm"
)

// ============ GORM Model ============

// SchoolNet 校网模型
type SchoolNet struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Code        string         `gorm:"size:50;uniqueIndex;not null" json:"code"`        // 校网编号（如：11, 12, 34）
	NameZhHant  string         `gorm:"size:200;not null;index" json:"name_zh_hant"`     // 中文繁体名称
	NameZhHans  string         `gorm:"size:200" json:"name_zh_hans,omitempty"`          // 中文简体名称
	NameEn      string         `gorm:"size:200" json:"name_en,omitempty"`               // 英文名称
	Type        string         `gorm:"size:20;not null;index" json:"type"`              // primary=小学, secondary=中学
	DistrictID  uint           `gorm:"not null;index" json:"district_id"`               // 所属地区ID
	Description string         `gorm:"type:text" json:"description,omitempty"`          // 校网描述
	Coverage    string         `gorm:"type:text" json:"coverage,omitempty"`             // 覆盖范围说明
	SchoolCount int            `gorm:"default:0" json:"school_count"`                   // 校网内学校数量
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	District *District `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Schools  []School  `gorm:"foreignKey:SchoolNetID" json:"schools,omitempty"`
}

func (SchoolNet) TableName() string {
	return "school_nets"
}

// School 学校模型
type School struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	NameZhHant    string         `gorm:"size:200;not null;index" json:"name_zh_hant"`      // 中文繁体名称
	NameZhHans    string         `gorm:"size:200" json:"name_zh_hans,omitempty"`           // 中文简体名称
	NameEn        string         `gorm:"size:200;index" json:"name_en,omitempty"`          // 英文名称
	Type          string         `gorm:"size:20;not null;index" json:"type"`               // primary=小学, secondary=中学
	Category      string         `gorm:"size:50;not null;index" json:"category"`           // government=官立, aided=资助, direct_subsidy=直资, private=私立, international=国际
	Gender        string         `gorm:"size:20;not null" json:"gender"`                   // coed=男女校, boys=男校, girls=女校
	SchoolNetID   uint           `gorm:"index" json:"school_net_id"`                       // 所属校网ID（可选，某些学校不在校网内）
	DistrictID    uint           `gorm:"not null;index" json:"district_id"`                // 所属地区ID
	Address       string         `gorm:"size:500;not null" json:"address"`                 // 学校地址
	Phone         string         `gorm:"size:50" json:"phone,omitempty"`                   // 联系电话
	Email         string         `gorm:"size:255" json:"email,omitempty"`                  // 电子邮件
	Website       string         `gorm:"size:500" json:"website,omitempty"`                // 学校网站
	EstablishedAt *time.Time     `json:"established_at,omitempty"`                         // 创校年份
	Principal     string         `gorm:"size:100" json:"principal,omitempty"`              // 校长姓名
	Religion      string         `gorm:"size:50" json:"religion,omitempty"`                // 宗教背景（如：基督教、天主教、佛教等）
	Curriculum    string         `gorm:"size:100" json:"curriculum,omitempty"`             // 课程类型（如：本地、IB、英式等）
	StudentCount  int            `gorm:"default:0" json:"student_count"`                   // 学生人数
	TeacherCount  int            `gorm:"default:0" json:"teacher_count"`                   // 教师人数
	Rating        float64        `gorm:"type:decimal(3,2)" json:"rating"`                  // 评分（0-5）
	Description   string         `gorm:"type:text" json:"description,omitempty"`           // 学校简介
	ViewCount     int            `gorm:"default:0" json:"view_count"`                      // 浏览次数
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	SchoolNet *SchoolNet `gorm:"foreignKey:SchoolNetID" json:"school_net,omitempty"`
	District  *District  `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
}

func (School) TableName() string {
	return "schools"
}

// ============ Request DTO ============

// ListSchoolNetsRequest 校网列表请求
type ListSchoolNetsRequest struct {
	Type       *string `form:"type" binding:"omitempty,oneof=primary secondary"`
	DistrictID *uint   `form:"district_id"`
	Keyword    string  `form:"keyword"`
	Page       int     `form:"page,default=1" binding:"min=1"`
	PageSize   int     `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListSchoolsRequest 学校列表请求
type ListSchoolsRequest struct {
	Type        *string `form:"type" binding:"omitempty,oneof=primary secondary"`
	Category    *string `form:"category" binding:"omitempty,oneof=government aided direct_subsidy private international"`
	Gender      *string `form:"gender" binding:"omitempty,oneof=coed boys girls"`
	SchoolNetID *uint   `form:"school_net_id"`
	DistrictID  *uint   `form:"district_id"`
	Keyword     string  `form:"keyword"`
	Page        int     `form:"page,default=1" binding:"min=1"`
	PageSize    int     `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ============ Response DTO ============

// SchoolNetResponse 校网响应
type SchoolNetResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	NameZhHant  string `json:"name_zh_hant"`
	NameZhHans  string `json:"name_zh_hans,omitempty"`
	NameEn      string `json:"name_en,omitempty"`
	Type        string `json:"type"`
	DistrictID  uint   `json:"district_id"`
	Description string `json:"description,omitempty"`
	Coverage    string `json:"coverage,omitempty"`
	SchoolCount int    `json:"school_count"`
}

// SchoolNetDetailResponse 校网详情响应
type SchoolNetDetailResponse struct {
	ID          uint              `json:"id"`
	Code        string            `json:"code"`
	NameZhHant  string            `json:"name_zh_hant"`
	NameZhHans  string            `json:"name_zh_hans,omitempty"`
	NameEn      string            `json:"name_en,omitempty"`
	Type        string            `json:"type"`
	DistrictID  uint              `json:"district_id"`
	Description string            `json:"description,omitempty"`
	Coverage    string            `json:"coverage,omitempty"`
	SchoolCount int               `json:"school_count"`
	District    *DistrictResponse `json:"district,omitempty"`
}

// SchoolResponse 学校响应
type SchoolResponse struct {
	ID           uint    `json:"id"`
	NameZhHant   string  `json:"name_zh_hant"`
	NameZhHans   string  `json:"name_zh_hans,omitempty"`
	NameEn       string  `json:"name_en,omitempty"`
	Type         string  `json:"type"`
	Category     string  `json:"category"`
	Gender       string  `json:"gender"`
	SchoolNetID  uint    `json:"school_net_id,omitempty"`
	DistrictID   uint    `json:"district_id"`
	Address      string  `json:"address"`
	Rating       float64 `json:"rating"`
	StudentCount int     `json:"student_count"`
}

// SchoolDetailResponse 学校详情响应
type SchoolDetailResponse struct {
	ID            uint                 `json:"id"`
	NameZhHant    string               `json:"name_zh_hant"`
	NameZhHans    string               `json:"name_zh_hans,omitempty"`
	NameEn        string               `json:"name_en,omitempty"`
	Type          string               `json:"type"`
	Category      string               `json:"category"`
	Gender        string               `json:"gender"`
	SchoolNetID   uint                 `json:"school_net_id,omitempty"`
	DistrictID    uint                 `json:"district_id"`
	Address       string               `json:"address"`
	Phone         string               `json:"phone,omitempty"`
	Email         string               `json:"email,omitempty"`
	Website       string               `json:"website,omitempty"`
	EstablishedAt *time.Time           `json:"established_at,omitempty"`
	Principal     string               `json:"principal,omitempty"`
	Religion      string               `json:"religion,omitempty"`
	Curriculum    string               `json:"curriculum,omitempty"`
	StudentCount  int                  `json:"student_count"`
	TeacherCount  int                  `json:"teacher_count"`
	Rating        float64              `json:"rating"`
	Description   string               `json:"description,omitempty"`
	ViewCount     int                  `json:"view_count"`
	SchoolNet     *SchoolNetResponse   `json:"school_net,omitempty"`
	District      *DistrictResponse    `json:"district,omitempty"`
}

// PaginatedSchoolNetsResponse 分页校网响应
type PaginatedSchoolNetsResponse struct {
	Items      []*SchoolNetResponse `json:"items"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// PaginatedSchoolsResponse 分页学校响应
type PaginatedSchoolsResponse struct {
	Items      []*SchoolResponse `json:"items"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}
