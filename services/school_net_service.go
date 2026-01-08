package services

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"gorm.io/gorm"
)

// SchoolNetService Methods:
// 0. NewSchoolNetService(schoolNetRepo *databases.SchoolNetRepo) -> 注入依赖
// 1. ListSchoolNets(ctx context.Context, req *models.ListSchoolNetsRequest) -> 校网列表
// 2. GetSchoolNet(ctx context.Context, id uint) -> 校网详情
// 3. GetSchoolsInNet(ctx context.Context, netID uint, page, pageSize int) -> 校网内学校
// 4. GetPropertiesInNet(ctx context.Context, netID uint, page, pageSize int) -> 校网内房源
// 5. GetEstatesInNet(ctx context.Context, netID uint, page, pageSize int) -> 校网内屋苑
// 6. SearchSchoolNets(ctx context.Context, keyword string, page, pageSize int) -> 搜索校网

type SchoolNetService struct {
	schoolNetRepo *databases.SchoolNetRepo
}

// 0. NewSchoolNetService 构造函数
func NewSchoolNetService(schoolNetRepo *databases.SchoolNetRepo) *SchoolNetService {
	return &SchoolNetService{
		schoolNetRepo: schoolNetRepo,
	}
}

// 1. ListSchoolNets 校网列表
func (s *SchoolNetService) ListSchoolNets(ctx context.Context, req *models.ListSchoolNetsRequest) (*models.PaginatedSchoolNetsResponse, error) {
	schoolNets, total, err := s.schoolNetRepo.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}

	var items []*models.SchoolNetResponse
	for _, net := range schoolNets {
		items = append(items, s.buildSchoolNetResponse(net))
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &models.PaginatedSchoolNetsResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// 2. GetSchoolNet 校网详情
func (s *SchoolNetService) GetSchoolNet(ctx context.Context, id uint) (*models.SchoolNetDetailResponse, error) {
	schoolNet, err := s.schoolNetRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	return s.buildSchoolNetDetailResponse(schoolNet), nil
}

// 3. GetSchoolsInNet 校网内学校
func (s *SchoolNetService) GetSchoolsInNet(ctx context.Context, netID uint, page, pageSize int) (*models.PaginatedSchoolsResponse, error) {
	schools, total, err := s.schoolNetRepo.GetSchoolsInNet(ctx, netID, page, pageSize)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	var items []*models.SchoolResponse
	for _, school := range schools {
		items = append(items, s.buildSchoolResponse(school))
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &models.PaginatedSchoolsResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// 4. GetPropertiesInNet 校网内房源
func (s *SchoolNetService) GetPropertiesInNet(ctx context.Context, netID uint, page, pageSize int) (interface{}, error) {
	// 获取校网信息确定类型
	schoolNet, err := s.schoolNetRepo.FindByID(ctx, netID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	properties, total, err := s.schoolNetRepo.GetPropertiesInNet(ctx, netID, schoolNet.Type, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return map[string]interface{}{
		"items":       properties,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	}, nil
}

// 5. GetEstatesInNet 校网内屋苑
func (s *SchoolNetService) GetEstatesInNet(ctx context.Context, netID uint, page, pageSize int) (interface{}, error) {
	// 获取校网信息确定类型
	schoolNet, err := s.schoolNetRepo.FindByID(ctx, netID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	estates, total, err := s.schoolNetRepo.GetEstatesInNet(ctx, netID, schoolNet.Type, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return map[string]interface{}{
		"items":       estates,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	}, nil
}

// 6. SearchSchoolNets 搜索校网
func (s *SchoolNetService) SearchSchoolNets(ctx context.Context, keyword string, page, pageSize int) (*models.PaginatedSchoolNetsResponse, error) {
	schoolNets, total, err := s.schoolNetRepo.Search(ctx, keyword, page, pageSize)
	if err != nil {
		return nil, err
	}

	var items []*models.SchoolNetResponse
	for _, net := range schoolNets {
		items = append(items, s.buildSchoolNetResponse(net))
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &models.PaginatedSchoolNetsResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// buildSchoolNetResponse 构建校网响应
func (s *SchoolNetService) buildSchoolNetResponse(net *models.SchoolNet) *models.SchoolNetResponse {
	return &models.SchoolNetResponse{
		ID:          net.ID,
		Code:        net.Code,
		NameZhHant:  net.NameZhHant,
		NameZhHans:  net.NameZhHans,
		NameEn:      net.NameEn,
		Type:        net.Type,
		DistrictID:  net.DistrictID,
		Description: net.Description,
		Coverage:    net.Coverage,
		SchoolCount: net.SchoolCount,
	}
}

// buildSchoolNetDetailResponse 构建校网详情响应
func (s *SchoolNetService) buildSchoolNetDetailResponse(net *models.SchoolNet) *models.SchoolNetDetailResponse {
	response := &models.SchoolNetDetailResponse{
		ID:          net.ID,
		Code:        net.Code,
		NameZhHant:  net.NameZhHant,
		NameZhHans:  net.NameZhHans,
		NameEn:      net.NameEn,
		Type:        net.Type,
		DistrictID:  net.DistrictID,
		Description: net.Description,
		Coverage:    net.Coverage,
		SchoolCount: net.SchoolCount,
	}

	if net.District != nil {
		response.District = net.District.ToDistrictResponse()
	}

	return response
}

// buildSchoolResponse 构建学校响应
func (s *SchoolNetService) buildSchoolResponse(school *models.School) *models.SchoolResponse {
	return &models.SchoolResponse{
		ID:           school.ID,
		NameZhHant:   school.NameZhHant,
		NameZhHans:   school.NameZhHans,
		NameEn:       school.NameEn,
		Type:         school.Type,
		Category:     school.Category,
		Gender:       school.Gender,
		SchoolNetID:  school.SchoolNetID,
		DistrictID:   school.DistrictID,
		Address:      school.Address,
		Rating:       school.Rating,
		StudentCount: school.StudentCount,
	}
}
