package services

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
)

// UserService 用户服务
type UserService struct {
	userRepo     *databases.UserRepo
	propertyRepo *databases.PropertyRepo
}

// NewUserService 创建用户服务
func NewUserService(userRepo *databases.UserRepo, propertyRepo *databases.PropertyRepo) *UserService {
	return &UserService{
		userRepo:     userRepo,
		propertyRepo: propertyRepo,
	}
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(ctx context.Context, id uint, req *models.UpdateUserRequest) (*models.User, error) {
	// 查找用户
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}

	// 保存更新
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserListings 获取用户的发布（房源、家具等）
func (s *UserService) GetUserListings(ctx context.Context, userID uint, listingType *string) (interface{}, error) {
	// 获取用户房产列表
	properties, err := s.propertyRepo.FindByPublisher(ctx, userID, listingType)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	propertyResponses := make([]models.PropertyResponse, len(properties))
	for i, p := range properties {
		propertyResponses[i] = *p.ToPropertyResponse()
	}

	return map[string]interface{}{
		"properties": propertyResponses,
		"furniture":  []interface{}{}, // 家具模块后续实现
	}, nil
}
