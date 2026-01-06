package services

// UserService Methods:
// 0. NewUserService(userRepo databases.UserRepository, propertyRepo databases.PropertyRepository) -> 注入依赖
// 1. GetCurrentUser(ctx context.Context, userID uint) -> 获取当前用户信息
// 2. UpdateCurrentUser(ctx context.Context, userID uint, req *models.User) -> 更新当前用户信息
// 3. ChangePassword(ctx context.Context, userID uint, req *models.ChangePasswordRequest) -> 修改密码
// 4. GetMyListings(ctx context.Context, userID uint, page, pageSize int) -> 获取我的发布
// 5. UpdateSettings(ctx context.Context, userID uint, req *map[string]interface{}) -> 更新设置

import (
	"context"
	"errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
)

var (
	ErrUserNotActive      = errors.New("user is not active")
	ErrOldPasswordInvalid = errors.New("old password is invalid")
)

// UserServiceInterface 用户服务接口
type UserServiceInterface interface {
	GetCurrentUser(ctx context.Context, userID uint) (*models.User, error)
	UpdateCurrentUser(ctx context.Context, userID uint, req *models.User) (*models.User, error)
	ChangePassword(ctx context.Context, userID uint, req *models.ChangePasswordRequest) error
	GetMyListings(ctx context.Context, userID uint, page, pageSize int) ([]*models.Property, error)
	UpdateSettings(ctx context.Context, userID uint, req *map[string]interface{}) (*map[string]interface{}, error)
}

// UserService 用户服务
type UserService struct {
	userRepo     databases.UserRepository
	propertyRepo databases.PropertyRepository
}

// 0. NewUserService 注入依赖
func NewUserService(userRepo databases.UserRepository, propertyRepo databases.PropertyRepository) *UserService {
	return &UserService{
		userRepo:     userRepo,
		propertyRepo: propertyRepo,
	}
}

// 1. GetCurrentUser 获取当前用户信息
func (s *UserService) GetCurrentUser(ctx context.Context, userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, tools.ErrNotFound
	}

	return &models.User{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		FullName:      user.FullName,
		Phone:         user.Phone,
		Avatar:        user.Avatar,
		UserType:      user.UserType,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		LastLoginAt:   user.LastLoginAt,
	}, nil
}

// 2. UpdateCurrentUser 更新当前用户信息
func (s *UserService) UpdateCurrentUser(ctx context.Context, userID uint, req *models.User) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, tools.ErrNotFound
	}

	// 更新字段
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &models.User{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		FullName:      user.FullName,
		Phone:         user.Phone,
		Avatar:        user.Avatar,
		UserType:      user.UserType,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		LastLoginAt:   user.LastLoginAt,
	}, nil
}

// 3. ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, userID uint, req *models.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return tools.ErrNotFound
	}

	// 验证旧密码
	if !tools.CheckPassword(user.Password, req.OldPassword) {
		return ErrOldPasswordInvalid
	}

	// 加密新密码
	hashedPassword, err := tools.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(ctx, user)
}

// 4. GetMyListings 获取我的发布
func (s *UserService) GetMyListings(ctx context.Context, userID uint, page, pageSize int) ([]*models.Property, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	properties, _, err := s.propertyRepo.ListByPublisher(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return properties, nil
}

// 5. UpdateSettings 更新设置
func (s *UserService) UpdateSettings(ctx context.Context, userID uint, req *map[string]interface{}) (*map[string]interface{}, error) {
	// TODO: 实现用户设置存储（可能需要单独的 user_settings 表）
	// 目前直接返回请求中的设置作为响应
	return req, nil
}
