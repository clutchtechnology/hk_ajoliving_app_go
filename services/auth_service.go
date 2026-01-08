package services

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
)

// AuthService 认证服务
type AuthService struct {
	userRepo *databases.UserRepo
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo *databases.UserRepo) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

// Register 用户注册
func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	// 检查邮箱是否已存在
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// 密码加密
	hashedPassword, err := tools.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &models.User{
		UserType:     req.UserType,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
		Phone:        req.Phone,
		Status:       "active",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	// 查找用户
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	// 验证密码
	if !tools.CheckPassword(user.PasswordHash, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.New("account is not active")
	}

	// 生成 JWT Token
	token, err := tools.GenerateToken(user.ID, user.Email, user.UserType)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	return &models.LoginResponse{
		Token: token,
		User:  user.ToUserResponse(),
	}, nil
}
