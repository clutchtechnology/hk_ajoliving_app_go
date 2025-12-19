package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
)

// AgencyService 代理公司服务接口
//
// AgencyService Methods:
// 0. NewAgencyService(repo repository.AgencyRepository, logger *zap.Logger) -> 注入依赖
// 1. ListAgencies(ctx context.Context, filter *request.ListAgenciesRequest) -> 获取代理公司列表
// 2. GetAgency(ctx context.Context, id uint) -> 获取代理公司详情
// 3. GetAgencyProperties(ctx context.Context, agencyID uint, page, pageSize int) -> 获取代理公司房源列表
// 4. ContactAgency(ctx context.Context, agencyID uint, req *request.ContactAgencyRequest) -> 联系代理公司
// 5. SearchAgencies(ctx context.Context, filter *request.SearchAgenciesRequest) -> 搜索代理公司
type AgencyService interface {
	ListAgencies(ctx context.Context, filter *request.ListAgenciesRequest) ([]*response.AgencyListItemResponse, int64, error)
	GetAgency(ctx context.Context, id uint) (*response.AgencyResponse, error)
	GetAgencyProperties(ctx context.Context, agencyID uint, page, pageSize int) ([]*response.PropertyListItemResponse, int64, error)
	ContactAgency(ctx context.Context, agencyID uint, req *request.ContactAgencyRequest) (*response.ContactAgencyResponse, error)
	SearchAgencies(ctx context.Context, filter *request.SearchAgenciesRequest) ([]*response.AgencyListItemResponse, int64, error)
}

type agencyService struct {
	repo   repository.AgencyRepository
	logger *zap.Logger
}

// 0. NewAgencyService 创建代理公司服务
func NewAgencyService(repo repository.AgencyRepository, logger *zap.Logger) AgencyService {
	return &agencyService{
		repo:   repo,
		logger: logger,
	}
}

// 1. ListAgencies 获取代理公司列表
func (s *agencyService) ListAgencies(ctx context.Context, filter *request.ListAgenciesRequest) ([]*response.AgencyListItemResponse, int64, error) {
	agencies, total, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list agencies", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*response.AgencyListItemResponse, 0, len(agencies))
	for _, agency := range agencies {
		// 获取房源数量
		propertyCount, err := s.repo.GetPropertyCount(ctx, agency.ID)
		if err != nil {
			s.logger.Warn("failed to get property count", zap.Error(err), zap.Uint("agency_id", agency.ID))
			propertyCount = 0
		}
		
		result = append(result, convertToAgencyListItemResponse(agency, propertyCount))
	}
	
	return result, total, nil
}

// 2. GetAgency 获取代理公司详情
func (s *agencyService) GetAgency(ctx context.Context, id uint) (*response.AgencyResponse, error) {
	agency, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get agency", zap.Error(err), zap.Uint("agency_id", id))
		return nil, err
	}
	if agency == nil {
		return nil, errors.ErrNotFound
	}
	
	// 获取房源数量
	propertyCount, err := s.repo.GetPropertyCount(ctx, id)
	if err != nil {
		s.logger.Warn("failed to get property count", zap.Error(err), zap.Uint("agency_id", id))
		propertyCount = 0
	}
	
	// 获取优秀代理人
	topAgents, err := s.repo.GetTopAgents(ctx, id, 5)
	if err != nil {
		s.logger.Warn("failed to get top agents", zap.Error(err), zap.Uint("agency_id", id))
		topAgents = []*model.Agent{}
	}
	
	return convertToAgencyResponse(agency, propertyCount, topAgents), nil
}

// 3. GetAgencyProperties 获取代理公司房源列表
func (s *agencyService) GetAgencyProperties(ctx context.Context, agencyID uint, page, pageSize int) ([]*response.PropertyListItemResponse, int64, error) {
	// 先检查代理公司是否存在
	agency, err := s.repo.GetByID(ctx, agencyID)
	if err != nil {
		s.logger.Error("failed to get agency", zap.Error(err), zap.Uint("agency_id", agencyID))
		return nil, 0, err
	}
	if agency == nil {
		return nil, 0, errors.ErrNotFound
	}
	
	// 获取代理公司的房源
	properties, total, err := s.repo.GetAgencyProperties(ctx, agencyID, page, pageSize)
	if err != nil {
		s.logger.Error("failed to get agency properties", zap.Error(err), zap.Uint("agency_id", agencyID))
		return nil, 0, err
	}
	
	result := make([]*response.PropertyListItemResponse, 0, len(properties))
	for _, property := range properties {
		result = append(result, convertToPropertyListItemResponse(property))
	}
	
	return result, total, nil
}

// 4. ContactAgency 联系代理公司
func (s *agencyService) ContactAgency(ctx context.Context, agencyID uint, req *request.ContactAgencyRequest) (*response.ContactAgencyResponse, error) {
	// 检查代理公司是否存在
	agency, err := s.repo.GetByID(ctx, agencyID)
	if err != nil {
		s.logger.Error("failed to get agency", zap.Error(err), zap.Uint("agency_id", agencyID))
		return nil, err
	}
	if agency == nil {
		return nil, errors.ErrNotFound
	}
	
	// TODO: 这里可以实现联系请求的逻辑，比如：
	// 1. 发送邮件通知代理公司
	// 2. 记录联系请求到数据库
	// 3. 发送短信通知
	// 4. 创建待办事项
	
	s.logger.Info("agency contact request received",
		zap.Uint("agency_id", agencyID),
		zap.String("name", req.Name),
		zap.String("phone", req.Phone),
		zap.String("subject", req.Subject),
	)
	
	// 返回成功响应
	return &response.ContactAgencyResponse{
		Success: true,
		Message: fmt.Sprintf("已收到您的咨询，%s 将尽快与您联系", agency.CompanyName),
	}, nil
}

