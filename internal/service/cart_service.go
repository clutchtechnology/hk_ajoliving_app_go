package service

// CartService Methods:
// 0. NewCartService(cartRepo repository.CartRepository, furnitureRepo repository.FurnitureRepository) -> 注入依赖
// 1. GetCart(ctx context.Context, userID uint) -> 获取购物车
// 2. AddToCart(ctx context.Context, userID uint, req *request.AddToCartRequest) -> 添加到购物车
// 3. UpdateCartItem(ctx context.Context, userID uint, id uint, req *request.UpdateCartItemRequest) -> 更新购物车项
// 4. RemoveFromCart(ctx context.Context, userID uint, id uint) -> 移除购物车项
// 5. ClearCart(ctx context.Context, userID uint) -> 清空购物车

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/request"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
)

var (
	ErrCartItemNotFound      = errors.New("cart item not found")
	ErrNotCartItemOwner      = errors.New("you are not the owner of this cart item")
	ErrFurnitureNotAvailable = errors.New("furniture is not available")
	ErrFurnitureAlreadyInCart = errors.New("furniture is already in cart")
)

// CartServiceInterface 购物车服务接口
type CartServiceInterface interface {
	GetCart(ctx context.Context, userID uint) (*response.CartResponse, error)
	AddToCart(ctx context.Context, userID uint, req *request.AddToCartRequest) (*response.AddToCartResponse, error)
	UpdateCartItem(ctx context.Context, userID uint, id uint, req *request.UpdateCartItemRequest) (*response.UpdateCartItemResponse, error)
	RemoveFromCart(ctx context.Context, userID uint, id uint) error
	ClearCart(ctx context.Context, userID uint) error
}

// CartService 购物车服务
type CartService struct {
	cartRepo      repository.CartRepository
	furnitureRepo repository.FurnitureRepository
}

// 0. NewCartService 注入依赖
func NewCartService(cartRepo repository.CartRepository, furnitureRepo repository.FurnitureRepository) *CartService {
	return &CartService{
		cartRepo:      cartRepo,
		furnitureRepo: furnitureRepo,
	}
}

// 1. GetCart 获取购物车
func (s *CartService) GetCart(ctx context.Context, userID uint) (*response.CartResponse, error) {
	cartItems, err := s.cartRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.convertToCartResponse(cartItems), nil
}

// 2. AddToCart 添加到购物车
func (s *CartService) AddToCart(ctx context.Context, userID uint, req *request.AddToCartRequest) (*response.AddToCartResponse, error) {
	// 验证家具是否存在
	furniture, err := s.furnitureRepo.GetByID(ctx, req.FurnitureID)
	if err != nil {
		return nil, err
	}
	if furniture == nil {
		return nil, ErrFurnitureNotFound
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

		totalPrice := furniture.Price * float64(existingItem.Quantity)
		return &response.AddToCartResponse{
			ID:          existingItem.ID,
			FurnitureID: existingItem.FurnitureID,
			Quantity:    existingItem.Quantity,
			TotalPrice:  totalPrice,
			CreatedAt:   existingItem.CreatedAt,
			Message:     "已更新购物车数量",
		}, nil
	}

	// 创建新的购物车项
	cartItem := &model.CartItem{
		UserID:      userID,
		FurnitureID: req.FurnitureID,
		Quantity:    req.Quantity,
	}

	if err := s.cartRepo.Create(ctx, cartItem); err != nil {
		return nil, err
	}

	totalPrice := furniture.Price * float64(cartItem.Quantity)
	return &response.AddToCartResponse{
		ID:          cartItem.ID,
		FurnitureID: cartItem.FurnitureID,
		Quantity:    cartItem.Quantity,
		TotalPrice:  totalPrice,
		CreatedAt:   cartItem.CreatedAt,
		Message:     "添加到购物车成功",
	}, nil
}

// 3. UpdateCartItem 更新购物车项
func (s *CartService) UpdateCartItem(ctx context.Context, userID uint, id uint, req *request.UpdateCartItemRequest) (*response.UpdateCartItemResponse, error) {
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

	totalPrice := float64(0)
	if cartItem.Furniture != nil {
		totalPrice = cartItem.Furniture.Price * float64(cartItem.Quantity)
	}

	return &response.UpdateCartItemResponse{
		ID:         cartItem.ID,
		Quantity:   cartItem.Quantity,
		TotalPrice: totalPrice,
		UpdatedAt:  cartItem.UpdatedAt,
		Message:    "购物车项更新成功",
	}, nil
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
func (s *CartService) convertToCartResponse(cartItems []*model.CartItem) *response.CartResponse {
	items := make([]response.CartItemResponse, 0, len(cartItems))
	totalQuantity := 0
	totalPrice := 0.0
	availableItems := 0
	unavailableItems := 0

	for _, item := range cartItems {
		itemResp := s.convertToCartItemResponse(item)
		items = append(items, itemResp)
		
		totalQuantity += item.Quantity
		totalPrice += itemResp.TotalPrice
		
		if itemResp.IsAvailable {
			availableItems++
		} else {
			unavailableItems++
		}
	}

	return &response.CartResponse{
		Items:            items,
		TotalItems:       len(items),
		TotalQuantity:    totalQuantity,
		TotalPrice:       totalPrice,
		AvailableItems:   availableItems,
		UnavailableItems: unavailableItems,
	}
}

// convertToCartItemResponse 转换为购物车项响应
func (s *CartService) convertToCartItemResponse(item *model.CartItem) response.CartItemResponse {
	resp := response.CartItemResponse{
		ID:          item.ID,
		FurnitureID: item.FurnitureID,
		Quantity:    item.Quantity,
		TotalPrice:  item.GetTotalPrice(),
		IsAvailable: item.IsAvailable(),
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}

	// 家具信息
	if item.Furniture != nil {
		furnitureResp := &response.CartFurnitureResponse{
			ID:          item.Furniture.ID,
			FurnitureNo: item.Furniture.FurnitureNo,
			Title:       item.Furniture.Title,
			Price:       item.Furniture.Price,
			Status:      string(item.Furniture.Status),
			IsAvailable: item.Furniture.IsAvailable(),
		}

		// 封面图片
		if len(item.Furniture.Images) > 0 {
			furnitureResp.CoverImage = &item.Furniture.Images[0].ImageURL
		}

		resp.Furniture = furnitureResp
	}

	return resp
}
