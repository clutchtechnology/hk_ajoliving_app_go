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

// AgentService 代理人服务接口
type AgentService interface {
	ListAgents(ctx context.Context, filter *request.ListAgentsRequest) ([]*response.AgentListItemResponse, int64, error)
	GetAgent(ctx context.Context, id uint) (*response.AgentResponse, error)
	GetAgentProperties(ctx context.Context, agentID uint, page, pageSize int) ([]*response.PropertyListItemResponse, int64, error)
	ContactAgent(ctx context.Context, agentID uint, userID *uint, req *request.ContactAgentRequest) (*response.AgentContactResponse, error)
}

type agentService struct {
	repo   repository.AgentRepository
	logger *zap.Logger
}

// NewAgentService 创建代理人服务
func NewAgentService(repo repository.AgentRepository, logger *zap.Logger) AgentService {
	return &agentService{
		repo:   repo,
		logger: logger,
	}
}

// ListAgents 获取代理人列表
func (s *agentService) ListAgents(ctx context.Context, filter *request.ListAgentsRequest) ([]*response.AgentListItemResponse, int64, error) {
	agents, total, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list agents", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*response.AgentListItemResponse, 0, len(agents))
	for _, agent := range agents {
		result = append(result, convertToAgentListItemResponse(agent))
	}
	
	return result, total, nil
}

// GetAgent 获取代理人详情
func (s *agentService) GetAgent(ctx context.Context, id uint) (*response.AgentResponse, error) {
	agent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get agent", zap.Error(err))
		return nil, err
	}
	if agent == nil {
		return nil, errors.ErrNotFound
	}
	
	return convertToAgentResponse(agent), nil
}

