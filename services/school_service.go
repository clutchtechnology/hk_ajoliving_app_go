package services

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"gorm.io/gorm"
)

// SchoolService Methods:
// 0. NewSchoolService(schoolRepo *databases.SchoolRepo) -> 注入依赖
// 1. ListSchools(ctx context.Context, req *models.ListSchoolsRequest) -> 学校列表
// 2. GetSchool(ctx context.Context, id uint) -> 学校详情
// 3. GetSchoolNet(ctx context.Context, schoolID uint) -> 获取学校所属校网
// 4. SearchSchools(ctx context.Context, keyword string, page, pageSize int) -> 搜索学校

type SchoolService struct {
	schoolRepo *databases.SchoolRepo
}

// 0. NewSchoolService 构造函数
func NewSchoolService(schoolRepo *databases.SchoolRepo) *SchoolService {
	return &SchoolService{
		schoolRepo: schoolRepo,
	}
}

// 1. ListSchools 学校列表
func (s *SchoolService) ListSchools(ctx context.Context, req *models.ListSchoolsRequest) (*models.PaginatedSchoolsResponse, error) {
	schools, total, err := s.schoolRepo.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}

	var items []*models.SchoolResponse
	for _, school := range schools {
		items = append(items, s.buildSchoolResponse(school))
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &models.PaginatedSchoolsResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// 2. GetSchool 学校详情
func (s *SchoolService) GetSchool(ctx context.Context, id uint) (*models.SchoolDetailResponse, error) {
	school, err := s.schoolRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// 增加浏览次数
	_ = s.schoolRepo.IncrementViewCount(ctx, id)

	return s.buildSchoolDetailResponse(school), nil
}

// 3. GetSchoolNet 获取学校所属校网
func (s *SchoolService) GetSchoolNet(ctx context.Context, schoolID uint) (*models.SchoolNetResponse, error) {
	school, err := s.schoolRepo.FindByID(ctx, schoolID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	if school.SchoolNet == nil {
		return nil, errors.New("school does not belong to any school net")
	}

	return &models.SchoolNetResponse{
		ID:          school.SchoolNet.ID,
		Code:        school.SchoolNet.Code,
		NameZhHant:  school.SchoolNet.NameZhHant,
		NameZhHans:  school.SchoolNet.NameZhHans,
		NameEn:      school.SchoolNet.NameEn,
		Type:        school.SchoolNet.Type,
		DistrictID:  school.SchoolNet.DistrictID,
		Description: school.SchoolNet.Description,
		Coverage:    school.SchoolNet.Coverage,
		SchoolCount: school.SchoolNet.SchoolCount,
	}, nil
}

// 4. SearchSchools 搜索学校
func (s *SchoolService) SearchSchools(ctx context.Context, keyword string, page, pageSize int) (*models.PaginatedSchoolsResponse, error) {
	schools, total, err := s.schoolRepo.Search(ctx, keyword, page, pageSize)
	if err != nil {
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

// buildSchoolResponse 构建学校响应
func (s *SchoolService) buildSchoolResponse(school *models.School) *models.SchoolResponse {
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

// buildSchoolDetailResponse 构建学校详情响应
func (s *SchoolService) buildSchoolDetailResponse(school *models.School) *models.SchoolDetailResponse {
	response := &models.SchoolDetailResponse{
		ID:            school.ID,
		NameZhHant:    school.NameZhHant,
		NameZhHans:    school.NameZhHans,
		NameEn:        school.NameEn,
		Type:          school.Type,
		Category:      school.Category,
		Gender:        school.Gender,
		SchoolNetID:   school.SchoolNetID,
		DistrictID:    school.DistrictID,
		Address:       school.Address,
		Phone:         school.Phone,
		Email:         school.Email,
		Website:       school.Website,
		EstablishedAt: school.EstablishedAt,
		Principal:     school.Principal,
		Religion:      school.Religion,
		Curriculum:    school.Curriculum,
		StudentCount:  school.StudentCount,
		TeacherCount:  school.TeacherCount,
		Rating:        school.Rating,
		Description:   school.Description,
		ViewCount:     school.ViewCount,
	}

	if school.SchoolNet != nil {
		response.SchoolNet = &models.SchoolNetResponse{
			ID:          school.SchoolNet.ID,
			Code:        school.SchoolNet.Code,
			NameZhHant:  school.SchoolNet.NameZhHant,
			NameZhHans:  school.SchoolNet.NameZhHans,
			NameEn:      school.SchoolNet.NameEn,
			Type:        school.SchoolNet.Type,
			DistrictID:  school.SchoolNet.DistrictID,
			Description: school.SchoolNet.Description,
			Coverage:    school.SchoolNet.Coverage,
			SchoolCount: school.SchoolNet.SchoolCount,
		}
	}

	if school.District != nil {
		response.District = school.District.ToDistrictResponse()
	}

	return response
}
