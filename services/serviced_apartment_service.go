package services

import (
	"context"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"go.uber.org/zap"
)

// ServicedApartmentService 服务式住宅服务接口
type ServicedApartmentService interface {
	ListServicedApartments(ctx context.Context, req *models.ListServicedApartmentsRequest) ([]*models.ServicedApartment, int64, error)
	GetServicedApartment(ctx context.Context, id uint) (*models.ServicedApartment, error)
	GetApartmentUnits(ctx context.Context, id uint) ([]*models.ServicedApartmentUnit, error)
	GetFeaturedApartments(ctx context.Context, limit int) ([]*models.ServicedApartment, error)
	CreateServicedApartment(ctx context.Context, userID uint, req *models.ServicedApartment) (*models.ServicedApartment, error)
	UpdateServicedApartment(ctx context.Context, id uint, req *models.ServicedApartment) (*models.ServicedApartment, error)
	DeleteServicedApartment(ctx context.Context, id uint) error
}

type servicedApartmentService struct {
	repo   databases.ServicedApartmentRepository
	logger *zap.Logger
}

func NewServicedApartmentService(repo databases.ServicedApartmentRepository, logger *zap.Logger) ServicedApartmentService {
	return &servicedApartmentService{
		repo:   repo,
		logger: logger,
	}
}

func (s *servicedApartmentService) ListServicedApartments(ctx context.Context, req *models.ListServicedApartmentsRequest) ([]*models.ServicedApartment, int64, error) {
	apartments, total, err := s.repo.List(ctx, req)
	if err != nil {
		s.logger.Error("failed to list serviced apartments", zap.Error(err))
		return nil, 0, tools.ErrInternalServer
	}

	result := make([]*models.ServicedApartment, 0, len(apartments))
	for _, apt := range apartments {
		result = append(result, s.toListItemResponse(apt))
	}

	return result, total, nil
}

func (s *servicedApartmentService) GetServicedApartment(ctx context.Context, id uint) (*models.ServicedApartment, error) {
	apartment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get serviced apartment", zap.Uint("id", id), zap.Error(err))
		return nil, tools.ErrNotFound
	}

	// 增加浏览次数
	go func() {
		if err := s.repo.IncrementViewCount(context.Background(), id); err != nil {
			s.logger.Warn("failed to increment view count", zap.Uint("id", id), zap.Error(err))
		}
	}()

	return s.toDetailResponse(apartment), nil
}

func (s *servicedApartmentService) GetApartmentUnits(ctx context.Context, id uint) ([]*models.ServicedApartmentUnit, error) {
	units, err := s.repo.GetUnits(ctx, id)
	if err != nil {
		s.logger.Error("failed to get apartment units", zap.Uint("apartment_id", id), zap.Error(err))
		return nil, tools.ErrInternalServer
	}

	result := make([]*models.ServicedApartmentUnit, 0, len(units))
	for _, unit := range units {
		result = append(result, s.toUnitResponse(&unit))
	}

	return result, nil
}

func (s *servicedApartmentService) GetFeaturedApartments(ctx context.Context, limit int) ([]*models.ServicedApartment, error) {
	apartments, err := s.repo.GetFeatured(ctx, limit)
	if err != nil {
		s.logger.Error("failed to get featured apartments", zap.Error(err))
		return nil, tools.ErrInternalServer
	}

	result := make([]*models.ServicedApartment, 0, len(apartments))
	for _, apt := range apartments {
		result = append(result, s.toListItemResponse(apt))
	}

	return result, nil
}

func (s *servicedApartmentService) CreateServicedApartment(ctx context.Context, userID uint, req *models.ServicedApartment) (*models.ServicedApartment, error) {
	apartment := &models.ServicedApartment{
		Name:         req.Name,
		Description:  req.Description,
		Address:      req.Address,
		DistrictID:   req.DistrictID,
		CompanyID:    userID,
		Phone:        req.Phone,
		Email:        req.Email,
		WebsiteURL:   req.WebsiteURL,
		CheckInTime:  req.CheckInTime,
		CheckOutTime: req.CheckOutTime,
		MinStayDays:  req.MinStayDays,
		Status:       models.ServicedApartmentStatusActive,
	}

	// 创建服务式住宅
	if err := s.repo.Create(ctx, apartment); err != nil {
		s.logger.Error("failed to create serviced apartment", zap.Error(err))
		return nil, tools.ErrInternalServer
	}

	// TODO: 处理图片、设施上传（需要在创建后添加）

	return s.GetServicedApartment(ctx, apartment.ID)
}

func (s *servicedApartmentService) UpdateServicedApartment(ctx context.Context, id uint, req *models.ServicedApartment) (*models.ServicedApartment, error) {
	apartment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, tools.ErrNotFound
	}

	// 更新字段
	if req.Name != "" {
		apartment.Name = req.Name
	}
	if req.Description != nil && *req.Description != "" {
		apartment.Description = req.Description
	}
	if req.Address != "" {
		apartment.Address = req.Address
	}
	if req.DistrictID != 0 {
		apartment.DistrictID = req.DistrictID
	}
	if req.Phone != "" {
		apartment.Phone = req.Phone
	}
	if req.Email != nil && *req.Email != "" {
		apartment.Email = req.Email
	}
	if req.WebsiteURL != nil && *req.WebsiteURL != "" {
		apartment.WebsiteURL = req.WebsiteURL
	}
	if req.CheckInTime != nil && *req.CheckInTime != "" {
		apartment.CheckInTime = req.CheckInTime
	}
	if req.CheckOutTime != nil && *req.CheckOutTime != "" {
		apartment.CheckOutTime = req.CheckOutTime
	}
	if req.MinStayDays != nil && *req.MinStayDays > 0 {
		apartment.MinStayDays = req.MinStayDays
	}
	if req.Status != "" {
		apartment.Status = req.Status
	}
	// IsFeatured 是 bool 类型，不需要指针检查
	apartment.IsFeatured = req.IsFeatured

	if err := s.repo.Update(ctx, apartment); err != nil {
		s.logger.Error("failed to update serviced apartment", zap.Uint("id", id), zap.Error(err))
		return nil, tools.ErrInternalServer
	}

	return s.GetServicedApartment(ctx, id)
}

func (s *servicedApartmentService) DeleteServicedApartment(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete serviced apartment", zap.Uint("id", id), zap.Error(err))
		return tools.ErrInternalServer
	}
	return nil
}

// 转换为列表项响应
func (s *servicedApartmentService) toListItemResponse(apt *models.ServicedApartment) *models.ServicedApartment {
	return apt
}

// 转换为详细响应
func (s *servicedApartmentService) toDetailResponse(apt *models.ServicedApartment) *models.ServicedApartment {
	return apt
}

// 转换为单元响应
func (s *servicedApartmentService) toUnitResponse(unit *models.ServicedApartmentUnit) *models.ServicedApartmentUnit {
	return unit
}
