package response

import "time"

// EstateResponse 屋苑详情响应
type EstateResponse struct {
	ID                           uint               `json:"id"`
	Name                         string             `json:"name"`
	NameEn                       string             `json:"name_en,omitempty"`
	Address                      string             `json:"address"`
	DistrictID                   uint               `json:"district_id"`
	District                     *DistrictResponse  `json:"district,omitempty"`
	TotalBlocks                  int                `json:"total_blocks,omitempty"`
	TotalUnits                   int                `json:"total_units,omitempty"`
	CompletionYear               int                `json:"completion_year,omitempty"`
	Age                          int                `json:"age"`
	Developer                    string             `json:"developer,omitempty"`
	ManagementCompany            string             `json:"management_company,omitempty"`
	PrimarySchoolNet             string             `json:"primary_school_net,omitempty"`
	SecondarySchoolNet           string             `json:"secondary_school_net,omitempty"`
	RecentTransactionsCount      int                `json:"recent_transactions_count"`
	ForSaleCount                 int                `json:"for_sale_count"`
	ForRentCount                 int                `json:"for_rent_count"`
	AvgTransactionPrice          float64            `json:"avg_transaction_price,omitempty"`
	AvgTransactionPriceUpdatedAt *time.Time         `json:"avg_transaction_price_updated_at,omitempty"`
	Description                  string             `json:"description,omitempty"`
	ViewCount                    int                `json:"view_count"`
	FavoriteCount                int                `json:"favorite_count"`
	IsFeatured                   bool               `json:"is_featured"`
	Images                       []EstateImageResponse `json:"images,omitempty"`
	Facilities                   []FacilityResponse `json:"facilities,omitempty"`
	CreatedAt                    time.Time          `json:"created_at"`
	UpdatedAt                    time.Time          `json:"updated_at"`
}

// EstateListItemResponse 屋苑列表项响应
type EstateListItemResponse struct {
	ID                      uint              `json:"id"`
	Name                    string            `json:"name"`
	NameEn                  string            `json:"name_en,omitempty"`
	Address                 string            `json:"address"`
	DistrictID              uint              `json:"district_id"`
	District                *DistrictResponse `json:"district,omitempty"`
	CompletionYear          int               `json:"completion_year,omitempty"`
	Age                     int               `json:"age"`
	RecentTransactionsCount int               `json:"recent_transactions_count"`
	ForSaleCount            int               `json:"for_sale_count"`
	ForRentCount            int               `json:"for_rent_count"`
	AvgTransactionPrice     float64           `json:"avg_transaction_price,omitempty"`
	ViewCount               int               `json:"view_count"`
	FavoriteCount           int               `json:"favorite_count"`
	IsFeatured              bool              `json:"is_featured"`
	CoverImage              string            `json:"cover_image,omitempty"`
	CreatedAt               time.Time         `json:"created_at"`
}

// EstateImageResponse 屋苑图片响应
type EstateImageResponse struct {
	ID        uint   `json:"id"`
	URL       string `json:"url"`
	ImageType string `json:"image_type"`
	Title     string `json:"title,omitempty"`
	SortOrder int    `json:"sort_order"`
}

// EstateStatisticsResponse 屋苑统计数据响应
type EstateStatisticsResponse struct {
	EstateID                uint      `json:"estate_id"`
	EstateName              string    `json:"estate_name"`
	TotalTransactions       int       `json:"total_transactions"`
	RecentTransactions      int       `json:"recent_transactions"`
	AvgTransactionPrice     float64   `json:"avg_transaction_price"`
	MedianTransactionPrice  float64   `json:"median_transaction_price"`
	MinTransactionPrice     float64   `json:"min_transaction_price"`
	MaxTransactionPrice     float64   `json:"max_transaction_price"`
	AvgPricePerSqft         float64   `json:"avg_price_per_sqft"`
	TotalListings           int       `json:"total_listings"`
	ForSaleCount            int       `json:"for_sale_count"`
	ForRentCount            int       `json:"for_rent_count"`
	AvgListingPrice         float64   `json:"avg_listing_price"`
	MedianListingPrice      float64   `json:"median_listing_price"`
	PriceChangePercentage   float64   `json:"price_change_percentage"`
	LastTransactionDate     *time.Time `json:"last_transaction_date,omitempty"`
}

// CreateEstateResponse 创建屋苑响应
type CreateEstateResponse struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}

// UpdateEstateResponse 更新屋苑响应
type UpdateEstateResponse struct {
	ID      uint   `json:"id"`
	Message string `json:"message"`
}
