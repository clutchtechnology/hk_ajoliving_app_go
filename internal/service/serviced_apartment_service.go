package service

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
	"go.uber.org/zap"
)

// ServicedApartmentService 服务式住宅服务接口
type ServicedApartmentService interface {
	ListServicedApartments(ctx context.Context, req *request.ListServicedApartmentsRequest) ([]*response.ServicedApartmentListItemResponse, int64, error)
	GetServicedApartment(ctx context.Context, id uint) (*response.ServicedApartmentResponse, error)
	GetApartmentUnits(ctx context.Context, id uint) ([]*response.ServicedApartmentUnitResponse, error)
	GetFeaturedApartments(ctx context.Context, limit int) ([]*response.ServicedApartmentListItemResponse, error)
	CreateServicedApartment(ctx context.Context, userID uint, req *request.CreateServicedApartmentRequest) (*response.ServicedApartmentResponse, error)
	UpdateServicedApartment(ctx context.Context, id uint, req *request.UpdateServicedApartmentRequest) (*response.ServicedApartmentResponse, error)
	DeleteServicedApartment(ctx context.Context, id uint) error
}

type servicedApartmentService struct {
	repo   repository.ServicedApartmentRepository
	logger *zap.Logger
}

func NewServicedApartmentService(repo repository.ServicedApartmentRepository, logger *zap.Logger) ServicedApartmentService {
	return &servicedApartmentService{
		repo:   repo,
		logger: logger,
	}
}

func (s *servicedApartmentService) ListServicedApartments(ctx context.Context, req *request.ListServicedApartmentsRequest) ([]*response.ServicedApartmentListItemResponse, int64, error) {
	apartments, total, err := s.repo.List(ctx, req)
	if err != nil {
		s.logger.Error("failed to list serviced apartments", zap.Error(err))
		return nil, 0, errors.ErrInternalServer
	}

	result := make([]*response.ServicedApartmentListItemResponse, 0, len(apartments))
	for _, apt := range apartments {
		result = append(result, s.toListItemResponse(apt))
	}

	return result, total, nil
}

func (s *servicedApartmentService) GetServicedApartment(ctx context.Context, id uint) (*response.ServicedApartmentResponse, error) {
	apartment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get serviced apartment", zap.Uint("id", id), zap.Error(err))
		return nil, errors.ErrNotFound
	}

	// 增加浏览次数
	go func() {
		if err := s.repo.IncrementViewCount(context.Background(), id); err != nil {
			s.logger.Warn("failed to increment view count", zap.Uint("id", id), zap.Error(err))
		}
	}()

	return s.toDetailResponse(apartment), nil
}

