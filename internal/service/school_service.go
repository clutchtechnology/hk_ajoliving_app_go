package service

import (
	"context"

	"go.uber.org/zap"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
)

// SchoolService 校网和学校服务接口
type SchoolService interface {
	// 校网相关
	ListSchoolNets(ctx context.Context, filter *request.ListSchoolNetsRequest) ([]*response.SchoolNetListItemResponse, int64, error)
	GetSchoolNet(ctx context.Context, id uint) (*response.SchoolNetResponse, error)
	GetSchoolsInNet(ctx context.Context, schoolNetID uint) ([]*response.SchoolListItemResponse, error)
	GetPropertiesInNet(ctx context.Context, schoolNetID uint) ([]*response.PropertyListItemResponse, error)
	GetEstatesInNet(ctx context.Context, schoolNetID uint) ([]*response.EstateListItemResponse, error)
	SearchSchoolNets(ctx context.Context, filter *request.SearchSchoolNetsRequest) ([]*response.SchoolNetListItemResponse, int64, error)
	
	// 学校相关
	ListSchools(ctx context.Context, filter *request.ListSchoolsRequest) ([]*response.SchoolListItemResponse, int64, error)
	GetSchoolNetBySchoolID(ctx context.Context, schoolID uint) (*response.SchoolNetResponse, error)
	SearchSchools(ctx context.Context, filter *request.SearchSchoolsRequest) ([]*response.SchoolListItemResponse, int64, error)
}

type schoolService struct {
	schoolRepo   repository.SchoolRepository
	propertyRepo repository.PropertyRepository
	estateRepo   repository.EstateRepository
	logger       *zap.Logger
}

// NewSchoolService 创建校网和学校服务
func NewSchoolService(
	schoolRepo repository.SchoolRepository,
	propertyRepo repository.PropertyRepository,
	estateRepo repository.EstateRepository,
	logger *zap.Logger,
) SchoolService {
	return &schoolService{
		schoolRepo:   schoolRepo,
		propertyRepo: propertyRepo,
		estateRepo:   estateRepo,
		logger:       logger,
	}
}

// 校网相关

func (s *schoolService) ListSchoolNets(ctx context.Context, filter *request.ListSchoolNetsRequest) ([]*response.SchoolNetListItemResponse, int64, error) {
	schoolNets, total, err := s.schoolRepo.ListSchoolNets(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list school nets", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*response.SchoolNetListItemResponse, 0, len(schoolNets))
	for _, sn := range schoolNets {
		result = append(result, convertToSchoolNetListItemResponse(sn))
	}
	
	return result, total, nil
}

func (s *schoolService) GetSchoolNet(ctx context.Context, id uint) (*response.SchoolNetResponse, error) {
	schoolNet, err := s.schoolRepo.GetSchoolNetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get school net", zap.Error(err))
		return nil, err
	}
	if schoolNet == nil {
		return nil, errors.ErrNotFound
	}
	
	return convertToSchoolNetResponse(schoolNet), nil
}

func (s *schoolService) GetSchoolsInNet(ctx context.Context, schoolNetID uint) ([]*response.SchoolListItemResponse, error) {
	// 检查校网是否存在
	schoolNet, err := s.schoolRepo.GetSchoolNetByID(ctx, schoolNetID)
	if err != nil {
		s.logger.Error("failed to get school net", zap.Error(err))
		return nil, err
	}
	if schoolNet == nil {
		return nil, errors.ErrNotFound
	}
	
	schools, err := s.schoolRepo.GetSchoolsBySchoolNetID(ctx, schoolNetID)
	if err != nil {
		s.logger.Error("failed to get schools in net", zap.Error(err))
		return nil, err
	}
	
	result := make([]*response.SchoolListItemResponse, 0, len(schools))
	for _, school := range schools {
		result = append(result, convertToSchoolListItemResponse(school))
	}
	
	return result, nil
}

