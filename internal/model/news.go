package model

import (
	"time"

	"gorm.io/gorm"
)

// NewsCategory 新闻分类
type NewsCategory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	NameZhHant  string         `gorm:"size:100;not null" json:"name_zh_hant"`           // 繁体中文名
	NameZhHans  string         `gorm:"size:100;not null" json:"name_zh_hans"`           // 简体中文名
	NameEn      string         `gorm:"size:100;not null" json:"name_en"`                // 英文名
	Slug        string         `gorm:"size:100;uniqueIndex;not null" json:"slug"`       // URL 别名
	Description string         `gorm:"type:text" json:"description"`                    // 分类描述
	SortOrder   int            `gorm:"default:0" json:"sort_order"`                     // 排序
	IsActive    bool           `gorm:"default:true" json:"is_active"`                   // 是否启用
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (NewsCategory) TableName() string {
	return "news_categories"
}

// News 新闻
type News struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CategoryID      uint           `gorm:"index;not null" json:"category_id"`                 // 分类ID
	Title           string         `gorm:"size:500;not null" json:"title"`                    // 标题
	Subtitle        string         `gorm:"size:500" json:"subtitle"`                          // 副标题
	Summary         string         `gorm:"type:text" json:"summary"`                          // 摘要
	Content         string         `gorm:"type:text;not null" json:"content"`                 // 内容（HTML）
	CoverImageURL   string         `gorm:"size:500" json:"cover_image_url"`                   // 封面图
	SourceName      string         `gorm:"size:200" json:"source_name"`                       // 来源名称
	SourceURL       string         `gorm:"size:1000" json:"source_url"`                       // 来源链接
	AuthorName      string         `gorm:"size:100" json:"author_name"`                       // 作者
	PublishedAt     *time.Time     `gorm:"index" json:"published_at"`                         // 发布时间
	ViewCount       int64          `gorm:"default:0" json:"view_count"`                       // 浏览量
	LikeCount       int64          `gorm:"default:0" json:"like_count"`                       // 点赞数
	CommentCount    int64          `gorm:"default:0" json:"comment_count"`                    // 评论数
	IsFeatured      bool           `gorm:"default:false;index" json:"is_featured"`            // 是否精选
	IsHot           bool           `gorm:"default:false;index" json:"is_hot"`                 // 是否热门
	IsTop           bool           `gorm:"default:false;index" json:"is_top"`                 // 是否置顶
	Status          string         `gorm:"size:20;default:'published';index" json:"status"`   // 状态: draft, published, archived
	Tags            string         `gorm:"size:500" json:"tags"`                              // 标签（逗号分隔）
	Keywords        string         `gorm:"size:500" json:"keywords"`                          // 关键词
	MetaDescription string         `gorm:"size:500" json:"meta_description"`                  // SEO 描述
	CrawlerSource   string         `gorm:"size:100" json:"crawler_source"`                    // 爬虫来源
	CrawledAt       *time.Time     `json:"crawled_at"`                                        // 爬取时间
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	Category        *NewsCategory  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (News) TableName() string {
	return "news"
}

// BeforeCreate GORM 钩子：创建前设置默认值
func (n *News) BeforeCreate(tx *gorm.DB) error {
	if n.Status == "" {
		n.Status = "published"
	}
	if n.PublishedAt == nil {
		now := time.Now()
		n.PublishedAt = &now
	}
	return nil
}

// 状态常量
const (
	NewsStatusDraft     = "draft"     // 草稿
	NewsStatusPublished = "published" // 已发布
	NewsStatusArchived  = "archived"  // 已归档
)

// 辅助方法

// IsPublished 是否已发布
func (n *News) IsPublished() bool {
	return n.Status == NewsStatusPublished
}

// IncrementViewCount 增加浏览量
func (n *News) IncrementViewCount() {
	n.ViewCount++
}

// IncrementLikeCount 增加点赞数
func (n *News) IncrementLikeCount() {
	n.LikeCount++
}

// GetTagList 获取标签列表
func (n *News) GetTagList() []string {
	if n.Tags == "" {
		return []string{}
	}
	tags := []string{}
	for _, tag := range splitByComma(n.Tags) {
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}

// splitByComma 按逗号分割字符串
func splitByComma(s string) []string {
	result := []string{}
	current := ""
	for _, c := range s {
		if c == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
