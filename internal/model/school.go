package model

import (
	"time"

	"gorm.io/gorm"
)

// SchoolNet 校网
type SchoolNet struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	NetCode     string         `gorm:"size:50;uniqueIndex;not null" json:"net_code"`     // 校网编号
	NameZhHant  string         `gorm:"size:200;not null" json:"name_zh_hant"`            // 繁体中文名
	NameZhHans  string         `gorm:"size:200;not null" json:"name_zh_hans"`            // 简体中文名
	NameEn      string         `gorm:"size:200;not null" json:"name_en"`                 // 英文名
	DistrictID  uint           `gorm:"index;not null" json:"district_id"`                // 所属地区ID
	Description string         `gorm:"type:text" json:"description"`                     // 描述
	Level       string         `gorm:"size:20;index" json:"level"`                       // 学校级别: primary, secondary
	SchoolCount int            `gorm:"default:0" json:"school_count"`                    // 学校数量
	MapData     string         `gorm:"type:text" json:"map_data"`                        // 地图数据（JSON格式）
	IsActive    bool           `gorm:"default:true;index" json:"is_active"`              // 是否启用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	District    *District      `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
	Schools     []School       `gorm:"foreignKey:SchoolNetID" json:"schools,omitempty"`
}

func (SchoolNet) TableName() string {
	return "school_nets"
}

// School 学校
type School struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SchoolNetID  uint           `gorm:"index;not null" json:"school_net_id"`              // 校网ID
	DistrictID   uint           `gorm:"index;not null" json:"district_id"`                // 所属地区ID
	NameZhHant   string         `gorm:"size:200;not null;index" json:"name_zh_hant"`      // 繁体中文名
	NameZhHans   string         `gorm:"size:200;not null" json:"name_zh_hans"`            // 简体中文名
	NameEn       string         `gorm:"size:200;not null" json:"name_en"`                 // 英文名
	SchoolCode   string         `gorm:"size:50;uniqueIndex" json:"school_code"`           // 学校编号
	Category     string         `gorm:"size:50;index" json:"category"`                    // 类别: government, aided, direct_subsidy, private, international
	Level        string         `gorm:"size:20;index" json:"level"`                       // 级别: kindergarten, primary, secondary
	Gender       string         `gorm:"size:20" json:"gender"`                            // 性别: co-ed, boys, girls
	Religion     string         `gorm:"size:50" json:"religion"`                          // 宗教
	Address      string         `gorm:"size:500" json:"address"`                          // 地址
	Phone        string         `gorm:"size:50" json:"phone"`                             // 电话
	Email        string         `gorm:"size:200" json:"email"`                            // 邮箱
	Website      string         `gorm:"size:500" json:"website"`                          // 网站
	Principal    string         `gorm:"size:100" json:"principal"`                        // 校长
	FoundedYear  int            `json:"founded_year"`                                     // 成立年份
	StudentCount int            `gorm:"default:0" json:"student_count"`                   // 学生人数
	TeacherCount int            `gorm:"default:0" json:"teacher_count"`                   // 教师人数
	Rating       float64        `gorm:"type:decimal(3,2);default:0" json:"rating"`        // 评分
	Features     string         `gorm:"type:text" json:"features"`                        // 特色（JSON数组）
	Facilities   string         `gorm:"type:text" json:"facilities"`                      // 设施（JSON数组）
	Latitude     float64        `gorm:"type:decimal(10,8)" json:"latitude"`               // 纬度
	Longitude    float64        `gorm:"type:decimal(11,8)" json:"longitude"`              // 经度
	LogoURL      string         `gorm:"size:500" json:"logo_url"`                         // Logo图片
	CoverImageURL string        `gorm:"size:500" json:"cover_image_url"`                  // 封面图片
	IsActive     bool           `gorm:"default:true;index" json:"is_active"`              // 是否启用
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	SchoolNet   *SchoolNet     `gorm:"foreignKey:SchoolNetID" json:"school_net,omitempty"`
	District    *District      `gorm:"foreignKey:DistrictID" json:"district,omitempty"`
}

func (School) TableName() string {
	return "schools"
}

// 学校级别常量
const (
	SchoolLevelKindergarten = "kindergarten" // 幼儿园
	SchoolLevelPrimary      = "primary"      // 小学
	SchoolLevelSecondary    = "secondary"    // 中学
)

// 学校类别常量
const (
	SchoolCategoryGovernment     = "government"      // 官立
	SchoolCategoryAided          = "aided"           // 资助
	SchoolCategoryDirectSubsidy  = "direct_subsidy"  // 直资
	SchoolCategoryPrivate        = "private"         // 私立
	SchoolCategoryInternational  = "international"   // 国际
)

// 性别常量
const (
	SchoolGenderCoEd = "co-ed" // 男女校
	SchoolGenderBoys = "boys"  // 男校
	SchoolGenderGirls = "girls" // 女校
)

// 辅助方法

// IsPrimary 是否小学
func (s *School) IsPrimary() bool {
	return s.Level == SchoolLevelPrimary
}

// IsSecondary 是否中学
func (s *School) IsSecondary() bool {
	return s.Level == SchoolLevelSecondary
}

// IsGovernmentSchool 是否官立学校
func (s *School) IsGovernmentSchool() bool {
	return s.Category == SchoolCategoryGovernment
}

// HasLocation 是否有位置信息
func (s *School) HasLocation() bool {
	return s.Latitude != 0 && s.Longitude != 0
}

// GetCategoryName 获取类别名称
func (s *School) GetCategoryName() string {
	categoryNames := map[string]string{
		SchoolCategoryGovernment:    "官立",
		SchoolCategoryAided:         "资助",
		SchoolCategoryDirectSubsidy: "直资",
		SchoolCategoryPrivate:       "私立",
		SchoolCategoryInternational: "国际",
	}
	if name, ok := categoryNames[s.Category]; ok {
		return name
	}
	return s.Category
}

// GetLevelName 获取级别名称
func (s *School) GetLevelName() string {
	levelNames := map[string]string{
		SchoolLevelKindergarten: "幼儿园",
		SchoolLevelPrimary:      "小学",
		SchoolLevelSecondary:    "中学",
	}
	if name, ok := levelNames[s.Level]; ok {
		return name
	}
	return s.Level
}
