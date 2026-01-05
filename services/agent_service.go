package services

import (
	"context"
	"go.uber.org/zap"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
)

// AgentService 代理人服务接口
type AgentService interface {
	ListAgents(ctx context.Context, filter *models.ListAgentsRequest) ([]*models.Agent, int64, error)
	GetAgent(ctx context.Context, id uint) (*models.Agent, error)
	GetAgentProperties(ctx context.Context, agentID uint, page, pageSize int) ([]*models.Property, int64, error)
	ContactAgent(ctx context.Context, agentID uint, userID *uint, req *models.ContactAgentRequest) (*models.AgentContactResponse, error)
}

type agentService struct {
	repo   databases.AgentRepository
	logger *zap.Logger
}

// NewAgentService 创建代理人服务
func NewAgentService(repo databases.AgentRepository, logger *zap.Logger) AgentService {
	return &agentService{
		repo:   repo,
		logger: logger,
	}
}

// ListAgents 获取代理人列表
func (s *agentService) ListAgents(ctx context.Context, filter *models.ListAgentsRequest) ([]*models.Agent, int64, error) {
	agents, total, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error("failed to list agents", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*models.Agent, 0, len(agents))
	for _, agent := range agents {
		result = append(result, convertToAgentListItemResponse(agent))
	}
	
	return result, total, nil
}

// GetAgent 获取代理人详情
func (s *agentService) GetAgent(ctx context.Context, id uint) (*models.Agent, error) {
	agent, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to get agent", zap.Error(err))
		return nil, err
	}
	if agent == nil {
		return nil, tools.ErrNotFound
	}
	
	return agent, nil, nil
}

// GetAgentProperties 获取代理人房源列表
func (s *agentService) GetAgentProperties(ctx context.Context, agentID uint, page, pageSize int) ([]*models.Property, int64, error) {
	// 先检查代理人是否存在
	agent, err := s.repo.GetByID(ctx, agentID)
	if err != nil {
		s.logger.Error("failed to get agent", zap.Error(err))
		return nil, 0, err
	}
	if agent == nil {
		return nil, 0, tools.ErrNotFound
	}
	
	// 获取代理人的房源
	properties, total, err := s.repo.GetAgentProperties(ctx, agentID, page, pageSize)
	if err != nil {
		s.logger.Error("failed to get agent properties", zap.Error(err))
		return nil, 0, err
	}
	
	result := make([]*models.Property, 0, len(properties))
	for _, property := range properties {
		result = append(result, convertPropertyToListItemResponse(property))
	}
	
	return result, total, nil
}

// ContactAgent 联系代理人
func (s *agentService) ContactAgent(ctx context.Context, agentID uint, userID *uint, req *models.ContactAgentRequest) (*models.AgentContactResponse, error) {
	// 检查代理人是否存在
	agent, err := s.repo.GetByID(ctx, agentID)
	if err != nil {
		s.logger.Error("failed to get agent", zap.Error(err))
		return nil, err
	}
	if agent == nil {
		return nil, tools.ErrNotFound
	}
	
	// 创建联系请求
	contactReq := &models.AgentContactRequest{
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
	
	return &models.AgentContactResponse{Success: true, Message: "已收到您的咨询"}, nil, nil
}

// 辅助函数

// convertToAgentListItemResponse 转换为代理人列表项响应
func convertToAgentListItemResponse(agent *models.Agent) *models.Agent {
	resp := &models.Agent{
		ID:               agent.ID,
		AgentName:        agent.AgentName,
		AgentNameEn:      agent.AgentNameEn,
		LicenseNo:        agent.LicenseNo,
		LicenseType:      agent.LicenseType,
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
		Status:           agent.Status,
		IsVerified:       agent.IsVerified,
		CreatedAt:        agent.CreatedAt,
	}
	
	// 设置代理公司名称
	if agent.Agency != nil {	
	return resp
}
}
