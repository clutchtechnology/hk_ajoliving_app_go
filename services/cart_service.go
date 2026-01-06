package services

// CartService Methods:
// 0. NewCartService(cartRepo databases.CartRepository, furnitureRepo databases.FurnitureRepository) -> 注入依赖
// 1. GetCart(ctx context.Context, userID uint) -> 获取购物车
// 2. AddToCart(ctx context.Context, userID uint, req *models.AddToCartRequest) -> 添加到购物车
// 3. UpdateCartItem(ctx context.Context, userID uint, id uint, req *models.UpdateCartItemRequest) -> 更新购物车项
// 4. RemoveFromCart(ctx context.Context, userID uint, id uint) -> 移除购物车项
// 5. ClearCart(ctx context.Context, userID uint) -> 清空购物车

import (
	"context"
	"errors"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
)

var (
	ErrCartItemNotFound       = errors.New("cart item not found")
	ErrNotCartItemOwner       = errors.New("you are not the owner of this cart item")
	ErrFurnitureNotFound      = errors.New("furniture not found")
	ErrFurnitureNotAvailable  = errors.New("furniture is not available")
	ErrFurnitureAlreadyInCart = errors.New("furniture is already in cart")
)

// CartServiceInterface 购物车服务接口
type CartServiceInterface interface {
	GetCart(ctx context.Context, userID uint) (*models.CartResponse, error)
	AddToCart(ctx context.Context, userID uint, req *models.AddToCartRequest) (*models.AddToCartResponse, error)
	UpdateCartItem(ctx context.Context, userID uint, id uint, req *models.UpdateCartItemRequest) (*models.CartResponse, error)
	RemoveFromCart(ctx context.Context, userID uint, id uint) error
	ClearCart(ctx context.Context, userID uint) error
}

// CartService 购物车服务
type CartService struct {
	cartRepo      databases.CartRepository
	furnitureRepo databases.FurnitureRepository
}

// 0. NewCartService 注入依赖
func NewCartService(cartRepo databases.CartRepository, furnitureRepo databases.FurnitureRepository) *CartService {
	return &CartService{
		cartRepo:      cartRepo,
		furnitureRepo: furnitureRepo,
	}
}

// 1. GetCart 获取购物车
func (s *CartService) GetCart(ctx context.Context, userID uint) (*models.CartResponse, error) {
	cartItems, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.convertToCartResponse(cartItems), nil
}

// 2. AddToCart 添加到购物车
func (s *CartService) AddToCart(ctx context.Context, userID uint, req *models.AddToCartRequest) (*models.AddToCartResponse, error) {
	// 验证家具是否存在
	furniture, err := s.furnitureRepo.GetByID(ctx, req.FurnitureID)
	if err != nil {
		return nil, err
	}
	if furniture == nil {
		return nil, tools.ErrNotFound
	}

	// 验证家具是否可用
	if !furniture.IsAvailable() {
		return nil, ErrFurnitureNotAvailable
	}

	// 检查是否已在购物车中
	existingItem, err := s.cartRepo.GetByUserAndFurniture(ctx, userID, req.FurnitureID)
	if err != nil {
		return nil, err
	}
	if existingItem != nil {
		// 如果已存在，更新数量
		existingItem.Quantity += req.Quantity
		if err := s.cartRepo.Update(ctx, existingItem); err != nil {
			return nil, err
		}

		return &models.AddToCartResponse{
			Success: true,
			Message: "Item quantity updated in cart",
			Item:    existingItem,
		}, nil
	}

	// 创建新的购物车项
	cartItem := &models.CartItem{
		UserID:      userID,
		FurnitureID: req.FurnitureID,
		Quantity:    req.Quantity,
	}

	if err := s.cartRepo.Create(ctx, cartItem); err != nil {
		return nil, err
	}

	return &models.AddToCartResponse{
		Success: true,
		Message: "Item added to cart successfully",
		Item:    cartItem,
	}, nil
}

// 3. UpdateCartItem 更新购物车项
func (s *CartService) UpdateCartItem(ctx context.Context, userID uint, id uint, req *models.UpdateCartItemRequest) (*models.CartResponse, error) {
	// 获取购物车项
	cartItem, err := s.cartRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cartItem == nil {
		return nil, ErrCartItemNotFound
	}

	// 验证权限
	if cartItem.UserID != userID {
		return nil, ErrNotCartItemOwner
	}

	// 更新数量
	cartItem.Quantity = req.Quantity

	if err := s.cartRepo.Update(ctx, cartItem); err != nil {
		return nil, err
	}

	// 返回完整的购物车
	return s.GetCart(ctx, userID)
}

// 4. RemoveFromCart 移除购物车项
func (s *CartService) RemoveFromCart(ctx context.Context, userID uint, id uint) error {
	// 获取购物车项
	cartItem, err := s.cartRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if cartItem == nil {
		return ErrCartItemNotFound
	}

	// 验证权限
	if cartItem.UserID != userID {
		return ErrNotCartItemOwner
	}

	// 删除购物车项
	return s.cartRepo.Delete(ctx, id)
}

// 5. ClearCart 清空购物车
func (s *CartService) ClearCart(ctx context.Context, userID uint) error {
	return s.cartRepo.DeleteByUserID(ctx, userID)
}

// 辅助方法

// convertToCartResponse 转换为购物车响应
func (s *CartService) convertToCartResponse(cartItems []*models.CartItem) *models.CartResponse {
	items := make([]models.CartItem, 0, len(cartItems))
	totalPrice := 0.0

	for _, item := range cartItems {
		items = append(items, *item)
		totalPrice += item.GetTotalPrice()
	}

	return &models.CartResponse{
		Items:      items,
		TotalItems: len(items),
		TotalPrice: totalPrice,
	}
}

