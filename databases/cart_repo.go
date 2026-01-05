package databases

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
)

// CartRepository 购物车数据访问接口
type CartRepository interface {
	GetByUserID(ctx context.Context, userID uint) ([]*models.CartItem, error)
	GetByID(ctx context.Context, id uint) (*models.CartItem, error)
	GetByUserAndFurniture(ctx context.Context, userID uint, furnitureID uint) (*models.CartItem, error)
	Create(ctx context.Context, cartItem *models.CartItem) error
	Update(ctx context.Context, cartItem *models.CartItem) error
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

func (r *cartRepository) GetByUserID(ctx context.Context, userID uint) ([]*models.CartItem, error) {
	var cartItems []*models.CartItem
	
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

func (r *cartRepository) GetByID(ctx context.Context, id uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	
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

func (r *cartRepository) GetByUserAndFurniture(ctx context.Context, userID uint, furnitureID uint) (*models.CartItem, error) {
	var cartItem models.CartItem
	
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

func (r *cartRepository) Create(ctx context.Context, cartItem *models.CartItem) error {
	return r.db.WithContext(ctx).Create(cartItem).Error
}

func (r *cartRepository) Update(ctx context.Context, cartItem *models.CartItem) error {
	return r.db.WithContext(ctx).Save(cartItem).Error
}

func (r *cartRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.CartItem{}, id).Error
}

func (r *cartRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.CartItem{}).Error
}

func (r *cartRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.CartItem{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
