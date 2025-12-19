package response

import "time"

// DistrictBasicResponse 地区基本响应
type DistrictBasicResponse struct {
	ID         uint   `json:"id"`
	NameZhHant string `json:"name_zh_hant"`
	NameZhHans string `json:"name_zh_hans"`
	NameEn     string `json:"name_en"`
}

// SchoolNetListItemResponse 校网列表项响应
type SchoolNetListItemResponse struct {
	ID          uint                   `json:"id"`
	NetCode     string                 `json:"net_code"`
	NameZhHant  string                 `json:"name_zh_hant"`
	NameZhHans  string                 `json:"name_zh_hans"`
	NameEn      string                 `json:"name_en"`
	DistrictID  uint                   `json:"district_id"`
	District    *DistrictBasicResponse `json:"district,omitempty"`
	Level       string                 `json:"level"`
	SchoolCount int                    `json:"school_count"`
	CreatedAt   time.Time              `json:"created_at"`
}

// SchoolNetResponse 校网详情响应
type SchoolNetResponse struct {
	ID          uint                   `json:"id"`
	NetCode     string                 `json:"net_code"`
	NameZhHant  string                 `json:"name_zh_hant"`
	NameZhHans  string                 `json:"name_zh_hans"`
	NameEn      string                 `json:"name_en"`
	DistrictID  uint                   `json:"district_id"`
	District    *DistrictBasicResponse `json:"district,omitempty"`
	Description string                 `json:"description"`
	Level       string                 `json:"level"`
	SchoolCount int                    `json:"school_count"`
	MapData     string                 `json:"map_data,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// SchoolListItemResponse 学校列表项响应
type SchoolListItemResponse struct {
	ID            uint                   `json:"id"`
	SchoolNetID   uint                   `json:"school_net_id"`
	DistrictID    uint                   `json:"district_id"`
	NameZhHant    string                 `json:"name_zh_hant"`
	NameZhHans    string                 `json:"name_zh_hans"`
	NameEn        string                 `json:"name_en"`
	SchoolCode    string                 `json:"school_code,omitempty"`
	Category      string                 `json:"category"`
	CategoryName  string                 `json:"category_name"`
	Level         string                 `json:"level"`
	LevelName     string                 `json:"level_name"`
	Gender        string                 `json:"gender,omitempty"`
	Religion      string                 `json:"religion,omitempty"`
	Address       string                 `json:"address"`
	Phone         string                 `json:"phone,omitempty"`
	Website       string                 `json:"website,omitempty"`
	StudentCount  int                    `json:"student_count"`
	Rating        float64                `json:"rating"`
	LogoURL       string                 `json:"logo_url,omitempty"`
	District      *DistrictBasicResponse `json:"district,omitempty"`
}

// SchoolResponse 学校详情响应
type SchoolResponse struct {
	ID            uint                   `json:"id"`
	SchoolNetID   uint                   `json:"school_net_id"`
	SchoolNet     *SchoolNetListItemResponse `json:"school_net,omitempty"`
	DistrictID    uint                   `json:"district_id"`
	District      *DistrictBasicResponse `json:"district,omitempty"`
	NameZhHant    string                 `json:"name_zh_hant"`
	NameZhHans    string                 `json:"name_zh_hans"`
	NameEn        string                 `json:"name_en"`
	SchoolCode    string                 `json:"school_code,omitempty"`
	Category      string                 `json:"category"`
	CategoryName  string                 `json:"category_name"`
	Level         string                 `json:"level"`
	LevelName     string                 `json:"level_name"`
	Gender        string                 `json:"gender,omitempty"`
	Religion      string                 `json:"religion,omitempty"`
	Address       string                 `json:"address"`
	Phone         string                 `json:"phone,omitempty"`
	Email         string                 `json:"email,omitempty"`
	Website       string                 `json:"website,omitempty"`
	Principal     string                 `json:"principal,omitempty"`
	FoundedYear   int                    `json:"founded_year,omitempty"`
	StudentCount  int                    `json:"student_count"`
	TeacherCount  int                    `json:"teacher_count"`
	Rating        float64                `json:"rating"`
	Features      string                 `json:"features,omitempty"`
	Facilities    string                 `json:"facilities,omitempty"`
	Latitude      float64                `json:"latitude,omitempty"`
	Longitude     float64                `json:"longitude,omitempty"`
	LogoURL       string                 `json:"logo_url,omitempty"`
	CoverImageURL string                 `json:"cover_image_url,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// SchoolNetWithPropertiesCountResponse 校网带房源数量响应
type SchoolNetWithPropertiesCountResponse struct {
	SchoolNetResponse
	PropertiesCount int `json:"properties_count"`
	EstatesCount    int `json:"estates_count"`
}
