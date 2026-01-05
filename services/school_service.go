package services

import (
	"context"
	"go.uber.org/zap"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
)

// SchoolService 校网和学校服务接口
type SchoolService interface {
	// 校网相关
	ListSchoolNets(ctx context.Context, filter *models.ListSchoolNetsRequest) ([]*models.SchoolNet, int64, error)
	GetSchoolNet(ctx context.Context, id uint) (*models.SchoolNet, error)
	GetSchoolsInNet(ctx context.Context, schoolNetID uint) ([]*models.School, error)
	GetPropertiesInNet(ctx context.Context, schoolNetID uint) ([]*models.Property, error)
	GetEstatesInNet(ctx context.Context, schoolNetID uint) ([]*models.Estate, error)
	SearchSchoolNets(ctx context.Context, filter *models.ListSchoolNetsRequest) ([]*models.SchoolNet, int64, error)
	
	// 学校相关
	ListSchools(ctx context.Context, filter *models.ListSchoolsRequest) ([]*models.School, int64, error)
	GetSchoolNetBySchoolID(ctx context.Context, schoolID uint) (*models.SchoolNet, error)
	SearchSchools(ctx context.Context, filter *models.ListSchoolsRequest) ([]*models.School, int64, error)
}

type schoolService struct {
	schoolRepo   databases.SchoolRepository
	propertyRepo databases.PropertyRepository
	estateRepo   databases.EstateRepository
	logger       *zap.Logger
}

