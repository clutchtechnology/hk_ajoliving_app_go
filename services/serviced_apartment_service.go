package services

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
)

// ServicedApartmentService 服务式住宅服务
// Methods:
// 1. ListServicedApartments(ctx context.Context, filter *models.ListServicedApartmentsRequest) -> 获取服务式住宅列表
// 2. GetServicedApartment(ctx context.Context, id uint) -> 获取服务式住宅详情
// 3. CreateServicedApartment(ctx context.Context, req *models.CreateServicedApartmentRequest, companyID uint) -> 创建服务式住宅
// 4. UpdateServicedApartment(ctx context.Context, id uint, req *models.UpdateServicedApartmentRequest, companyID uint) -> 更新服务式住宅
// 5. DeleteServicedApartment(ctx context.Context, id uint, companyID uint) -> 删除服务式住宅
// 6. GetServicedApartmentUnits(ctx context.Context, apartmentID uint) -> 获取房型列表
// 7. GetServicedApartmentImages(ctx context.Context, apartmentID uint) -> 获取图片列表
type ServicedApartmentService struct {
	repo *databases.ServicedApartmentRepo
}

// NewServicedApartmentService 创建服务式住宅服务
func NewServicedApartmentService(repo *databases.ServicedApartmentRepo) *ServicedApartmentService {
	return &ServicedApartmentService{repo: repo}
}

// ListServicedApartments 获取服务式住宅列表
func (s *ServicedApartmentService) ListServicedApartments(ctx context.Context, filter *models.ListServicedApartmentsRequest) (*models.PaginatedServicedApartmentsResponse, error) {
	// 默认参数
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}

	apartments, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]models.ServicedApartmentResponse, len(apartments))
	for i, sa := range apartments {
		items[i] = *sa.ToServicedApartmentResponse()
	}

	// 计算总页数
	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &models.PaginatedServicedApartmentsResponse{
		Data:       items,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetServicedApartment 获取服务式住宅详情
func (s *ServicedApartmentService) GetServicedApartment(ctx context.Context, id uint) (*models.ServicedApartmentDetailResponse, error) {
	apartment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 异步增加浏览次数
	go func() {
		_ = s.repo.IncrementViewCount(context.Background(), id)
	}()

	return apartment.ToServicedApartmentDetailResponse(), nil
}

// CreateServicedApartment 创建服务式住宅
func (s *ServicedApartmentService) CreateServicedApartment(ctx context.Context, req *models.CreateServicedApartmentRequest, companyID uint) (*models.ServicedApartmentDetailResponse, error) {
	apartment := &models.ServicedApartment{
		Name:          req.Name,
		NameEn:        req.NameEn,
		Description:   req.Description,
		Address:       req.Address,
		DistrictID:    req.DistrictID,
		Phone:         req.Phone,
		WebsiteURL:    req.WebsiteURL,
		Email:         req.Email,
		CompanyID:     companyID,
		CheckInTime:   req.CheckInTime,
		CheckOutTime:  req.CheckOutTime,
		MinStayDays:   req.MinStayDays,
		Status:        "active",
		IsFeatured:    false,
		Rating:        0,
		ReviewCount:   0,
		ViewCount:     0,
		FavoriteCount: 0,
	}

	if err := s.repo.Create(ctx, apartment); err != nil {
		return nil, err
	}

	// 返回完整信息
	result, err := s.repo.FindByID(ctx, apartment.ID)
	if err != nil {
		return nil, err
	}

	return result.ToServicedApartmentDetailResponse(), nil
}

// UpdateServicedApartment 更新服务式住宅
func (s *ServicedApartmentService) UpdateServicedApartment(ctx context.Context, id uint, req *models.UpdateServicedApartmentRequest, companyID uint) (*models.ServicedApartmentDetailResponse, error) {
	// 查找现有记录
	apartment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 验证权限（只能修改自己公司的住宅）
	if apartment.CompanyID != companyID {
		return nil, errors.New("permission denied")
	}

	// 更新字段
	if req.Name != nil {
		apartment.Name = *req.Name
	}
	if req.Description != nil {
		apartment.Description = *req.Description
	}
	if req.Phone != nil {
		apartment.Phone = *req.Phone
	}
	if req.WebsiteURL != nil {
		apartment.WebsiteURL = *req.WebsiteURL
	}
	if req.Email != nil {
		apartment.Email = *req.Email
	}
	if req.MinStayDays != nil {
		apartment.MinStayDays = *req.MinStayDays
	}
	if req.CheckInTime != nil {
		apartment.CheckInTime = *req.CheckInTime
	}
	if req.CheckOutTime != nil {
		apartment.CheckOutTime = *req.CheckOutTime
	}
	if req.Status != nil {
		apartment.Status = *req.Status
	}

	if err := s.repo.Update(ctx, apartment); err != nil {
		return nil, err
	}

	// 返回更新后的信息
	result, err := s.repo.FindByID(ctx, apartment.ID)
	if err != nil {
		return nil, err
	}

	return result.ToServicedApartmentDetailResponse(), nil
}

// DeleteServicedApartment 删除服务式住宅
func (s *ServicedApartmentService) DeleteServicedApartment(ctx context.Context, id uint, companyID uint) error {
	// 查找现有记录
	apartment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 验证权限
	if apartment.CompanyID != companyID {
		return errors.New("permission denied")
	}

	return s.repo.Delete(ctx, id)
}

// GetServicedApartmentUnits 获取房型列表
func (s *ServicedApartmentService) GetServicedApartmentUnits(ctx context.Context, apartmentID uint) ([]models.ServicedApartmentUnit, error) {
	// 先检查住宅是否存在
	_, err := s.repo.FindByID(ctx, apartmentID)
	if err != nil {
		return nil, err
	}

	return s.repo.FindUnits(ctx, apartmentID)
}

// GetServicedApartmentImages 获取图片列表
func (s *ServicedApartmentService) GetServicedApartmentImages(ctx context.Context, apartmentID uint) ([]models.ServicedApartmentImage, error) {
	// 先检查住宅是否存在
	_, err := s.repo.FindByID(ctx, apartmentID)
	if err != nil {
		return nil, err
	}

	return s.repo.FindImages(ctx, apartmentID)
}
