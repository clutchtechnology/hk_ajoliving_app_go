package model

import (
	"time"
)

// SearchType 搜索类型
type SearchType string

const (
	SearchTypeGlobal    SearchType = "global"    // 全局搜索
	SearchTypeProperty  SearchType = "property"  // 房产搜索
	SearchTypeEstate    SearchType = "estate"    // 屋苑搜索
	SearchTypeAgent     SearchType = "agent"     // 代理人搜索
	SearchTypeAgency    SearchType = "agency"    // 代理公司搜索
	SearchTypeNews      SearchType = "news"      // 新闻搜索
	SearchTypeSchool    SearchType = "school"    // 学校搜索
	SearchTypeFurniture SearchType = "furniture" // 家具搜索
)

// SearchHistory 搜索历史记录模型
type SearchHistory struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     *uint      `gorm:"index" json:"user_id,omitempty"`                    // 用户ID（可选，未登录用户为空）
	Keyword    string     `gorm:"type:varchar(200);not null;index" json:"keyword"`   // 搜索关键词
	SearchType SearchType `gorm:"type:varchar(20);not null;index" json:"search_type"` // 搜索类型
	ResultCount int       `gorm:"default:0" json:"result_count"`                      // 结果数量
	IPAddress  *string    `gorm:"type:varchar(50)" json:"ip_address,omitempty"`      // IP地址
	UserAgent  *string    `gorm:"type:varchar(500)" json:"user_agent,omitempty"`     // 用户代理
	CreatedAt  time.Time  `gorm:"autoCreateTime;index" json:"created_at"`

	// 关联
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
}

// TableName 指定表名
func (SearchHistory) TableName() string {
	return "search_histories"
}

// IsGlobalSearch 判断是否为全局搜索
func (sh *SearchHistory) IsGlobalSearch() bool {
	return sh.SearchType == SearchTypeGlobal
}

// HasResults 判断是否有搜索结果
func (sh *SearchHistory) HasResults() bool {
	return sh.ResultCount > 0
}

// IsAuthenticated 判断是否为已登录用户的搜索
func (sh *SearchHistory) IsAuthenticated() bool {
	return sh.UserID != nil
}
