package databases

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

type SchoolRepo struct {
	db *gorm.DB
}

func NewSchoolRepo(db *gorm.DB) *SchoolRepo {
	return &SchoolRepo{db: db}
}

// FindAll 查询所有学校（分页+筛选）
func (r *SchoolRepo) FindAll(ctx context.Context, req *models.ListSchoolsRequest) ([]*models.School, int64, error) {
	var schools []*models.School
	var total int64

	query := r.db.WithContext(ctx).Model(&models.School{})

	// 应用筛选条件
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}
	if req.Category != nil {
		query = query.Where("category = ?", *req.Category)
	}
	if req.Gender != nil {
		query = query.Where("gender = ?", *req.Gender)
	}
	if req.SchoolNetID != nil {
		query = query.Where("school_net_id = ?", *req.SchoolNetID)
	}
	if req.DistrictID != nil {
		query = query.Where("district_id = ?", *req.DistrictID)
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("name_zh_hant LIKE ? OR name_zh_hans LIKE ? OR name_en LIKE ?",
			keyword, keyword, keyword)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err := query.
		Preload("SchoolNet").
		Preload("District").
		Offset(offset).
		Limit(req.PageSize).
		Order("name_zh_hant ASC").
		Find(&schools).Error

	return schools, total, err
}

// FindByID 根据ID查询学校
func (r *SchoolRepo) FindByID(ctx context.Context, id uint) (*models.School, error) {
	var school models.School
	err := r.db.WithContext(ctx).
		Preload("SchoolNet").
		Preload("District").
		First(&school, id).Error
	return &school, err
}

// Search 搜索学校
func (r *SchoolRepo) Search(ctx context.Context, keyword string, page, pageSize int) ([]*models.School, int64, error) {
	var schools []*models.School
	var total int64

	searchPattern := "%" + keyword + "%"
	query := r.db.WithContext(ctx).Model(&models.School{}).
		Where("name_zh_hant LIKE ? OR name_zh_hans LIKE ? OR name_en LIKE ?",
			searchPattern, searchPattern, searchPattern)

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.
		Preload("SchoolNet").
		Preload("District").
		Offset(offset).
		Limit(pageSize).
		Order("name_zh_hant ASC").
		Find(&schools).Error

	return schools, total, err
}

// IncrementViewCount 增加浏览次数
func (r *SchoolRepo) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&models.School{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).
		Error
}