func (s *schoolService) GetPropertiesInNet(ctx context.Context, schoolNetID uint) ([]*response.PropertyListItemResponse, error) {
	// 检查校网是否存在
	schoolNet, err := s.schoolRepo.GetSchoolNetByID(ctx, schoolNetID)
	if err != nil {
		s.logger.Error("failed to get school net", zap.Error(err))
		return nil, err
	}
	if schoolNet == nil {
		return nil, errors.ErrNotFound
	}
	
	// 获取校网所在地区的房源
	filter := &request.ListPropertiesRequest{
		DistrictID: &schoolNet.DistrictID,
		Page:       1,
		PageSize:   100,
		SortBy:     "created_at",
		SortOrder:  "desc",
	}
	
	properties, _, err := s.propertyRepo.List(ctx, filter)
	if err != nil {
		s.logger.Error("failed to get properties in net", zap.Error(err))
		return nil, err
	}
	
	result := make([]*response.PropertyListItemResponse, 0, len(properties))
	for _, property := range properties {
		result = append(result, convertToPropertyListItemResponse(property))
	}
	
	return result, nil
}

func (s *schoolService) GetEstatesInNet(ctx context.Context, schoolNetID uint) ([]*response.EstateListItemResponse, error) {
	// 检查校网是否存在
	schoolNet, err := s.schoolRepo.GetSchoolNetByID(ctx, schoolNetID)
	if err != nil {
		s.logger.Error("failed to get school net", zap.Error(err))
		return nil, err
	}
	if schoolNet == nil {
		return nil, errors.ErrNotFound
	}
	
	// 获取校网所在地区的屋苑
	filter := &request.ListEstatesRequest{
		DistrictID: &schoolNet.DistrictID,
		Page:       1,
		PageSize:   100,
		SortBy:     "name_zh_hant",
		SortOrder:  "asc",
	}
	
	estates, _, err := s.estateRepo.List(ctx, filter)
	if err != nil {
		s.logger.Error("failed to get estates in net", zap.Error(err))
		return nil, err
	}
	
	result := make([]*response.EstateListItemResponse, 0, len(estates))
	for _, estate := range estates {
		result = append(result, convertToEstateListItemResponse(estate))
	}
	
	return result, nil
}

