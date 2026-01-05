package services

// AuthService Methods:
// 0. NewAuthService(userRepo databases.UserRepository, jwtManager *tools.JWTManager) -> 注入依赖
// 1. Register(ctx context.Context, req *models.RegisterRequest) -> 用户注册
// 2. Login(ctx context.Context, req *models.LoginRequest) -> 用户登录
// 3. Logout(ctx context.Context, userID uint) -> 用户登出
// 4. RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) -> 刷新令牌
// 5. ForgotPassword(ctx context.Context, req *models.ForgotPasswordRequest) -> 忘记密码
// 6. ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) -> 重置密码
// 7. VerifyCode(ctx context.Context, req *models.VerifyCodeRequest) -> 验证码验证

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
)

var (
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrUserNotFound            = errors.New("user not found")
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidVerificationCode = errors.New("invalid verification code")
)

// AuthServiceInterface 认证服务接口
type AuthServiceInterface interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error)       // 1. 用户注册
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)                 // 2. 用户登录
	Logout(ctx context.Context, userID uint) error                                                        // 3. 用户登出
	RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.AuthResponse, error)   // 4. 刷新令牌
	ForgotPassword(ctx context.Context, req *models.ForgotPasswordRequest) (*models.ForgotPasswordResponse, error) // 5. 忘记密码
	ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) (*models.ResetPasswordResponse, error)    // 6. 重置密码
	VerifyCode(ctx context.Context, req *models.VerifyCodeRequest) (*models.VerifyCodeResponse, error) // 7. 验证码验证
}

// AuthService 认证服务
type AuthService struct {
	userRepo   databases.UserRepository
	jwtManager *tools.JWTManager
}

// 0. NewAuthService 注入依赖
func NewAuthService(userRepo databases.UserRepository, jwtManager *tools.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// 1. Register 用户注册
func (s *AuthService) Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error) {
	// 检查用户是否已存在
	exists, err := s.userRepo.Exists(ctx, req.Email, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUserAlreadyExists
	}

	// 加密密码
	hashedPassword, err := tools.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 设置默认用户类型
	userType := req.UserType
	if userType == "" {
		userType = "buyer"
	}

	// 创建用户
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Phone:    req.Phone,
		FullName: req.FullName,
		UserType: userType,
		Status:   "active",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &models.RegisterResponse{
		User: &models.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			Phone:     user.Phone,
			Avatar:    user.Avatar,
			UserType:  user.UserType,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// 2. Login 用户登录
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	// 查找用户
	user, err := s.userRepo.GetByEmailOrUsername(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// 验证密码
	if !tools.CheckPassword(user.Password, req.Password) {
		return nil, ErrInvalidCredentials
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, errors.New("user account is not active")
	}

	// 生成令牌
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username, user.Email, user.UserType)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		// 记录日志但不返回错误
		fmt.Printf("failed to update last login time: %v\n", err)
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.jwtManager.GetAccessExpireSeconds(),
		User: &models.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			Phone:     user.Phone,
			Avatar:    user.Avatar,
			UserType:  user.UserType,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// 3. Logout 用户登出
func (s *AuthService) Logout(ctx context.Context, userID uint) error {
	// 实际项目中可以将令牌加入黑名单（Redis）
	// 这里简单返回成功
	return nil
}

// 4. RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(ctx context.Context, req *models.RefreshTokenRequest) (*models.AuthResponse, error) {
	// 解析刷新令牌
	claims, err := s.jwtManager.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 获取用户信息
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 生成新的访问令牌
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username, user.Email, user.UserType)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    s.jwtManager.GetAccessExpireSeconds(),
		User: &models.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			Phone:     user.Phone,
			Avatar:    user.Avatar,
			UserType:  user.UserType,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// 5. ForgotPassword 忘记密码
func (s *AuthService) ForgotPassword(ctx context.Context, req *models.ForgotPasswordRequest) (*models.ForgotPasswordResponse, error) {
	// 查找用户
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		// 为了安全，即使用户不存在也返回成功消息
		return &models.ForgotPasswordResponse{
			Email:   req.Email,
		}, nil
	}

	// TODO: 生成重置令牌并发送邮件
	// 实际项目中应该:
	// 1. 生成随机令牌
	// 2. 将令牌存储到 Redis 并设置过期时间（如 1 小时）
	// 3. 发送包含重置链接的邮件

	return &models.ForgotPasswordResponse{
		Email:   req.Email,
	}, nil
}

// 6. ResetPassword 重置密码
func (s *AuthService) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) (*models.ResetPasswordResponse, error) {
	// TODO: 验证重置令牌（从 Redis 中获取）
	// 这里简化处理，实际项目中应该:
	// 1. 从 Redis 获取令牌对应的用户 ID
	// 2. 验证令牌是否过期
	// 3. 删除已使用的令牌

	// 解析令牌获取用户信息（简化版）
	claims, err := s.jwtManager.ParseToken(req.Token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 获取用户
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// 加密新密码
	hashedPassword, err := tools.HashPassword(req.NewPassword)
	if err != nil {
		return nil, err
	}

	// 更新密码
	user.Password = hashedPassword
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &models.ResetPasswordResponse{
		}, nil
}

// 7. VerifyCode 验证码验证
func (s *AuthService) VerifyCode(ctx context.Context, req *models.VerifyCodeRequest) (*models.VerifyCodeResponse, error) {
	// TODO: 从 Redis 中验证验证码
	// 实际项目中应该:
	// 1. 从 Redis 获取该邮箱对应的验证码
	// 2. 比对验证码是否匹配
	// 3. 检查验证码是否过期
	// 4. 验证成功后删除验证码

	// 简化处理：假设验证码正确
	isValid := req.Code == "123456" // 仅用于演示

	if !isValid {
		return &models.VerifyCodeResponse{
			Valid:   false,
			}, nil
	}

	// 生成临时令牌用于后续操作
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	token, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username, user.Email, user.UserType)
	if err != nil {
		return nil, err
	}

	return &models.VerifyCodeResponse{
		Valid:   true,
		Token:   token,
	}, nil
}