// GetAgentProperties 获取代理人房源列表
func (s *agentService) GetAgentProperties(ctx context.Context, agentID uint, page, pageSize int) ([]*response.PropertyListItemResponse, int64, error) {
	// 先检查代理人是否存在
	agent, err := s.repo.GetByID(ctx, agentID)
	if err != nil {
		s.logger.Error("failed to get agent", zap.Error(err))
		return nil, 0, err
	}
	if agent == nil {
		return nil, 0, errors.ErrNotFound
	}
	
	// 获取代理人的房源
	properties, total, err := s.repo.GetAgentProperties(ctx, agentID, page, pageSize)
	if err != nil {
		s.logger.Error("failed to get agent properties", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*response.PropertyListItemResponse, 0, len(properties))
	for _, property := range properties {
		result = append(result, convertToPropertyListItemResponse(property))
	}
	
	return result, total, nil
}

// ContactAgent 联系代理人
func (s *agentService) ContactAgent(ctx context.Context, agentID uint, userID *uint, req *request.ContactAgentRequest) (*response.AgentContactResponse, error) {
	// 检查代理人是否存在
	agent, err := s.repo.GetByID(ctx, agentID)
	if err != nil {
		s.logger.Error("failed to get agent", zap.Error(err))
		return nil, err
	}
	if agent == nil {
		return nil, errors.ErrNotFound
	}
	
	// 创建联系请求
	contactReq := &model.AgentContactRequest{
		AgentID:     agentID,
		UserID:      userID,
		PropertyID:  req.PropertyID,
		Name:        req.Name,
		Phone:       req.Phone,
		Email:       req.Email,
		Message:     req.Message,
		ContactType: req.ContactType,
	}
	
	if err := s.repo.CreateContactRequest(ctx, contactReq); err != nil {
		s.logger.Error("failed to create contact request", zap.Error(err))
		return nil, err
	}
	
	// 重新查询以获取关联数据
	contactReq, err = s.repo.GetContactRequestByID(ctx, contactReq.ID)
	if err != nil {
		s.logger.Error("failed to get contact request", zap.Error(err))
		return nil, err
	}
	
	return convertToAgentContactResponse(contactReq), nil
}

// 辅助函数

// convertToAgentListItemResponse 转换为代理人列表项响应
func convertToAgentListItemResponse(agent *model.Agent) *response.AgentListItemResponse {
	resp := &response.AgentListItemResponse{
		ID:               agent.ID,
		AgentName:        agent.AgentName,
		AgentNameEn:      agent.AgentNameEn,
		LicenseNo:        agent.LicenseNo,
		LicenseType:      string(agent.LicenseType),
		AgencyID:         agent.AgencyID,
		Phone:            agent.Phone,
		Email:            agent.Email,
		ProfilePhoto:     agent.ProfilePhoto,
		Specialization:   agent.Specialization,
		YearsExperience:  agent.YearsExperience,
		Rating:           agent.Rating,
		ReviewCount:      agent.ReviewCount,
		PropertiesSold:   agent.PropertiesSold,
		PropertiesRented: agent.PropertiesRented,
		Status:           string(agent.Status),
		IsVerified:       agent.IsVerified,
		CreatedAt:        agent.CreatedAt,
	}
	
	// 设置代理公司名称
	if agent.Agency != nil {
		agencyName := agent.Agency.Username
		resp.AgencyName = &agencyName
	}
	
	return resp
}

// convertToAgentResponse 转换为代理人详情响应
func convertToAgentResponse(agent *model.Agent) *response.AgentResponse {
	resp := &response.AgentResponse{
		ID:                agent.ID,
		UserID:            agent.UserID,
		AgentName:         agent.AgentName,
		AgentNameEn:       agent.AgentNameEn,
		LicenseNo:         agent.LicenseNo,
		LicenseType:       string(agent.LicenseType),
		LicenseExpiryDate: agent.LicenseExpiryDate,
		AgencyID:          agent.AgencyID,
		Phone:             agent.Phone,
		Mobile:            agent.Mobile,
		Email:             agent.Email,
		WechatID:          agent.WechatID,
		Whatsapp:          agent.Whatsapp,
		OfficeAddress:     agent.OfficeAddress,
		Specialization:    agent.Specialization,
		YearsExperience:   agent.YearsExperience,
		ProfilePhoto:      agent.ProfilePhoto,
		Bio:               agent.Bio,
		Rating:            agent.Rating,
		ReviewCount:       agent.ReviewCount,
		PropertiesSold:    agent.PropertiesSold,
		PropertiesRented:  agent.PropertiesRented,
		Status:            string(agent.Status),
		IsVerified:        agent.IsVerified,
		VerifiedAt:        agent.VerifiedAt,
		CreatedAt:         agent.CreatedAt,
		UpdatedAt:         agent.UpdatedAt,
	}
	
	// 设置代理公司名称
	if agent.Agency != nil {
		agencyName := agent.Agency.Username
		resp.AgencyName = &agencyName
	}
	
	// 设置服务区域
	if len(agent.ServiceAreas) > 0 {
		serviceAreas := make([]response.AgentServiceAreaResponse, 0, len(agent.ServiceAreas))
		for _, sa := range agent.ServiceAreas {
			area := response.AgentServiceAreaResponse{
				ID:         sa.ID,
				DistrictID: sa.DistrictID,
			}
			if sa.District != nil {
				area.DistrictName = sa.District.NameZhHant
			}
			serviceAreas = append(serviceAreas, area)
		}
		resp.ServiceAreas = serviceAreas
	}
	
	return resp
}

// convertToAgentContactResponse 转换为联系请求响应
func convertToAgentContactResponse(contactReq *model.AgentContactRequest) *response.AgentContactResponse {
	resp := &response.AgentContactResponse{
		ID:          contactReq.ID,
		AgentID:     contactReq.AgentID,
		PropertyID:  contactReq.PropertyID,
		Name:        contactReq.Name,
		Phone:       contactReq.Phone,
		Email:       contactReq.Email,
		Message:     contactReq.Message,
		ContactType: contactReq.ContactType,
		Status:      contactReq.Status,
		ContactedAt: contactReq.ContactedAt,
		CreatedAt:   contactReq.CreatedAt,
	}
	
	if contactReq.Agent != nil {
		resp.AgentName = contactReq.Agent.AgentName
	}
	
	return resp
}

// convertToPropertyListItemResponse 转换为房产列表项响应
func convertToPropertyListItemResponse(property *model.Property) *response.PropertyListItemResponse {
	resp := &response.PropertyListItemResponse{
		ID:           property.ID,
		Title:        property.Title,
		Price:        property.Price,
		Area:         property.Area,
		Bedrooms:     property.Bedrooms,
		Bathrooms:    property.Bathrooms,
		PropertyType: property.PropertyType,
		ListingType:  property.ListingType,
		Address:      property.Address,
		DistrictID:   property.DistrictID,
		Status:       property.Status,
		CreatedAt:    property.CreatedAt,
		UpdatedAt:    property.UpdatedAt,
	}
	
	if property.District != nil {
		resp.DistrictName = property.District.NameZhHant
	}
	
	if property.Estate != nil {
		estateName := property.Estate.NameZhHant
		resp.EstateName = &estateName
	}
	
	// 设置主图
	if len(property.Images) > 0 {
		resp.CoverImage = property.Images[0].ImageURL
	}
	
	return resp
}
