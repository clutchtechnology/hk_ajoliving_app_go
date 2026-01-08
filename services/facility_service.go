package services

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
)

// FacilityService Methods:
// 0. NewFacilityService(repo *databases.FacilityRepo) -> 注入依赖
// 1. ListFacilities(ctx context.Context, category string) -> 获取设施列表
// 2. GetFacility(ctx context.Context, id uint) -> 获取设施详情
// 3. CreateFacility(ctx context.Context, req *models.CreateFacilityRequest) -> 创建设施
// 4. UpdateFacility(ctx context.Context, id uint, req *models.UpdateFacilityRequest) -> 更新设施
// 5. DeleteFacility(ctx context.Context, id uint) -> 删除设施

type FacilityService struct {
	repo *databases.FacilityRepo
}

// 0. NewFacilityService 构造函数
func NewFacilityService(repo *databases.FacilityRepo) *FacilityService {
	return &FacilityService{repo: repo}
}

// 1. ListFacilities 获取设施列表
func (s *FacilityService) ListFacilities(ctx context.Context, category string) ([]models.FacilityResponse, error) {
	facilities, err := s.repo.FindAll(ctx, category)
	if err != nil {
		return nil, err
	}

	var responses []models.FacilityResponse
	for _, facility := range facilities {
		responses = append(responses, models.FacilityResponse{
			ID:         facility.ID,
			NameZhHant: facility.NameZhHant,
			NameZhHans: facility.NameZhHans,
			NameEn:     facility.NameEn,
			Icon:       facility.Icon,
			Category:   facility.Category,
			SortOrder:  facility.SortOrder,
			CreatedAt:  facility.CreatedAt,
			UpdatedAt:  facility.UpdatedAt,
		})
	}

	return responses, nil
}

// 2. GetFacility 获取设施详情
func (s *FacilityService) GetFacility(ctx context.Context, id uint) (*models.FacilityResponse, error) {
	facility, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	return &models.FacilityResponse{
		ID:         facility.ID,
		NameZhHant: facility.NameZhHant,
		NameZhHans: facility.NameZhHans,
		NameEn:     facility.NameEn,
		Icon:       facility.Icon,
		Category:   facility.Category,
		SortOrder:  facility.SortOrder,
		CreatedAt:  facility.CreatedAt,
		UpdatedAt:  facility.UpdatedAt,
	}, nil
}

// 3. CreateFacility 创建设施
func (s *FacilityService) CreateFacility(ctx context.Context, req *models.CreateFacilityRequest) (*models.FacilityResponse, error) {
	// 检查名称是否已存在
	exists, err := s.repo.CheckNameExists(ctx, req.NameZhHant, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("facility name already exists")
	}

	facility := &models.Facility{
		NameZhHant: req.NameZhHant,
		NameZhHans: req.NameZhHans,
		NameEn:     req.NameEn,
		Icon:       req.Icon,
		Category:   req.Category,
		SortOrder:  req.SortOrder,
	}

	if err := s.repo.Create(ctx, facility); err != nil {
		return nil, err
	}

	return &models.FacilityResponse{
		ID:         facility.ID,
		NameZhHant: facility.NameZhHant,
		NameZhHans: facility.NameZhHans,
		NameEn:     facility.NameEn,
		Icon:       facility.Icon,
		Category:   facility.Category,
		SortOrder:  facility.SortOrder,
		CreatedAt:  facility.CreatedAt,
		UpdatedAt:  facility.UpdatedAt,
	}, nil
}

// 4. UpdateFacility 更新设施
func (s *FacilityService) UpdateFacility(ctx context.Context, id uint, req *models.UpdateFacilityRequest) (*models.FacilityResponse, error) {
	// 验证设施存在
	facility, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// 检查名称是否重复（如果要更新名称）
	if req.NameZhHant != nil && *req.NameZhHant != facility.NameZhHant {
		exists, err := s.repo.CheckNameExists(ctx, *req.NameZhHant, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("facility name already exists")
		}
	}

	// 构建更新数据
	updates := make(map[string]interface{})
	if req.NameZhHant != nil {
		updates["name_zh_hant"] = *req.NameZhHant
	}
	if req.NameZhHans != nil {
		updates["name_zh_hans"] = *req.NameZhHans
	}
	if req.NameEn != nil {
		updates["name_en"] = *req.NameEn
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) > 0 {
		if err := s.repo.Update(ctx, id, updates); err != nil {
			return nil, err
		}
	}

	// 重新查询返回最新数据
	updatedFacility, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.FacilityResponse{
		ID:         updatedFacility.ID,
		NameZhHant: updatedFacility.NameZhHant,
		NameZhHans: updatedFacility.NameZhHans,
		NameEn:     updatedFacility.NameEn,
		Icon:       updatedFacility.Icon,
		Category:   updatedFacility.Category,
		SortOrder:  updatedFacility.SortOrder,
		CreatedAt:  updatedFacility.CreatedAt,
		UpdatedAt:  updatedFacility.UpdatedAt,
	}, nil
}

// 5. DeleteFacility 删除设施
func (s *FacilityService) DeleteFacility(ctx context.Context, id uint) error {
	// 验证设施存在
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			return tools.ErrNotFound
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}
