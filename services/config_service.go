package services

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ConfigService 配置服务接口
type ConfigService interface {
	GetConfig(ctx context.Context) (*map[string]interface{}, error)
	GetRegions(ctx context.Context) (*[]models.District, error)
	GetPropertyTypes(ctx context.Context) (*map[string]interface{}, error)
}

type configService struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewConfigService 创建配置服务实例
func NewConfigService(db *gorm.DB, logger *zap.Logger) ConfigService {
	return &configService{
		db:     db,
		logger: logger,
	}
}

// GetConfig 获取系统配置
func (s *configService) GetConfig(ctx context.Context) (*map[string]interface{}, error) {
	config := map[string]interface{}{
		"system": map[string]interface{}{
			"version":     "1.0.0",
			"environment": "production",
			"api_base_url": "https://api.ajoliving.com",
			"web_url":     "https://ajoliving.com",
			"timezone":    "Asia/Hong_Kong",
			"language":    "zh-HK",
			"currency":    "HKD",
		},
		"app": map[string]interface{}{
			"name":          "AJO Living",
			"description":   "香港地產服務平台",
			"logo":          "https://cdn.ajoliving.com/logo.png",
			"favicon":       "https://cdn.ajoliving.com/favicon.ico",
			"copyright":     "© 2025 AJO Living. All rights reserved.",
			"support_email": "support@ajoliving.com",
			"support_phone": "+852 1234 5678",
		},
		"features": map[string]interface{}{
			"enable_registration":     true,
			"enable_social_login":     true,
			"enable_property_listing": true,
			"enable_furniture_store":  true,
			"enable_mortgage":         true,
			"enable_valuation":        true,
			"enable_news":             true,
			"enable_school_net":       true,
			"enable_price_index":      true,
		},
		"api": map[string]interface{}{
			"version":       "v1",
			"rate_limit":    60,
			"timeout":       30,
			"max_page_size": 100,
			"allowed_origins": []string{
				"https://ajoliving.com",
				"https://www.ajoliving.com",
				"http://localhost:3000",
			},
		},
		"ui": map[string]interface{}{
			"theme":         "light",
			"primary_color": "#2563EB",
			"map_provider":  "google",
			"languages": []map[string]string{
				{"code": "zh-HK", "name": "繁體中文", "flag": "🇭🇰"},
				{"code": "zh-CN", "name": "简体中文", "flag": "🇨🇳"},
				{"code": "en", "name": "English", "flag": "🇺🇸"},
			},
		},
	}

	return &config, nil
}

// GetRegions 获取区域配置
func (s *configService) GetRegions(ctx context.Context) (*[]models.District, error) {
	// 查询所有地区并按区域分组
	var districts []models.District
	if err := s.db.WithContext(ctx).Order("sort_order ASC").Find(&districts).Error; err != nil {
		s.logger.Error("查询地区失败", zap.Error(err))
		return nil, err
	}

	return &districts, nil
}

// GetPropertyTypes 获取房产类型配置
func (s *configService) GetPropertyTypes(ctx context.Context) (*map[string]interface{}, error) {
	// 房产类型配置
	propertyTypes := []map[string]interface{}{
		{
			"code":          "apartment",
			"name_zh_hant":  "公寓",
			"name_zh_hans":  "公寓",
			"name_en":       "Apartment",
			"icon":          "🏢",
			"display_order": 1,
			"description":   "標準住宅公寓",
		},
		{
			"code":          "villa",
			"name_zh_hant":  "別墅",
			"name_zh_hans":  "别墅",
			"name_en":       "Villa",
			"icon":          "🏡",
			"display_order": 2,
			"description":   "獨立別墅",
		},
		{
			"code":          "townhouse",
			"name_zh_hant":  "聯排別墅",
			"name_zh_hans":  "联排别墅",
			"name_en":       "Townhouse",
			"icon":          "🏘️",
			"display_order": 3,
			"description":   "聯排式住宅",
		},
		{
			"code":          "penthouse",
			"name_zh_hant":  "頂層豪宅",
			"name_zh_hans":  "顶层豪宅",
			"name_en":       "Penthouse",
			"icon":          "🏰",
			"display_order": 4,
			"description":   "頂層豪華公寓",
		},
		{
			"code":          "studio",
			"name_zh_hant":  "開放式單位",
			"name_zh_hans":  "开放式单位",
			"name_en":       "Studio",
			"icon":          "🚪",
			"display_order": 5,
			"description":   "開放式設計",
		},
		{
			"code":          "shophouse",
			"name_zh_hant":  "商住兩用",
			"name_zh_hans":  "商住两用",
			"name_en":       "Shop House",
			"icon":          "🏬",
			"display_order": 6,
			"description":   "商業住宅混合",
		},
	}

	// 房源类型配置
	listingTypes := []map[string]interface{}{
		{
			"code":         "rent",
			"name_zh_hant": "租賃",
			"name_zh_hans": "租赁",
			"name_en":      "For Rent",
			"icon":         "🔑",
			"color":        "#10B981",
		},
		{
			"code":         "sale",
			"name_zh_hant": "出售",
			"name_zh_hans": "出售",
			"name_en":      "For Sale",
			"icon":         "💰",
			"color":        "#F59E0B",
		},
	}

	// 状态配置
	statuses := []map[string]interface{}{
		{
			"code":         "active",
			"name_zh_hant": "活躍",
			"name_zh_hans": "活跃",
			"name_en":      "Active",
			"color":        "#22C55E",
			"description":  "正在出租/出售",
		},
		{
			"code":         "pending",
			"name_zh_hant": "待審核",
			"name_zh_hans": "待审核",
			"name_en":      "Pending",
			"color":        "#F59E0B",
			"description":  "等待審核中",
		},
		{
			"code":         "sold",
			"name_zh_hant": "已售",
			"name_zh_hans": "已售",
			"name_en":      "Sold",
			"color":        "#EF4444",
			"description":  "已成功售出",
		},
		{
			"code":         "rented",
			"name_zh_hant": "已租",
			"name_zh_hans": "已租",
			"name_en":      "Rented",
			"color":        "#3B82F6",
			"description":  "已成功出租",
		},
		{
			"code":         "inactive",
			"name_zh_hant": "未啟用",
			"name_zh_hans": "未启用",
			"name_en":      "Inactive",
			"color":        "#9CA3AF",
			"description":  "暫時下架",
		},
		{
			"code":         "expired",
			"name_zh_hant": "已過期",
			"name_zh_hans": "已过期",
			"name_en":      "Expired",
			"color":        "#6B7280",
			"description":  "刊登已過期",
		},
	}

	result := map[string]interface{}{
		"property_types": propertyTypes,
		"listing_types":  listingTypes,
		"statuses":       statuses,
	}

	return &result, nil
}
