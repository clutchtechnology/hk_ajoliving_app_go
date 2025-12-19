package response

import "time"

// NewsCategoryResponse 新闻分类响应
type NewsCategoryResponse struct {
	ID          uint      `json:"id"`
	NameZhHant  string    `json:"name_zh_hant"`
	NameZhHans  string    `json:"name_zh_hans"`
	NameEn      string    `json:"name_en"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order"`
	NewsCount   int64     `json:"news_count,omitempty"` // 该分类下的新闻数量
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewsListItemResponse 新闻列表项响应
type NewsListItemResponse struct {
	ID            uint                  `json:"id"`
	CategoryID    uint                  `json:"category_id"`
	Category      *NewsCategoryResponse `json:"category,omitempty"`
	Title         string                `json:"title"`
	Subtitle      string                `json:"subtitle,omitempty"`
	Summary       string                `json:"summary"`
	CoverImageURL string                `json:"cover_image_url,omitempty"`
	SourceName    string                `json:"source_name,omitempty"`
	AuthorName    string                `json:"author_name,omitempty"`
	PublishedAt   *time.Time            `json:"published_at"`
	ViewCount     int64                 `json:"view_count"`
	LikeCount     int64                 `json:"like_count"`
	CommentCount  int64                 `json:"comment_count"`
	IsFeatured    bool                  `json:"is_featured"`
	IsHot         bool                  `json:"is_hot"`
	IsTop         bool                  `json:"is_top"`
	Tags          []string              `json:"tags,omitempty"`
	CreatedAt     time.Time             `json:"created_at"`
}

// NewsResponse 新闻详情响应
type NewsResponse struct {
	ID              uint                  `json:"id"`
	CategoryID      uint                  `json:"category_id"`
	Category        *NewsCategoryResponse `json:"category,omitempty"`
	Title           string                `json:"title"`
	Subtitle        string                `json:"subtitle,omitempty"`
	Summary         string                `json:"summary"`
	Content         string                `json:"content"`
	CoverImageURL   string                `json:"cover_image_url,omitempty"`
	SourceName      string                `json:"source_name,omitempty"`
	SourceURL       string                `json:"source_url,omitempty"`
	AuthorName      string                `json:"author_name,omitempty"`
	PublishedAt     *time.Time            `json:"published_at"`
	ViewCount       int64                 `json:"view_count"`
	LikeCount       int64                 `json:"like_count"`
	CommentCount    int64                 `json:"comment_count"`
	IsFeatured      bool                  `json:"is_featured"`
	IsHot           bool                  `json:"is_hot"`
	IsTop           bool                  `json:"is_top"`
	Status          string                `json:"status"`
	Tags            []string              `json:"tags,omitempty"`
	Keywords        string                `json:"keywords,omitempty"`
	MetaDescription string                `json:"meta_description,omitempty"`
	CrawlerSource   string                `json:"crawler_source,omitempty"`
	CrawledAt       *time.Time            `json:"crawled_at,omitempty"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

// RelatedNewsResponse 相关新闻响应
type RelatedNewsResponse struct {
	ID            uint       `json:"id"`
	Title         string     `json:"title"`
	Summary       string     `json:"summary"`
	CoverImageURL string     `json:"cover_image_url,omitempty"`
	PublishedAt   *time.Time `json:"published_at"`
	ViewCount     int64      `json:"view_count"`
}