// NewSchoolService 创建校网和学校服务
func NewSchoolService(
	schoolRepo databases.SchoolRepository,
	propertyRepo databases.PropertyRepository,
	estateRepo databases.EstateRepository,
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

func (s *schoolService) ListSchoolNets(ctx context.Context, filter *models.ListSchoolNetsRequest) ([]*models.SchoolNet, int64, error) {
	schoolNets, total, err := s.schoolRepo.ListSchoolNets(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list school nets", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*models.SchoolNet, 0, len(schoolNets))
	for _, sn := range schoolNets {
		result = append(result, convertToSchoolNetListItemResponse(sn))
	}
	
	return result, total, nil
}

func (s *schoolService) GetSchoolNet(ctx context.Context, id uint) (*models.SchoolNet, error) {
	schoolNet, err := s.schoolRepo.GetSchoolNetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get school net", zap.Error(err))
		return nil, err
	}
	if schoolNet == nil {
		return nil, tools.ErrNotFound
	}
	
	return convertToSchoolNetResponse(schoolNet), nil
}

func (s *schoolService) GetSchoolsInNet(ctx context.Context, schoolNetID uint) ([]*models.School, error) {
	// 检查校网是否存在
	schoolNet, err := s.schoolRepo.GetSchoolNetByID(ctx, schoolNetID)
	if err != nil {
		s.logger.Error("failed to get school net", zap.Error(err))
		return nil, err
	}
	if schoolNet == nil {
		return nil, tools.ErrNotFound
	}
	
	schools, err := s.schoolRepo.GetSchoolsBySchoolNetID(ctx, schoolNetID)
	if err != nil {
		s.logger.Error("failed to get schools in net", zap.Error(err))
		return nil, err
	}
	
	result := make([]*models.School, 0, len(schools))
	for _, school := range schools {
		result = append(result, convertToSchoolListItemResponse(school))
	}
	
	return result, nil
}

func (s *schoolService) GetPropertiesInNet(ctx context.Context, schoolNetID uint) ([]*models.Property, error) {
	// 检查校网是否存在
	schoolNet, err := s.schoolRepo.GetSchoolNetByID(ctx, schoolNetID)
	if err != nil {
		s.logger.Error("failed to get school net", zap.Error(err))
		return nil, err
	}
	if schoolNet == nil {
		return nil, tools.ErrNotFound
	}
	
	// 获取校网所在地区的房源
	filter := &models.ListPropertiesRequest{
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
	
	return properties, nil
}

func (s *schoolService) GetEstatesInNet(ctx context.Context, schoolNetID uint) ([]*models.Estate, error) {
	// 检查校网是否存在
	schoolNet, err := s.schoolRepo.GetSchoolNetByID(ctx, schoolNetID)
	if err != nil {
		s.logger.Error("failed to get school net", zap.Error(err))
		return nil, err
	}
	if schoolNet == nil {
		return nil, tools.ErrNotFound
	}
	
	// 获取校网所在地区的屋苑
	filter := &models.ListEstatesRequest{
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
	
	return estates, nil
}

func (s *schoolService) SearchSchoolNets(ctx context.Context, filter *models.ListSchoolNetsRequest) ([]*models.SchoolNet, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	
	keyword := ""
	if filter.Keyword != nil {
		keyword = *filter.Keyword
	}
	
	schoolNets, total, err := s.schoolRepo.SearchSchoolNets(ctx, keyword, pageSize, offset)
	if err != nil {
		s.logger.Error("failed to search school nets", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*models.SchoolNet, 0, len(schoolNets))
	for _, sn := range schoolNets {
		result = append(result, convertToSchoolNetListItemResponse(sn))
	}
	
	return result, total, nil
}

// 学校相关

func (s *schoolService) ListSchools(ctx context.Context, filter *models.ListSchoolsRequest) ([]*models.School, int64, error) {
	schools, total, err := s.schoolRepo.ListSchools(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list schools", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*models.School, 0, len(schools))
	for _, school := range schools {
		result = append(result, convertToSchoolListItemResponse(school))
	}
	
	return result, total, nil
}

func (s *schoolService) GetSchoolNetBySchoolID(ctx context.Context, schoolID uint) (*models.SchoolNet, error) {
	school, err := s.schoolRepo.GetSchoolByID(ctx, schoolID)
	if err != nil {
		s.logger.Error("failed to get school", zap.Error(err))
		return nil, err
	}
	if school == nil {
		return nil, tools.ErrNotFound
	}
	
	schoolNet, err := s.schoolRepo.GetSchoolNetByID(ctx, school.SchoolNetID)
	if err != nil {
		s.logger.Error("failed to get school net", zap.Error(err))
		return nil, err
	}
	if schoolNet == nil {
		return nil, tools.ErrNotFound
	}
	
	return convertToSchoolNetResponse(schoolNet), nil
}

func (s *schoolService) SearchSchools(ctx context.Context, filter *models.ListSchoolsRequest) ([]*models.School, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	
	keyword := ""
	if filter.Keyword != nil {
		keyword = *filter.Keyword
	}
	
	schools, total, err := s.schoolRepo.SearchSchools(ctx, keyword, pageSize, offset)
	if err != nil {
		s.logger.Error("failed to search schools", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*models.School, 0, len(schools))
	for _, school := range schools {
		result = append(result, convertToSchoolListItemResponse(school))
	}
	
	return result, total, nil
}

// 辅助函数

func convertToSchoolNetListItemResponse(sn *models.SchoolNet) *models.SchoolNet {
	resp := &models.SchoolNet{
		ID:          sn.ID,
		NetCode:     sn.NetCode,
		NameZhHant:  sn.NameZhHant,
		NameZhHans:  sn.NameZhHans,
		NameEn:      sn.NameEn,
		DistrictID:  sn.DistrictID,
		Level:       sn.Level,
		SchoolCount: sn.SchoolCount,
		CreatedAt:   sn.CreatedAt,
		District:    sn.District,
	}
	
	return resp
}

func convertToSchoolNetResponse(sn *models.SchoolNet) *models.SchoolNet {
	resp := &models.SchoolNet{
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
		District:    sn.District,
	}
	
	return resp
}

func convertToSchoolListItemResponse(school *models.School) *models.School {
	resp := &models.School{
		ID:           school.ID,
		SchoolNetID:  school.SchoolNetID,
		DistrictID:   school.DistrictID,
		NameZhHant:   school.NameZhHant,
		NameZhHans:   school.NameZhHans,
		NameEn:       school.NameEn,
		SchoolCode:   school.SchoolCode,
		Category:     school.Category,
		Level:        school.Level,
		Gender:       school.Gender,
		Religion:     school.Religion,
		Address:      school.Address,
		Phone:        school.Phone,
		Website:      school.Website,
		StudentCount: school.StudentCount,
		Rating:       school.Rating,
		LogoURL:      school.LogoURL,
		District:     school.District,
	}
	
	return resp
}

