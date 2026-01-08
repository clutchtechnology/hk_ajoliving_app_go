package services

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
)

// AgencyService Methods:
// 0. NewAgencyService(repo *databases.AgencyRepo) -> 注入依赖
// 1. ListAgencies(ctx context.Context, filter *models.ListAgenciesRequest) -> 获取代理公司列表
// 2. GetAgency(ctx context.Context, id uint) -> 获取代理公司详情
// 3. GetAgencyProperties(ctx context.Context, id uint, page, pageSize int) -> 获取代理公司房源列表
// 4. ContactAgency(ctx context.Context, agencyID uint, userID *uint, req *models.ContactAgencyRequest) -> 联系代理公司
// 5. SearchAgencies(ctx context.Context, req *models.SearchAgenciesRequest) -> 搜索代理公司

type AgencyService struct {
	repo *databases.AgencyRepo
}

// 0. NewAgencyService 构造函数
func NewAgencyService(repo *databases.AgencyRepo) *AgencyService {
	return &AgencyService{repo: repo}
}

// 1. ListAgencies 获取代理公司列表
func (s *AgencyService) ListAgencies(ctx context.Context, filter *models.ListAgenciesRequest) (*models.PaginatedAgenciesResponse, error) {
	agencies, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	var agencyResponses []models.AgencyResponse
	for _, agency := range agencies {
		agencyResponses = append(agencyResponses, models.AgencyResponse{
			ID:              agency.ID,
			CompanyName:     agency.CompanyName,
			CompanyNameEn:   agency.CompanyNameEn,
			LogoURL:         agency.LogoURL,
			Address:         agency.Address,
			Phone:           agency.Phone,
			Email:           agency.Email,
			WebsiteURL:      agency.WebsiteURL,
			EstablishedYear: agency.EstablishedYear,
			AgentCount:      agency.AgentCount,
			Rating:          agency.Rating,
			ReviewCount:     agency.ReviewCount,
			IsVerified:      agency.IsVerified,
		})
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize > 0 {
		totalPages++
	}

	return &models.PaginatedAgenciesResponse{
		Agencies:   agencyResponses,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

// 2. GetAgency 获取代理公司详情
func (s *AgencyService) GetAgency(ctx context.Context, id uint) (*models.AgencyDetailResponse, error) {
	agency, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// 获取代理公司的代理人列表
	agents, err := s.repo.GetAgencyAgents(ctx, agency.UserID)
	if err != nil {
		return nil, err
	}

	// 更新代理人数量（如果不一致）
	if len(agents) != agency.AgentCount {
		agency.AgentCount = len(agents)
	}

	return &models.AgencyDetailResponse{
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
		CreatedAt:              agency.CreatedAt,
		UpdatedAt:              agency.UpdatedAt,
	}, nil
}

// 3. GetAgencyProperties 获取代理公司房源列表
func (s *AgencyService) GetAgencyProperties(ctx context.Context, id uint, page, pageSize int) ([]models.Property, int64, error) {
	// 验证代理公司存在
	agency, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			return nil, 0, tools.ErrNotFound
		}
		return nil, 0, err
	}

	// 获取房源列表
	properties, total, err := s.repo.GetAgencyProperties(ctx, agency.UserID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return properties, total, nil
}

// 4. ContactAgency 联系代理公司
func (s *AgencyService) ContactAgency(ctx context.Context, agencyID uint, userID *uint, req *models.ContactAgencyRequest) (*models.ContactAgencyResponse, error) {
	// 验证代理公司存在
	agency, err := s.repo.FindByID(ctx, agencyID)
	if err != nil {
		if errors.Is(err, tools.ErrNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// 验证代理公司的用户状态
	if agency.User != nil && agency.User.Status != "active" {
		return nil, errors.New("agency is not active")
	}

	// 创建联系记录
	contact := &models.AgencyContact{
		AgencyID:   agency.UserID,
		UserID:     userID,
		Name:       req.Name,
		Phone:      req.Phone,
		Email:      req.Email,
		Message:    req.Message,
		PropertyID: req.PropertyID,
	}

	if err := s.repo.CreateContact(ctx, contact); err != nil {
		return nil, err
	}

	return &models.ContactAgencyResponse{
		ID:        contact.ID,
		AgencyID:  contact.AgencyID,
		Message:   "Contact request submitted successfully",
		CreatedAt: contact.CreatedAt,
	}, nil
}

// 5. SearchAgencies 搜索代理公司
func (s *AgencyService) SearchAgencies(ctx context.Context, req *models.SearchAgenciesRequest) (*models.PaginatedAgenciesResponse, error) {
	agencies, total, err := s.repo.Search(ctx, req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	var agencyResponses []models.AgencyResponse
	for _, agency := range agencies {
		agencyResponses = append(agencyResponses, models.AgencyResponse{
			ID:              agency.ID,
			CompanyName:     agency.CompanyName,
			CompanyNameEn:   agency.CompanyNameEn,
			LogoURL:         agency.LogoURL,
			Address:         agency.Address,
			Phone:           agency.Phone,
			Email:           agency.Email,
			WebsiteURL:      agency.WebsiteURL,
			EstablishedYear: agency.EstablishedYear,
			AgentCount:      agency.AgentCount,
			Rating:          agency.Rating,
			ReviewCount:     agency.ReviewCount,
			IsVerified:      agency.IsVerified,
		})
	}

	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize > 0 {
		totalPages++
	}

	return &models.PaginatedAgenciesResponse{
		Agencies:   agencyResponses,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}
