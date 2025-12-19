package response

import "time"

// UserResponse 用户信息响应
type UserResponse struct {
	ID            uint       `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	FullName      string     `json:"full_name,omitempty"`
	Phone         string     `json:"phone,omitempty"`
	Avatar        string     `json:"avatar,omitempty"`
	UserType      string     `json:"user_type"`
	Status        string     `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
}

// UserSettingsResponse 用户设置响应
type UserSettingsResponse struct {
	Language            string `json:"language"`
	NotificationEnabled bool   `json:"notification_enabled"`
	EmailNotification   bool   `json:"email_notification"`
	SmsNotification     bool   `json:"sms_notification"`
}

// MyListingsResponse 我的发布列表响应
type MyListingsResponse struct {
	Properties []PropertyBriefResponse `json:"properties"`
	Total      int64                   `json:"total"`
}

// PropertyBriefResponse 房产简要信息响应
type PropertyBriefResponse struct {
	ID            uint    `json:"id"`
	PropertyNo    string  `json:"property_no"`
	Title         string  `json:"title"`
	Price         float64 `json:"price"`
	Area          float64 `json:"area"`
	Bedrooms      int     `json:"bedrooms"`
	ListingType   string  `json:"listing_type"`
	Status        string  `json:"status"`
	CoverImage    string  `json:"cover_image,omitempty"`
	ViewCount     int     `json:"view_count"`
	FavoriteCount int     `json:"favorite_count"`
	CreatedAt     string  `json:"created_at"`
}
