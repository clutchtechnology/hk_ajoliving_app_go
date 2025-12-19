package repository

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
)

// NewsRepository 新闻数据访问接口
type NewsRepository interface {
	// 新闻相关
	List(ctx context.Context, filter *request.ListNewsRequest) ([]*model.News, int64, error)
	GetByID(ctx context.Context, id uint) (*model.News, error)
	GetHotNews(ctx context.Context, limit int) ([]*model.News, error)
	GetFeaturedNews(ctx context.Context, limit int) ([]*model.News, error)
	GetLatestNews(ctx context.Context, limit int) ([]*model.News, error)
	GetRelatedNews(ctx context.Context, newsID uint, categoryID uint, limit int) ([]*model.News, error)
	IncrementViewCount(ctx context.Context, id uint) error
	
	// 分类相关
	GetAllCategories(ctx context.Context) ([]*model.NewsCategory, error)
	GetCategoryByID(ctx context.Context, id uint) (*model.NewsCategory, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*model.NewsCategory, error)
}

type newsRepository struct {
	db *gorm.DB
}

// NewNewsRepository 创建新闻仓库
func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{db: db}
}

// 新闻相关

func (r *newsRepository) List(ctx context.Context, filter *request.ListNewsRequest) ([]*model.News, int64, error) {
	var news []*model.News
	var total int64
	
	query := r.db.WithContext(ctx).Model(&model.News{})
	
	// 默认只查询已发布的新闻
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	} else {
		query = query.Where("status = ?", model.NewsStatusPublished)
	}
	
	// 应用筛选条件
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filter.IsFeatured)
	}
	if filter.IsHot != nil {
		query = query.Where("is_hot = ?", *filter.IsHot)
	}
	if filter.IsTop != nil {
		query = query.Where("is_top = ?", *filter.IsTop)
	}
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("title LIKE ? OR summary LIKE ? OR content LIKE ?", keyword, keyword, keyword)
	}
	if filter.Tag != "" {
		tagPattern := "%" + filter.Tag + "%"
		query = query.Where("tags LIKE ?", tagPattern)
	}
	
	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// 设置默认值
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	sortBy := filter.SortBy
	if sortBy == "" {
		sortBy = "published_at"
	}
	sortOrder := filter.SortOrder
	if sortOrder == "" {
		sortOrder = "desc"
	}
	
	// 置顶新闻优先
	query = query.Order("is_top DESC")
	
	// 分页和排序
	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)
	query = query.Order(sortBy + " " + sortOrder)
	
	// 预加载关联
	query = query.Preload("Category")
	
	if err := query.Find(&news).Error; err != nil {
		return nil, 0, err
	}
	
	return news, total, nil
}

func (r *newsRepository) GetByID(ctx context.Context, id uint) (*model.News, error) {
	var news model.News
	
	err := r.db.WithContext(ctx).
		Preload("Category").
		First(&news, id).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &news, nil
}

func (r *newsRepository) GetHotNews(ctx context.Context, limit int) ([]*model.News, error) {
	var news []*model.News
	
	if limit <= 0 {
		limit = 10
	}
	
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("status = ? AND is_hot = ?", model.NewsStatusPublished, true).
		Order("view_count DESC, published_at DESC").
		Limit(limit).
		Find(&news).Error
		
	if err != nil {
		return nil, err
	}
	
	return news, nil
}

func (r *newsRepository) GetFeaturedNews(ctx context.Context, limit int) ([]*model.News, error) {
	var news []*model.News
	
	if limit <= 0 {
		limit = 10
	}
	
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("status = ? AND is_featured = ?", model.NewsStatusPublished, true).
		Order("published_at DESC").
		Limit(limit).
		Find(&news).Error
		
	if err != nil {
		return nil, err
	}
	
	return news, nil
}

func (r *newsRepository) GetLatestNews(ctx context.Context, limit int) ([]*model.News, error) {
	var news []*model.News
	
	if limit <= 0 {
		limit = 10
	}
	
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("status = ?", model.NewsStatusPublished).
		Order("published_at DESC").
		Limit(limit).
		Find(&news).Error
		
	if err != nil {
		return nil, err
	}
	
	return news, nil
}

func (r *newsRepository) GetRelatedNews(ctx context.Context, newsID uint, categoryID uint, limit int) ([]*model.News, error) {
	var news []*model.News
	
	if limit <= 0 {
		limit = 5
	}
	
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("status = ? AND id != ? AND category_id = ?", model.NewsStatusPublished, newsID, categoryID).
		Order("published_at DESC").
		Limit(limit).
		Find(&news).Error
		
	if err != nil {
		return nil, err
	}
	
	return news, nil
}

func (r *newsRepository) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&model.News{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).
		Error
}

// 分类相关

func (r *newsRepository) GetAllCategories(ctx context.Context) ([]*model.NewsCategory, error) {
	var categories []*model.NewsCategory
	
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("sort_order ASC, name_zh_hant ASC").
		Find(&categories).Error
		
	if err != nil {
		return nil, err
	}
	
	return categories, nil
}

func (r *newsRepository) GetCategoryByID(ctx context.Context, id uint) (*model.NewsCategory, error) {
	var category model.NewsCategory
	
	err := r.db.WithContext(ctx).First(&category, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &category, nil
}

func (r *newsRepository) GetCategoryBySlug(ctx context.Context, slug string) (*model.NewsCategory, error) {
	var category model.NewsCategory
	
	slug = strings.TrimSpace(strings.ToLower(slug))
	
	err := r.db.WithContext(ctx).
		Where("LOWER(slug) = ?", slug).
		First(&category).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &category, nil
}
