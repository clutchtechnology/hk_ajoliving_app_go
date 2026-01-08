package databases

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/gorm"
)

type CartRepo struct {
	db *gorm.DB
}

func NewCartRepo(db *gorm.DB) *CartRepo {
	return &CartRepo{db: db}
}

// GetUserCart 获取用户购物车
func (r *CartRepo) GetUserCart(ctx context.Context, userID uint) ([]*models.CartItem, error) {
	var items []*models.CartItem
	err := r.db.WithContext(ctx).
		Preload("Furniture").
		Preload("Furniture.Category").
		Preload("Furniture.DeliveryDistrict").
		Preload("Furniture.Images", "is_cover = ?", true).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

// FindByID 根据ID查找购物车项
func (r *CartRepo) FindByID(ctx context.Context, id uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).
		Preload("Furniture").
		First(&item, id).Error
	return &item, err
}

// FindByUserAndFurniture 查找用户特定家具的购物车项
func (r *CartRepo) FindByUserAndFurniture(ctx context.Context, userID, furnitureID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND furniture_id = ?", userID, furnitureID).
		First(&item).Error
	return &item, err
}

// Create 创建购物车项
func (r *CartRepo) Create(ctx context.Context, item *models.CartItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

// Update 更新购物车项
func (r *CartRepo) Update(ctx context.Context, item *models.CartItem) error {
	return r.db.WithContext(ctx).Save(item).Error
}

// Delete 删除购物车项
func (r *CartRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.CartItem{}, id).Error
}

// ClearUserCart 清空用户购物车
func (r *CartRepo) ClearUserCart(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.CartItem{}).Error
}

// CountUserCartItems 统计用户购物车项数量
func (r *CartRepo) CountUserCartItems(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.CartItem{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}
