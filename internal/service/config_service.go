package service

import (
	"context"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/dto/response"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ConfigService 配置服务接口
type ConfigService interface {
	GetConfig(ctx context.Context) (*response.ConfigResponse, error)
	GetRegions(ctx context.Context) (*response.RegionsResponse, error)
	GetPropertyTypes(ctx context.Context) (*response.PropertyTypesResponse, error)
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
func (s *configService) GetConfig(ctx context.Context) (*response.ConfigResponse, error) {
	config := &response.ConfigResponse{
		System: &response.SystemConfig{
			Version:     "1.0.0",
			Environment: "production",
			APIBaseURL:  "https://api.ajoliving.com",
			WebURL:      "https://ajoliving.com",
			Timezone:    "Asia/Hong_Kong",
			Language:    "zh-HK",
			Currency:    "HKD",
		},
		App: &response.AppConfig{
			Name:         "AJO Living",
			Description:  "香港地產服務平台",
			Logo:         "https://cdn.ajoliving.com/logo.png",
			Favicon:      "https://cdn.ajoliving.com/favicon.ico",
			Copyright:    "© 2025 AJO Living. All rights reserved.",
			SupportEmail: "support@ajoliving.com",
			SupportPhone: "+852 1234 5678",
		},
		Features: &response.FeaturesConfig{
			EnableRegistration:    true,
			EnableSocialLogin:     true,
			EnablePropertyListing: true,
			EnableFurnitureStore:  true,
			EnableMortgage:        true,
			EnableValuation:       true,
			EnableNews:            true,
			EnableSchoolNet:       true,
			EnablePriceIndex:      true,
		},
		API: &response.APIConfig{
			Version:     "v1",
			RateLimit:   60,
			Timeout:     30,
			MaxPageSize: 100,
			AllowedOrigins: []string{
				"https://ajoliving.com",
				"https://www.ajoliving.com",
				"http://localhost:3000",
			},
		},
		UI: &response.UIConfig{
			Theme:        "light",
			PrimaryColor: "#2563EB",
			MapProvider:  "google",
			MapAPIKey:    "", // 从环境变量读取
			Languages: []response.LanguageOption{
				{Code: "zh-HK", Name: "繁體中文", Flag: "🇭🇰"},
				{Code: "zh-CN", Name: "简体中文", Flag: "🇨🇳"},
				{Code: "en", Name: "English", Flag: "🇺🇸"},
			},
		},
	}

	return config, nil
}

// GetRegions 获取区域配置
func (s *configService) GetRegions(ctx context.Context) (*response.RegionsResponse, error) {
	// 查询所有地区并按区域分组
	var districts []model.District
	if err := s.db.WithContext(ctx).Order("display_order ASC").Find(&districts).Error; err != nil {
		s.logger.Error("查询地区失败", zap.Error(err))
		return nil, err
	}

	// 统计每个地区的房产和屋苑数量
	propertyCountMap := make(map[uint]int)
	estateCountMap := make(map[uint]int)

	type CountResult struct {
		DistrictID uint
		Count      int64
	}

	// 统计房产数量
	var propertyCounts []CountResult
	s.db.WithContext(ctx).Model(&model.Property{}).
		Select("district_id, COUNT(*) as count").
		Group("district_id").
		Scan(&propertyCounts)
	for _, pc := range propertyCounts {
		propertyCountMap[pc.DistrictID] = int(pc.Count)
	}

	// 统计屋苑数量
	var estateCounts []CountResult
	s.db.WithContext(ctx).Model(&model.Estate{}).
		Select("district_id, COUNT(*) as count").
		Group("district_id").
		Scan(&estateCounts)
	for _, ec := range estateCounts {
		estateCountMap[ec.DistrictID] = int(ec.Count)
	}

	// 按区域分组
	regionMap := make(map[string]*response.RegionConfig)
	regionOrder := []string{"HK", "KLN", "NT"}
	regionNames := map[string]struct {
		ZhHant string
		ZhHans string
		En     string
		Type   string
		Order  int
	}{
		"HK":  {"香港島", "香港岛", "Hong Kong Island", "island", 1},
		"KLN": {"九龍", "九龙", "Kowloon", "peninsula", 2},
		"NT":  {"新界", "新界", "New Territories", "territories", 3},
	}

	// 初始化区域
	for code, info := range regionNames {
		regionMap[code] = &response.RegionConfig{
			Code:         code,
			NameZhHant:   info.ZhHant,
			NameZhHans:   info.ZhHans,
			NameEn:       info.En,
			Type:         info.Type,
			DisplayOrder: info.Order,
			Districts:    []*response.DistrictConfig{},
		}
	}

	// 将地区归类到区域
	for _, district := range districts {
		regionCode := getRegionCode(func() string { if district.NameEn != nil { return *district.NameEn }; return "" }())
		if region, ok := regionMap[regionCode]; ok {
			districtConfig := &response.DistrictConfig{
				ID:            district.ID,
				RegionID:      0, // 可以添加 region_id 字段到 District 模型
				// Code:          district.Code, // TODO: 添加Code字段到District model
				NameZhHant:    district.NameZhHant,
				NameZhHans:    func() string { if district.NameZhHans != nil { return *district.NameZhHans }; return "" }(),
				NameEn:        func() string { if district.NameEn != nil { return *district.NameEn }; return "" }(),
				DisplayOrder:  district.SortOrder,
				PropertyCount: propertyCountMap[district.ID],
				EstateCount:   estateCountMap[district.ID],
			}
			region.Districts = append(region.Districts, districtConfig)
		}
	}

	// 构建响应
	var regions []*response.RegionConfig
	for _, code := range regionOrder {
		if region, ok := regionMap[code]; ok {
			regions = append(regions, region)
		}
	}

	return &response.RegionsResponse{
		Regions: regions,
	}, nil
}

// GetPropertyTypes 获取房产类型配置
func (s *configService) GetPropertyTypes(ctx context.Context) (*response.PropertyTypesResponse, error) {
	// 房产类型配置
	propertyTypes := []*response.PropertyTypeConfig{
		{
			Code:         "apartment",
			NameZhHant:   "公寓",
			NameZhHans:   "公寓",
			NameEn:       "Apartment",
			Icon:         "🏢",
			DisplayOrder: 1,
			Description:  "標準住宅公寓",
		},
		{
			Code:         "villa",
			NameZhHant:   "別墅",
			NameZhHans:   "别墅",
			NameEn:       "Villa",
			Icon:         "🏡",
			DisplayOrder: 2,
			Description:  "獨立別墅",
		},
		{
			Code:         "townhouse",
			NameZhHant:   "聯排別墅",
			NameZhHans:   "联排别墅",
			NameEn:       "Townhouse",
			Icon:         "🏘️",
			DisplayOrder: 3,
			Description:  "聯排式住宅",
		},
		{
			Code:         "penthouse",
			NameZhHant:   "頂層豪宅",
			NameZhHans:   "顶层豪宅",
			NameEn:       "Penthouse",
			Icon:         "🏰",
			DisplayOrder: 4,
			Description:  "頂層豪華公寓",
		},
		{
			Code:         "studio",
			NameZhHant:   "開放式單位",
			NameZhHans:   "开放式单位",
			NameEn:       "Studio",
			Icon:         "🚪",
			DisplayOrder: 5,
			Description:  "開放式設計",
		},
		{
			Code:         "shophouse",
			NameZhHant:   "商住兩用",
			NameZhHans:   "商住两用",
			NameEn:       "Shop House",
			Icon:         "🏬",
			DisplayOrder: 6,
			Description:  "商業住宅混合",
		},
	}

	// 房源类型配置
	listingTypes := []*response.ListingTypeConfig{
		{
			Code:       "rent",
			NameZhHant: "租賃",
			NameZhHans: "租赁",
			NameEn:     "For Rent",
			Icon:       "🔑",
			Color:      "#10B981",
		},
		{
			Code:       "sale",
			NameZhHant: "出售",
			NameZhHans: "出售",
			NameEn:     "For Sale",
			Icon:       "💰",
			Color:      "#F59E0B",
		},
	}

	// 状态配置
	statuses := []*response.StatusConfig{
		{
			Code:        "active",
			NameZhHant:  "活躍",
			NameZhHans:  "活跃",
			NameEn:      "Active",
			Color:       "#22C55E",
			Description: "正在出租/出售",
		},
		{
			Code:        "pending",
			NameZhHant:  "待審核",
			NameZhHans:  "待审核",
			NameEn:      "Pending",
			Color:       "#F59E0B",
			Description: "等待審核中",
		},
		{
			Code:        "sold",
			NameZhHant:  "已售",
			NameZhHans:  "已售",
			NameEn:      "Sold",
			Color:       "#EF4444",
			Description: "已成功售出",
		},
		{
			Code:        "rented",
			NameZhHant:  "已租",
			NameZhHans:  "已租",
			NameEn:      "Rented",
			Color:       "#3B82F6",
			Description: "已成功出租",
		},
		{
			Code:        "inactive",
			NameZhHant:  "未啟用",
			NameZhHans:  "未启用",
			NameEn:      "Inactive",
			Color:       "#9CA3AF",
			Description: "暫時下架",
		},
		{
			Code:        "expired",
			NameZhHant:  "已過期",
			NameZhHans:  "已过期",
			NameEn:      "Expired",
			Color:       "#6B7280",
			Description: "刊登已過期",
		},
	}

	return &response.PropertyTypesResponse{
		PropertyTypes: propertyTypes,
		ListingTypes:  listingTypes,
		Statuses:      statuses,
	}, nil
}

// getRegionCode 根据地区英文名获取区域代码
func getRegionCode(districtNameEn string) string {
	hkIslandDistricts := map[string]bool{
		"Central and Western": true,
		"Wan Chai":            true,
		"Eastern":             true,
		"Southern":            true,
	}

	kowloonDistricts := map[string]bool{
		"Yau Tsim Mong":  true,
		"Sham Shui Po":   true,
		"Kowloon City":   true,
		"Wong Tai Sin":   true,
		"Kwun Tong":      true,
	}

	if hkIslandDistricts[districtNameEn] {
		return "HK"
	} else if kowloonDistricts[districtNameEn] {
		return "KLN"
	}
	return "NT" // New Territories
}
