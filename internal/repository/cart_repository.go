package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
)

// CartRepository 购物车数据访问接口
type CartRepository interface {
	GetByUserID(ctx context.Context, userID uint) ([]*model.CartItem, error)
	GetByID(ctx context.Context, id uint) (*model.CartItem, error)
	GetByUserAndFurniture(ctx context.Context, userID uint, furnitureID uint) (*model.CartItem, error)
	Create(ctx context.Context, cartItem *model.CartItem) error
	Update(ctx context.Context, cartItem *model.CartItem) error
	Delete(ctx context.Context, id uint) error
	DeleteByUserID(ctx context.Context, userID uint) error
	CountByUserID(ctx context.Context, userID uint) (int64, error)
}

type cartRepository struct {
	db *gorm.DB
}

// NewCartRepository 创建购物车仓库
func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetByUserID(ctx context.Context, userID uint) ([]*model.CartItem, error) {
	var cartItems []*model.CartItem
	
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Furniture").
		Preload("Furniture.Category").
		Preload("Furniture.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_cover = ?", true).Order("sort_order ASC").Limit(1)
		}).
		Order("created_at DESC").
		Find(&cartItems).Error
		
	if err != nil {
		return nil, err
	}
	
	return cartItems, nil
}

func (r *cartRepository) GetByID(ctx context.Context, id uint) (*model.CartItem, error) {
	var cartItem model.CartItem
	
	err := r.db.WithContext(ctx).
		Preload("Furniture").
		Preload("Furniture.Category").
		Preload("Furniture.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_cover = ?", true).Order("sort_order ASC").Limit(1)
		}).
		First(&cartItem, id).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &cartItem, nil
}

func (r *cartRepository) GetByUserAndFurniture(ctx context.Context, userID uint, furnitureID uint) (*model.CartItem, error) {
	var cartItem model.CartItem
	
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND furniture_id = ?", userID, furnitureID).
		Preload("Furniture").
		First(&cartItem).Error
		
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	
	return &cartItem, nil
}

func (r *cartRepository) Create(ctx context.Context, cartItem *model.CartItem) error {
	return r.db.WithContext(ctx).Create(cartItem).Error
}

func (r *cartRepository) Update(ctx context.Context, cartItem *model.CartItem) error {
	return r.db.WithContext(ctx).Save(cartItem).Error
}

func (r *cartRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.CartItem{}, id).Error
}

func (r *cartRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.CartItem{}).Error
}

func (r *cartRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.CartItem{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
