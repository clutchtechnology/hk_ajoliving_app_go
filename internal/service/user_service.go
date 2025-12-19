package service

// UserService Methods:
// 0. NewUserService(userRepo repository.UserRepository, propertyRepo repository.PropertyRepository) -> 注入依赖
// 1. GetCurrentUser(ctx context.Context, userID uint) -> 获取当前用户信息
// 2. UpdateCurrentUser(ctx context.Context, userID uint, req *request.UpdateUserRequest) -> 更新当前用户信息
// 3. ChangePassword(ctx context.Context, userID uint, req *request.ChangePasswordRequest) -> 修改密码
// 4. GetMyListings(ctx context.Context, userID uint, page, pageSize int) -> 获取我的发布
// 5. UpdateSettings(ctx context.Context, userID uint, req *request.UpdateSettingsRequest) -> 更新设置

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/utils"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
)

var (
	ErrUserNotActive      = errors.New("user is not active")
	ErrOldPasswordInvalid = errors.New("old password is invalid")
)

// UserServiceInterface 用户服务接口
type UserServiceInterface interface {
	GetCurrentUser(ctx context.Context, userID uint) (*response.UserResponse, error)
	UpdateCurrentUser(ctx context.Context, userID uint, req *request.UpdateUserRequest) (*response.UserResponse, error)
	ChangePassword(ctx context.Context, userID uint, req *request.ChangePasswordRequest) error
	GetMyListings(ctx context.Context, userID uint, page, pageSize int) (*response.MyListingsResponse, error)
	UpdateSettings(ctx context.Context, userID uint, req *request.UpdateSettingsRequest) (*response.UserSettingsResponse, error)
}

// UserService 用户服务
type UserService struct {
	userRepo     repository.UserRepository
	propertyRepo repository.PropertyRepository
}

// 0. NewUserService 注入依赖
func NewUserService(userRepo repository.UserRepository, propertyRepo repository.PropertyRepository) *UserService {
	return &UserService{
		userRepo:     userRepo,
		propertyRepo: propertyRepo,
	}
}

// 1. GetCurrentUser 获取当前用户信息
func (s *UserService) GetCurrentUser(ctx context.Context, userID uint) (*response.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &response.UserResponse{
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
func (s *UserService) UpdateCurrentUser(ctx context.Context, userID uint, req *request.UpdateUserRequest) (*response.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
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

	return &response.UserResponse{
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
func (s *UserService) ChangePassword(ctx context.Context, userID uint, req *request.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// 验证旧密码
	if !utils.CheckPassword(user.Password, req.OldPassword) {
		return ErrOldPasswordInvalid
	}

	// 加密新密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.userRepo.Update(ctx, user)
}

// 4. GetMyListings 获取我的发布
func (s *UserService) GetMyListings(ctx context.Context, userID uint, page, pageSize int) (*response.MyListingsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	properties, total, err := s.propertyRepo.ListByPublisher(ctx, userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	items := make([]response.PropertyBriefResponse, 0, len(properties))
	for _, p := range properties {
		coverImage := ""
		for _, img := range p.Images {
			if img.IsCover {
				coverImage = img.URL
				break
			}
		}
		if coverImage == "" && len(p.Images) > 0 {
			coverImage = p.Images[0].URL
		}

		items = append(items, response.PropertyBriefResponse{
			ID:            p.ID,
			PropertyNo:    p.PropertyNo,
			Title:         p.Title,
			Price:         p.Price,
			Area:          p.Area,
			Bedrooms:      p.Bedrooms,
			ListingType:   string(p.ListingType),
			Status:        string(p.Status),
			CoverImage:    coverImage,
			ViewCount:     p.ViewCount,
			FavoriteCount: p.FavoriteCount,
			CreatedAt:     p.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &response.MyListingsResponse{
		Properties: items,
		Total:      total,
	}, nil
}

// 5. UpdateSettings 更新设置
func (s *UserService) UpdateSettings(ctx context.Context, userID uint, req *request.UpdateSettingsRequest) (*response.UserSettingsResponse, error) {
	// TODO: 实现用户设置存储（可能需要单独的 user_settings 表）
	// 目前返回请求中的设置作为响应
	return &response.UserSettingsResponse{
		Language:            req.Language,
		NotificationEnabled: req.NotificationEnabled,
		EmailNotification:   req.EmailNotification,
		SmsNotification:     req.SmsNotification,
	}, nil
}