func (s *servicedApartmentService) GetApartmentUnits(ctx context.Context, id uint) ([]*response.ServicedApartmentUnitResponse, error) {
	units, err := s.repo.GetUnits(ctx, id)
	if err != nil {
		s.logger.Error("failed to get apartment units", zap.Uint("apartment_id", id), zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	result := make([]*response.ServicedApartmentUnitResponse, 0, len(units))
	for _, unit := range units {
		result = append(result, s.toUnitResponse(&unit))
	}

	return result, nil
}

func (s *servicedApartmentService) GetFeaturedApartments(ctx context.Context, limit int) ([]*response.ServicedApartmentListItemResponse, error) {
	apartments, err := s.repo.GetFeatured(ctx, limit)
	if err != nil {
		s.logger.Error("failed to get featured apartments", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	result := make([]*response.ServicedApartmentListItemResponse, 0, len(apartments))
	for _, apt := range apartments {
		result = append(result, s.toListItemResponse(apt))
	}

	return result, nil
}

func (s *servicedApartmentService) CreateServicedApartment(ctx context.Context, userID uint, req *request.CreateServicedApartmentRequest) (*response.ServicedApartmentResponse, error) {
	apartment := &model.ServicedApartment{
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
		Status:       model.ServicedApartmentStatusActive,
	}

	// 创建服务式住宅
	if err := s.repo.Create(ctx, apartment); err != nil {
		s.logger.Error("failed to create serviced apartment", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	// TODO: 处理图片、设施上传（需要在创建后添加）

	return s.GetServicedApartment(ctx, apartment.ID)
}

func (s *servicedApartmentService) UpdateServicedApartment(ctx context.Context, id uint, req *request.UpdateServicedApartmentRequest) (*response.ServicedApartmentResponse, error) {
	apartment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrNotFound
	}

	// 更新字段
	if req.Name != nil {
		apartment.Name = *req.Name
	}
	if req.Description != nil {
		apartment.Description = *req.Description
	}
	if req.Address != nil {
		apartment.Address = *req.Address
	}
	if req.DistrictID != nil {
		apartment.DistrictID = *req.DistrictID
	}
	if req.Phone != nil {
		apartment.Phone = *req.Phone
	}
	if req.Email != nil {
		apartment.Email = *req.Email
	}
	if req.WebsiteURL != nil {
		apartment.WebsiteURL = *req.WebsiteURL
	}
	if req.CheckInTime != nil {
		apartment.CheckInTime = *req.CheckInTime
	}
	if req.CheckOutTime != nil {
		apartment.CheckOutTime = *req.CheckOutTime
	}
	if req.MinStayDays != nil {
		apartment.MinStayDays = *req.MinStayDays
	}
	if req.Status != nil {
		apartment.Status = *req.Status
	}
	if req.IsFeatured != nil {
		apartment.IsFeatured = *req.IsFeatured
	}

	if err := s.repo.Update(ctx, apartment); err != nil {
		s.logger.Error("failed to update serviced apartment", zap.Uint("id", id), zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	return s.GetServicedApartment(ctx, id)
}

func (s *servicedApartmentService) DeleteServicedApartment(ctx context.Context, id uint) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete serviced apartment", zap.Uint("id", id), zap.Error(err))
		return errors.ErrInternalServer
	}
	return nil
}

// 转换为列表项响应
func (s *servicedApartmentService) toListItemResponse(apt *model.ServicedApartment) *response.ServicedApartmentListItemResponse {
	resp := &response.ServicedApartmentListItemResponse{
		ID:           apt.ID,
		Name:         apt.Name,
		Address:      apt.Address,
		DistrictName: "",
		MinStayDays:  apt.MinStayDays,
		Status:       apt.Status,
		Rating:       apt.Rating,
		ReviewCount:  apt.ReviewCount,
		ViewCount:    apt.ViewCount,
		IsFeatured:   apt.IsFeatured,
		CoverImage:   "",
	}

	if apt.District != nil {
		resp.DistrictName = apt.District.Name
	}

	if len(apt.Images) > 0 {
		resp.CoverImage = apt.Images[0].ImageURL
	}

	// 计算价格范围
	minPrice, maxPrice := apt.GetMinPrice()
	resp.MinPrice = minPrice
	resp.MaxPrice = maxPrice

	return resp
}

// 转换为详细响应
func (s *servicedApartmentService) toDetailResponse(apt *model.ServicedApartment) *response.ServicedApartmentResponse {
	resp := &response.ServicedApartmentResponse{
		ID:           apt.ID,
		Name:         apt.Name,
		Description:  apt.Description,
		Address:      apt.Address,
		DistrictID:   apt.DistrictID,
		DistrictName: "",
		Phone:        apt.Phone,
		Email:        apt.Email,
		WebsiteURL:   apt.WebsiteURL,
		CheckInTime:  apt.CheckInTime,
		CheckOutTime: apt.CheckOutTime,
		MinStayDays:  apt.MinStayDays,
		Status:       apt.Status,
		Rating:       apt.Rating,
		ReviewCount:  apt.ReviewCount,
		ViewCount:    apt.ViewCount,
		IsFeatured:   apt.IsFeatured,
		CreatedAt:    apt.CreatedAt,
		UpdatedAt:    apt.UpdatedAt,
		Images:       []response.ServicedApartmentImageResponse{},
		Facilities:   []string{},
	}

	if apt.District != nil {
		resp.DistrictName = apt.District.Name
	}

	// 图片
	for _, img := range apt.Images {
		resp.Images = append(resp.Images, response.ServicedApartmentImageResponse{
			ID:        img.ID,
			ImageURL:  img.ImageURL,
			ImageType: img.ImageType,
			SortOrder: img.SortOrder,
		})
	}

	// 设施
	for _, facility := range apt.Facilities {
		resp.Facilities = append(resp.Facilities, facility.Name)
	}

	// 价格范围
	minPrice, maxPrice := apt.GetMinPrice()
	resp.MinPrice = minPrice
	resp.MaxPrice = maxPrice

	return resp
}

// 转换为单元响应
func (s *servicedApartmentService) toUnitResponse(unit *model.ServicedApartmentUnit) *response.ServicedApartmentUnitResponse {
	return &response.ServicedApartmentUnitResponse{
		ID:              unit.ID,
		UnitType:        unit.UnitType,
		Area:            unit.Area,
		Bedrooms:        unit.Bedrooms,
		Bathrooms:       unit.Bathrooms,
		MaxOccupancy:    unit.MaxOccupancy,
		DailyRate:       unit.DailyRate,
		WeeklyRate:      unit.WeeklyRate,
		MonthlyRate:     unit.MonthlyRate,
		IsAvailable:     unit.IsAvailable,
		AvailableFrom:   unit.AvailableFrom,
		Features:        unit.Features,
		FloorPlanURL:    unit.FloorPlanURL,
	}
}
