package request

// GetOverviewStatisticsRequest 获取总览统计请求
type GetOverviewStatisticsRequest struct {
	Period string `form:"period" binding:"omitempty,oneof=day week month year"` // 统计周期
}

// GetPropertyStatisticsRequest 获取房产统计请求
type GetPropertyStatisticsRequest struct {
	ListingType *string `form:"listing_type" binding:"omitempty,oneof=rent sale"`
	DistrictID  *uint   `form:"district_id"`
	StartDate   *string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate     *string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Period      string  `form:"period,default=month" binding:"omitempty,oneof=day week month year"`
}

// GetTransactionStatisticsRequest 获取成交统计请求
type GetTransactionStatisticsRequest struct {
	DistrictID *uint   `form:"district_id"`
	EstateID   *uint   `form:"estate_id"`
	StartDate  *string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate    *string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Period     string  `form:"period,default=month" binding:"omitempty,oneof=day week month year"`
}

// GetUserStatisticsRequest 获取用户统计请求
type GetUserStatisticsRequest struct {
	StartDate *string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   *string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Period    string  `form:"period,default=month" binding:"omitempty,oneof=day week month year"`
}