// 5. SearchAgencies 搜索代理公司
func (s *agencyService) SearchAgencies(ctx context.Context, filter *request.SearchAgenciesRequest) ([]*response.AgencyListItemResponse, int64, error) {
	agencies, total, err := s.repo.Search(ctx, filter)
	if err != nil {
		s.logger.Error("failed to search agencies", zap.Error(err), zap.String("keyword", filter.Keyword))
		return nil, 0, err
	}
	
	result := make([]*response.AgencyListItemResponse, 0, len(agencies))
	for _, agency := range agencies {
		// 获取房源数量
		propertyCount, err := s.repo.GetPropertyCount(ctx, agency.ID)
		if err != nil {
			s.logger.Warn("failed to get property count", zap.Error(err), zap.Uint("agency_id", agency.ID))
			propertyCount = 0
		}
		
		result = append(result, convertToAgencyListItemResponse(agency, propertyCount))
	}
	
	return result, total, nil
}

// ========== 转换函数 ==========

// convertToAgencyListItemResponse 转换为代理公司列表项响应
func convertToAgencyListItemResponse(agency *model.AgencyDetail, propertyCount int) *response.AgencyListItemResponse {
	return &response.AgencyListItemResponse{
		ID:              agency.ID,
		CompanyName:     agency.CompanyName,
		CompanyNameEn:   agency.CompanyNameEn,
		LogoURL:         agency.LogoURL,
		Address:         agency.Address,
		Phone:           agency.Phone,
		Email:           agency.Email,
		AgentCount:      agency.AgentCount,
		Rating:          agency.Rating,
		ReviewCount:     agency.ReviewCount,
		IsVerified:      agency.IsVerified,
		EstablishedYear: agency.EstablishedYear,
		PropertyCount:   propertyCount,
	}
}

// convertToAgencyResponse 转换为代理公司详情响应
func convertToAgencyResponse(agency *model.AgencyDetail, propertyCount int, topAgents []*model.Agent) *response.AgencyResponse {
	// 转换优秀代理人信息
	agentInfos := make([]response.AgentBasicInfo, 0, len(topAgents))
	for _, agent := range topAgents {
		agentInfos = append(agentInfos, response.AgentBasicInfo{
			ID:           agent.ID,
			AgentName:    agent.AgentName,
			ProfilePhoto: agent.ProfilePhoto,
			Phone:        agent.Phone,
			Email:        agent.Email,
			Rating:       agent.Rating,
			ReviewCount:  agent.ReviewCount,
		})
	}
	
	return &response.AgencyResponse{
		ID:                     agency.ID,
		CompanyName:            agency.CompanyName,
		CompanyNameEn:          agency.CompanyNameEn,
		LicenseNo:              agency.LicenseNo,
		BusinessRegistrationNo: agency.BusinessRegistrationNo,
		Address:                agency.Address,
		Phone:                  agency.Phone,
		Fax:                    agency.Fax,
		Email:                  agency.Email,
		WebsiteURL:             agency.WebsiteURL,
		EstablishedYear:        agency.EstablishedYear,
		AgentCount:             agency.AgentCount,
		Description:            agency.Description,
		LogoURL:                agency.LogoURL,
		CoverImageURL:          agency.CoverImageURL,
		Rating:                 agency.Rating,
		ReviewCount:            agency.ReviewCount,
		IsVerified:             agency.IsVerified,
		VerifiedAt:             agency.VerifiedAt,
		PropertyCount:          propertyCount,
		TopAgents:              agentInfos,
		ServiceDistricts:       []response.DistrictInfo{}, // TODO: 如果需要显示服务地区，需要从数据库查询
		CreatedAt:              agency.CreatedAt,
		UpdatedAt:              agency.UpdatedAt,
	}
}

// convertToPropertyListItemResponse 转换为房产列表项响应
func convertToPropertyListItemResponse(property *model.Property) *response.PropertyListItemResponse {
	resp := &response.PropertyListItemResponse{
		ID:           property.ID,
		PropertyNo:   property.PropertyNo,
		ListingType:  property.ListingType,
		Title:        property.Title,
		Price:        property.Price,
		Area:         property.Area,
		Address:      property.Address,
		DistrictID:   property.DistrictID,
		BuildingName: property.BuildingName,
		Bedrooms:     property.Bedrooms,
		Bathrooms:    property.Bathrooms,
		PropertyType: property.PropertyType,
		Status:       property.Status,
		ViewCount:    property.ViewCount,
		FavoriteCount: property.FavoriteCount,
		CreatedAt:    property.CreatedAt,
	}
	
	// 设置地区信息
	if property.District != nil {
		resp.District = &response.DistrictResponse{
			ID:         property.District.ID,
			NameZhHant: property.District.NameZhHant,
			Region:     string(property.District.Region),
		}
		if property.District.NameZhHans != nil {
			resp.District.NameZhHans = *property.District.NameZhHans
		}
		if property.District.NameEn != nil {
			resp.District.NameEn = *property.District.NameEn
		}
	}
	
	// 设置主图
	if len(property.Images) > 0 {
		for _, img := range property.Images {
			if img.IsCover {
				resp.CoverImage = img.URL
				break
			}
		}
		// 如果没有设置封面，使用第一张图片
		if resp.CoverImage == "" && len(property.Images) > 0 {
			resp.CoverImage = property.Images[0].URL
		}
	}
	
	return resp
}
