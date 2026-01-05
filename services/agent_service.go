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
	
	return agent, nil
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
	
	return properties, total, nil
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
	
	// TODO: 实现联系请求逻辑
	// 1. 保存联系记录到数据库
	// 2. 发送通知给代理人
	// 3. 发送确认邮件/短信给用户
	
	s.logger.Info("agent contact request received",
		zap.Uint("agent_id", agentID),
		zap.String("name", req.Name),
		zap.String("phone", req.Phone),
	)
	
	return &models.AgentContactResponse{
		Success: true,
		Message: "已收到您的咨询，代理人将尽快与您联系",
	}, nil
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
	
	return resp
}