func (s *schoolService) SearchSchoolNets(ctx context.Context, filter *request.SearchSchoolNetsRequest) ([]*response.SchoolNetListItemResponse, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	
	schoolNets, total, err := s.schoolRepo.SearchSchoolNets(ctx, filter.Keyword, pageSize, offset)
	if err != nil {
		s.logger.Error("failed to search school nets", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*response.SchoolNetListItemResponse, 0, len(schoolNets))
	for _, sn := range schoolNets {
		result = append(result, convertToSchoolNetListItemResponse(sn))
	}
	
	return result, total, nil
}

// 学校相关

func (s *schoolService) ListSchools(ctx context.Context, filter *request.ListSchoolsRequest) ([]*response.SchoolListItemResponse, int64, error) {
	schools, total, err := s.schoolRepo.ListSchools(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list schools", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*response.SchoolListItemResponse, 0, len(schools))
	for _, school := range schools {
		result = append(result, convertToSchoolListItemResponse(school))
	}
	
	return result, total, nil
}

func (s *schoolService) GetSchoolNetBySchoolID(ctx context.Context, schoolID uint) (*response.SchoolNetResponse, error) {
	school, err := s.schoolRepo.GetSchoolByID(ctx, schoolID)
	if err != nil {
		s.logger.Error("failed to get school", zap.Error(err))
		return nil, err
	}
	if school == nil {
		return nil, errors.ErrNotFound
	}
	
	schoolNet, err := s.schoolRepo.GetSchoolNetByID(ctx, school.SchoolNetID)
	if err != nil {
		s.logger.Error("failed to get school net", zap.Error(err))
		return nil, err
	}
	if schoolNet == nil {
		return nil, errors.ErrNotFound
	}
	
	return convertToSchoolNetResponse(schoolNet), nil
}

func (s *schoolService) SearchSchools(ctx context.Context, filter *request.SearchSchoolsRequest) ([]*response.SchoolListItemResponse, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	
	schools, total, err := s.schoolRepo.SearchSchools(ctx, filter.Keyword, pageSize, offset)
	if err != nil {
		s.logger.Error("failed to search schools", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*response.SchoolListItemResponse, 0, len(schools))
	for _, school := range schools {
		result = append(result, convertToSchoolListItemResponse(school))
	}
	
	return result, total, nil
}

// 辅助函数

func convertToSchoolNetListItemResponse(sn *model.SchoolNet) *response.SchoolNetListItemResponse {
	resp := &response.SchoolNetListItemResponse{
		ID:          sn.ID,
		NetCode:     sn.NetCode,
		NameZhHant:  sn.NameZhHant,
		NameZhHans:  sn.NameZhHans,
		NameEn:      sn.NameEn,
		DistrictID:  sn.DistrictID,
		Level:       sn.Level,
		SchoolCount: sn.SchoolCount,
		CreatedAt:   sn.CreatedAt,
	}
	
	if sn.District != nil {
		resp.District = &response.DistrictBasicResponse{
			ID:         sn.District.ID,
			NameZhHant: sn.District.NameZhHant,
			NameZhHans: sn.District.NameZhHans,
			NameEn:     sn.District.NameEn,
		}
	}
	
	return resp
}

func convertToSchoolNetResponse(sn *model.SchoolNet) *response.SchoolNetResponse {
	resp := &response.SchoolNetResponse{
		ID:          sn.ID,
		NetCode:     sn.NetCode,
		NameZhHant:  sn.NameZhHant,
		NameZhHans:  sn.NameZhHans,
		NameEn:      sn.NameEn,
		DistrictID:  sn.DistrictID,
		Description: sn.Description,
		Level:       sn.Level,
		SchoolCount: sn.SchoolCount,
		MapData:     sn.MapData,
		CreatedAt:   sn.CreatedAt,
		UpdatedAt:   sn.UpdatedAt,
	}
	
	if sn.District != nil {
		resp.District = &response.DistrictBasicResponse{
			ID:         sn.District.ID,
			NameZhHant: sn.District.NameZhHant,
			NameZhHans: sn.District.NameZhHans,
			NameEn:     sn.District.NameEn,
		}
	}
	
	return resp
}

func convertToSchoolListItemResponse(school *model.School) *response.SchoolListItemResponse {
	resp := &response.SchoolListItemResponse{
		ID:           school.ID,
		SchoolNetID:  school.SchoolNetID,
		DistrictID:   school.DistrictID,
		NameZhHant:   school.NameZhHant,
		NameZhHans:   school.NameZhHans,
		NameEn:       school.NameEn,
		SchoolCode:   school.SchoolCode,
		Category:     school.Category,
		CategoryName: school.GetCategoryName(),
		Level:        school.Level,
		LevelName:    school.GetLevelName(),
		Gender:       school.Gender,
		Religion:     school.Religion,
		Address:      school.Address,
		Phone:        school.Phone,
		Website:      school.Website,
		StudentCount: school.StudentCount,
		Rating:       school.Rating,
		LogoURL:      school.LogoURL,
	}
	
	if school.District != nil {
		resp.District = &response.DistrictBasicResponse{
			ID:         school.District.ID,
			NameZhHant: school.District.NameZhHant,
			NameZhHans: school.District.NameZhHans,
			NameEn:     school.District.NameEn,
		}
	}
	
	return resp
}

// 这些转换函数需要在对应的 service 中实现，这里为了编译通过临时添加
func convertToPropertyListItemResponse(property *model.Property) *response.PropertyListItemResponse {
	// 这个函数应该在 property_service.go 中实现
	return &response.PropertyListItemResponse{}
}

func convertToEstateListItemResponse(estate *model.Estate) *response.EstateListItemResponse {
	// 这个函数应该在 estate_service.go 中实现
	return &response.EstateListItemResponse{}
}
