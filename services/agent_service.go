package services

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"gorm.io/gorm"
)

// AgentService Methods:
// 0. NewAgentService(agentRepo *databases.AgentRepo) -> 注入依赖
// 1. ListAgents(ctx context.Context, req *models.ListAgentsRequest) -> 代理人列表
// 2. GetAgent(ctx context.Context, id uint) -> 代理人详情
// 3. GetAgentProperties(ctx context.Context, agentID uint, page, pageSize int) -> 代理人房源列表
// 4. ContactAgent(ctx context.Context, agentID uint, userID *uint, req *models.ContactAgentRequest) -> 联系代理人

type AgentService struct {
	agentRepo *databases.AgentRepo
}

// 0. NewAgentService 构造函数
func NewAgentService(agentRepo *databases.AgentRepo) *AgentService {
	return &AgentService{
		agentRepo: agentRepo,
	}
}

// 1. ListAgents 代理人列表
func (s *AgentService) ListAgents(ctx context.Context, req *models.ListAgentsRequest) (*models.PaginatedAgentsResponse, error) {
	agents, total, err := s.agentRepo.FindAll(ctx, req)
	if err != nil {
		return nil, err
	}

	var items []*models.AgentResponse
	for _, agent := range agents {
		items = append(items, s.buildAgentResponse(agent))
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &models.PaginatedAgentsResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// 2. GetAgent 代理人详情
func (s *AgentService) GetAgent(ctx context.Context, id uint) (*models.AgentDetailResponse, error) {
	agent, err := s.agentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// 获取服务区域
	serviceAreas, _ := s.agentRepo.GetServiceAreas(ctx, id)

	return s.buildAgentDetailResponse(agent, serviceAreas), nil
}

// 3. GetAgentProperties 代理人房源列表
func (s *AgentService) GetAgentProperties(ctx context.Context, agentID uint, page, pageSize int) (interface{}, error) {
	// 验证代理人是否存在
	_, err := s.agentRepo.FindByID(ctx, agentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	properties, total, err := s.agentRepo.GetAgentProperties(ctx, agentID, page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return map[string]interface{}{
		"items":       properties,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": totalPages,
	}, nil
}

// 4. ContactAgent 联系代理人
func (s *AgentService) ContactAgent(ctx context.Context, agentID uint, userID *uint, req *models.ContactAgentRequest) error {
	// 验证代理人是否存在且状态为活跃
	agent, err := s.agentRepo.FindByID(ctx, agentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tools.ErrNotFound
		}
		return err
	}

	if agent.Status != "active" {
		return errors.New("agent is not active")
	}

	// 创建联系记录
	contact := &models.AgentContact{
		AgentID: agentID,
		UserID:  userID,
		Name:    req.Name,
		Phone:   req.Phone,
		Email:   req.Email,
		Message: req.Message,
	}

	return s.agentRepo.CreateContact(ctx, contact)
}

// buildAgentResponse 构建代理人响应
func (s *AgentService) buildAgentResponse(agent *models.Agent) *models.AgentResponse {
	return &models.AgentResponse{
		ID:               agent.ID,
		AgentName:        agent.AgentName,
		AgentNameEn:      agent.AgentNameEn,
		LicenseNo:        agent.LicenseNo,
		LicenseType:      agent.LicenseType,
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
	}
}

// buildAgentDetailResponse 构建代理人详情响应
func (s *AgentService) buildAgentDetailResponse(agent *models.Agent, serviceAreas []*models.AgentServiceArea) *models.AgentDetailResponse {
	response := &models.AgentDetailResponse{
		ID:                agent.ID,
		AgentName:         agent.AgentName,
		AgentNameEn:       agent.AgentNameEn,
		LicenseNo:         agent.LicenseNo,
		LicenseType:       agent.LicenseType,
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
		Status:            agent.Status,
		IsVerified:        agent.IsVerified,
		VerifiedAt:        agent.VerifiedAt,
	}

	// 添加服务区域
	if len(serviceAreas) > 0 {
		var areas []*models.DistrictBrief
		for _, area := range serviceAreas {
			if area.District != nil {
				areas = append(areas, &models.DistrictBrief{
					ID:         area.District.ID,
					NameZhHant: area.District.NameZhHant,
					NameEn:     area.District.NameEn,
					Region:     area.District.Region,
				})
			}
		}
		response.ServiceAreas = areas
	}

	return response
}
