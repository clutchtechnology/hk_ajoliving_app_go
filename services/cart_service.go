package services

import (
	"context"
	"errors"

	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
	"gorm.io/gorm"
)

// CartService Methods:
// 0. NewCartService(cartRepo *databases.CartRepo, furnitureRepo *databases.FurnitureRepo) -> 注入依赖
// 1. GetCart(ctx context.Context, userID uint) -> 获取购物车
// 2. AddToCart(ctx context.Context, userID uint, req *models.AddToCartRequest) -> 添加到购物车
// 3. UpdateCartItem(ctx context.Context, userID, itemID uint, req *models.UpdateCartItemRequest) -> 更新购物车项
// 4. RemoveFromCart(ctx context.Context, userID, itemID uint) -> 移除购物车项
// 5. ClearCart(ctx context.Context, userID uint) -> 清空购物车

type CartService struct {
	cartRepo      *databases.CartRepo
	furnitureRepo *databases.FurnitureRepo
}

// 0. NewCartService 构造函数
func NewCartService(cartRepo *databases.CartRepo, furnitureRepo *databases.FurnitureRepo) *CartService {
	return &CartService{
		cartRepo:      cartRepo,
		furnitureRepo: furnitureRepo,
	}
}

// 1. GetCart 获取购物车
func (s *CartService) GetCart(ctx context.Context, userID uint) (*models.CartResponse, error) {
	items, err := s.cartRepo.GetUserCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	var cartItems []*models.CartItemResponse
	var totalPrice float64
	totalItems := 0

	for _, item := range items {
		if item.Furniture == nil {
			continue
		}

		// 构建购物车项响应
		furnitureInCart := &models.FurnitureInCart{
			ID:          item.Furniture.ID,
			FurnitureNo: item.Furniture.FurnitureNo,
			Title:       item.Furniture.Title,
			Price:       item.Furniture.Price,
			Status:      item.Furniture.Status,
		}

		// 添加分类名称
		if item.Furniture.Category != nil {
			furnitureInCart.CategoryName = item.Furniture.Category.NameZhHant
		}

		// 添加地区名称
		if item.Furniture.DeliveryDistrict != nil {
			furnitureInCart.DistrictName = item.Furniture.DeliveryDistrict.NameZhHant
		}

		// 添加封面图
		if len(item.Furniture.Images) > 0 {
			furnitureInCart.CoverImageURL = item.Furniture.Images[0].ImageURL
		}

		cartItem := &models.CartItemResponse{
			ID:          item.ID,
			FurnitureID: item.FurnitureID,
			Quantity:    item.Quantity,
			Furniture:   furnitureInCart,
			CreatedAt:   item.CreatedAt,
		}

		cartItems = append(cartItems, cartItem)
		totalPrice += item.Furniture.Price * float64(item.Quantity)
		totalItems += item.Quantity
	}

	return &models.CartResponse{
		Items:      cartItems,
		TotalItems: totalItems,
		TotalPrice: totalPrice,
	}, nil
}

// 2. AddToCart 添加到购物车
func (s *CartService) AddToCart(ctx context.Context, userID uint, req *models.AddToCartRequest) (*models.CartItemResponse, error) {
	// 验证家具是否存在且可用
	furniture, err := s.furnitureRepo.FindByID(ctx, req.FurnitureID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// 检查家具状态
	if furniture.Status != "available" {
		return nil, errors.New("furniture is not available for purchase")
	}

	// 检查是否已在购物车中
	existingItem, err := s.cartRepo.FindByUserAndFurniture(ctx, userID, req.FurnitureID)
	if err == nil && existingItem != nil {
		// 更新数量
		existingItem.Quantity += req.Quantity
		if err := s.cartRepo.Update(ctx, existingItem); err != nil {
			return nil, err
		}
		
		// 重新加载关联数据
		existingItem, _ = s.cartRepo.FindByID(ctx, existingItem.ID)
		return s.buildCartItemResponse(existingItem), nil
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

	// 重新加载关联数据
	cartItem, _ = s.cartRepo.FindByID(ctx, cartItem.ID)
	return s.buildCartItemResponse(cartItem), nil
}

// 3. UpdateCartItem 更新购物车项
func (s *CartService) UpdateCartItem(ctx context.Context, userID, itemID uint, req *models.UpdateCartItemRequest) (*models.CartItemResponse, error) {
	// 查找购物车项
	item, err := s.cartRepo.FindByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tools.ErrNotFound
		}
		return nil, err
	}

	// 验证权限
	if item.UserID != userID {
		return nil, tools.ErrForbidden
	}

	// 更新数量
	item.Quantity = req.Quantity
	if err := s.cartRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	// 重新加载关联数据
	item, _ = s.cartRepo.FindByID(ctx, item.ID)
	return s.buildCartItemResponse(item), nil
}

// 4. RemoveFromCart 移除购物车项
func (s *CartService) RemoveFromCart(ctx context.Context, userID, itemID uint) error {
	// 查找购物车项
	item, err := s.cartRepo.FindByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return tools.ErrNotFound
		}
		return err
	}

	// 验证权限
	if item.UserID != userID {
		return tools.ErrForbidden
	}

	return s.cartRepo.Delete(ctx, itemID)
}

// 5. ClearCart 清空购物车
func (s *CartService) ClearCart(ctx context.Context, userID uint) error {
	return s.cartRepo.ClearUserCart(ctx, userID)
}

// buildCartItemResponse 构建购物车项响应
func (s *CartService) buildCartItemResponse(item *models.CartItem) *models.CartItemResponse {
	if item.Furniture == nil {
		return &models.CartItemResponse{
			ID:          item.ID,
			FurnitureID: item.FurnitureID,
			Quantity:    item.Quantity,
			CreatedAt:   item.CreatedAt,
		}
	}

	furnitureInCart := &models.FurnitureInCart{
		ID:          item.Furniture.ID,
		FurnitureNo: item.Furniture.FurnitureNo,
		Title:       item.Furniture.Title,
		Price:       item.Furniture.Price,
		Status:      item.Furniture.Status,
	}

	if item.Furniture.Category != nil {
		furnitureInCart.CategoryName = item.Furniture.Category.NameZhHant
	}

	if item.Furniture.DeliveryDistrict != nil {
		furnitureInCart.DistrictName = item.Furniture.DeliveryDistrict.NameZhHant
	}

	if len(item.Furniture.Images) > 0 {
		furnitureInCart.CoverImageURL = item.Furniture.Images[0].ImageURL
	}

	return &models.CartItemResponse{
		ID:          item.ID,
		FurnitureID: item.FurnitureID,
		Quantity:    item.Quantity,
		Furniture:   furnitureInCart,
		CreatedAt:   item.CreatedAt,
	}
}
