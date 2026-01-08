package routes

import (
	"github.com/clutchtechnology/hk_ajoliving_app_go/controllers"
	"github.com/clutchtechnology/hk_ajoliving_app_go/middlewares"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(
	r *gin.Engine,
	healthCtrl *controllers.HealthController,
	authCtrl *controllers.AuthController,
	userCtrl *controllers.UserController,
	propertyCtrl *controllers.PropertyController,
	newDevelopmentCtrl *controllers.NewDevelopmentController,
	servicedApartmentCtrl *controllers.ServicedApartmentController,
	estateCtrl *controllers.EstateController,
	valuationCtrl *controllers.ValuationController,
	furnitureCtrl *controllers.FurnitureController,
	cartCtrl *controllers.CartController,
	schoolNetCtrl *controllers.SchoolNetController,
	schoolCtrl *controllers.SchoolController,
	agentCtrl *controllers.AgentController,
	agencyCtrl *controllers.AgencyController,
	districtCtrl *controllers.DistrictController,
	facilityCtrl *controllers.FacilityController,
	searchCtrl *controllers.SearchController,
	statisticsCtrl *controllers.StatisticsController,
) {
	// API v1 路由组
	v1 := r.Group("/api/v1")

	// ========== 基础路由（无需认证） ==========
	v1.GET("/health", healthCtrl.HealthCheck)     // 健康检查
	v1.GET("/version", healthCtrl.Version)        // 版本信息

	// ========== 认证路由（无需认证） ==========
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/register", authCtrl.Register) // 用户注册
		authGroup.POST("/login", authCtrl.Login)       // 用户登录
		authGroup.POST("/logout", authCtrl.Logout)     // 用户登出（可选认证）
	}

	// ========== 用户路由（需要认证） ==========
	userGroup := v1.Group("/users")
	userGroup.Use(middlewares.JWTAuth()) // 使用 JWT 认证中间件
	{
		userGroup.GET("/me", userCtrl.GetCurrentUser)          // 获取当前用户信息
		userGroup.PUT("/me", userCtrl.UpdateCurrentUser)       // 更新当前用户信息
		userGroup.GET("/me/listings", userCtrl.GetMyListings)  // 获取我的发布
	}

	// ========== 房产路由 ==========
	propertyGroup := v1.Group("/properties")
	{
		// 公开接口（无需认证）
		propertyGroup.GET("", propertyCtrl.ListProperties)                    // 房产列表
		propertyGroup.GET("/featured", propertyCtrl.GetFeaturedProperties)    // 精选房源
		propertyGroup.GET("/hot", propertyCtrl.GetHotProperties)              // 热门房源
		propertyGroup.GET("/:id", propertyCtrl.GetProperty)                   // 房产详情
		propertyGroup.GET("/:id/similar", propertyCtrl.GetSimilarProperties)  // 相似房源

		// 买房分类
		buyGroup := propertyGroup.Group("/buy")
		{
			buyGroup.GET("", propertyCtrl.ListBuyProperties)                    // 买房房源列表
			buyGroup.GET("/new", propertyCtrl.ListNewProperties)                // 新房列表
			buyGroup.GET("/secondhand", propertyCtrl.ListSecondhandProperties)  // 二手房列表
		}

		// 租房分类
		rentGroup := propertyGroup.Group("/rent")
		{
			rentGroup.GET("", propertyCtrl.ListRentProperties)           // 租房房源列表
			rentGroup.GET("/short-term", propertyCtrl.ListShortTermRent) // 短租房源
			rentGroup.GET("/long-term", propertyCtrl.ListLongTermRent)   // 长租房源
		}

		// 需要认证的接口
		authenticated := propertyGroup.Group("")
		authenticated.Use(middlewares.JWTAuth())
		{
			authenticated.POST("", propertyCtrl.CreateProperty)            // 创建房产
			authenticated.PUT("/:id", propertyCtrl.UpdateProperty)         // 更新房产
			authenticated.DELETE("/:id", propertyCtrl.DeleteProperty)      // 删除房产
		}
	}

	// ========== 新盘路由 ==========
	newPropertyGroup := v1.Group("/new-properties")
	{
		newPropertyGroup.GET("", newDevelopmentCtrl.ListNewDevelopments)         // 新盘列表
		newPropertyGroup.GET("/:id", newDevelopmentCtrl.GetNewDevelopment)       // 新盘详情
		newPropertyGroup.GET("/:id/layouts", newDevelopmentCtrl.GetDevelopmentLayouts) // 户型列表
	}

	// ========== 服务式住宅路由 ==========
	servicedApartmentGroup := v1.Group("/serviced-apartments")
	{
		// 公开接口（无需认证）
		servicedApartmentGroup.GET("", servicedApartmentCtrl.ListServicedApartments)          // 服务式住宅列表
		servicedApartmentGroup.GET("/:id", servicedApartmentCtrl.GetServicedApartment)        // 服务式住宅详情
		servicedApartmentGroup.GET("/:id/units", servicedApartmentCtrl.GetServicedApartmentUnits)   // 房型列表
		servicedApartmentGroup.GET("/:id/images", servicedApartmentCtrl.GetServicedApartmentImages) // 图片列表

		// 需要认证的接口
		authenticated := servicedApartmentGroup.Group("")
		authenticated.Use(middlewares.JWTAuth())
		{
			authenticated.POST("", servicedApartmentCtrl.CreateServicedApartment)             // 创建服务式住宅
			authenticated.PUT("/:id", servicedApartmentCtrl.UpdateServicedApartment)          // 更新服务式住宅
			authenticated.DELETE("/:id", servicedApartmentCtrl.DeleteServicedApartment)       // 删除服务式住宅
		}
	}

	// ========== 屋苑路由 ==========
	estateGroup := v1.Group("/estates")
	{
		// 公开接口（无需认证）
		estateGroup.GET("", estateCtrl.ListEstates)                           // 屋苑列表
		estateGroup.GET("/featured", estateCtrl.GetFeaturedEstates)           // 精选屋苑
		estateGroup.GET("/:id", estateCtrl.GetEstate)                         // 屋苑详情
		estateGroup.GET("/:id/properties", estateCtrl.GetEstateProperties)    // 屋苑内房源列表
		estateGroup.GET("/:id/images", estateCtrl.GetEstateImages)            // 屋苑图片
		estateGroup.GET("/:id/facilities", estateCtrl.GetEstateFacilities)    // 屋苑设施
		estateGroup.GET("/:id/transactions", estateCtrl.GetEstateTransactions) // 屋苑成交记录
		estateGroup.GET("/:id/statistics", estateCtrl.GetEstateStatistics)    // 屋苑统计数据

		// 需要认证的接口
		authenticated := estateGroup.Group("")
		authenticated.Use(middlewares.JWTAuth())
		{
			authenticated.POST("", estateCtrl.CreateEstate)              // 创建屋苑
			authenticated.PUT("/:id", estateCtrl.UpdateEstate)           // 更新屋苑
			authenticated.DELETE("/:id", estateCtrl.DeleteEstate)        // 删除屋苑
		}
	}

	// ========== 物业估价路由 ==========
	valuationGroup := v1.Group("/valuation")
	{
		// 公开接口（无需认证）
		valuationGroup.GET("", valuationCtrl.ListValuations)                           // 获取屋苑估价列表
		valuationGroup.GET("/search", valuationCtrl.SearchValuations)                  // 搜索屋苑估价
		valuationGroup.GET("/:estateId", valuationCtrl.GetEstateValuation)             // 获取指定屋苑估价参考
		valuationGroup.GET("/districts/:districtId", valuationCtrl.GetDistrictValuations) // 获取地区屋苑估价列表
	}

	// ========== 家具商城路由 ==========
	furnitureGroup := v1.Group("/furniture")
	{
		// 公开接口（无需认证）
		furnitureGroup.GET("", furnitureCtrl.ListFurniture)                    // 家具列表
		furnitureGroup.GET("/categories", furnitureCtrl.GetFurnitureCategories) // 家具分类
		furnitureGroup.GET("/featured", furnitureCtrl.GetFeaturedFurniture)    // 精选家具
		furnitureGroup.GET("/:id", furnitureCtrl.GetFurniture)                 // 家具详情
		furnitureGroup.GET("/:id/images", furnitureCtrl.GetFurnitureImages)    // 家具图片

		// 需要认证的接口
		authenticated := furnitureGroup.Group("")
		authenticated.Use(middlewares.JWTAuth())
		{
			authenticated.POST("", furnitureCtrl.CreateFurniture)                  // 发布家具
			authenticated.PUT("/:id", furnitureCtrl.UpdateFurniture)               // 更新家具
			authenticated.DELETE("/:id", furnitureCtrl.DeleteFurniture)            // 删除家具
			authenticated.PUT("/:id/status", furnitureCtrl.UpdateFurnitureStatus)  // 更新家具状态
		}
	}

	// ========== 购物车路由（需要认证） ==========
	cartGroup := v1.Group("/cart")
	cartGroup.Use(middlewares.JWTAuth())
	{
		cartGroup.GET("", cartCtrl.GetCart)                       // 获取购物车
		cartGroup.DELETE("", cartCtrl.ClearCart)                  // 清空购物车
		cartGroup.POST("/items", cartCtrl.AddToCart)             // 添加到购物车
		cartGroup.PUT("/items/:id", cartCtrl.UpdateCartItem)     // 更新购物车项
		cartGroup.DELETE("/items/:id", cartCtrl.RemoveFromCart)  // 移除购物车项
	}

	// ========== 校网路由（公开） ==========
	schoolNetGroup := v1.Group("/school-nets")
	{
		schoolNetGroup.GET("", schoolNetCtrl.ListSchoolNets)                     // 校网列表
		schoolNetGroup.GET("/search", schoolNetCtrl.SearchSchoolNets)            // 搜索校网
		schoolNetGroup.GET("/:id", schoolNetCtrl.GetSchoolNet)                   // 校网详情
		schoolNetGroup.GET("/:id/schools", schoolNetCtrl.GetSchoolsInNet)        // 校网内学校
		schoolNetGroup.GET("/:id/properties", schoolNetCtrl.GetPropertiesInNet)  // 校网内房源
		schoolNetGroup.GET("/:id/estates", schoolNetCtrl.GetEstatesInNet)        // 校网内屋苑
	}

	// ========== 学校路由（公开） ==========
	schoolGroup := v1.Group("/schools")
	{
		schoolGroup.GET("", schoolCtrl.ListSchools)                   // 学校列表
		schoolGroup.GET("/search", schoolCtrl.SearchSchools)          // 搜索学校
		schoolGroup.GET("/:id", schoolCtrl.GetSchool)                 // 学校详情
		schoolGroup.GET("/:id/school-net", schoolCtrl.GetSchoolNet)   // 获取学校所属校网
	}

	// ========== 代理人路由（公开） ==========
	agentGroup := v1.Group("/agents")
	{
		agentGroup.GET("", agentCtrl.ListAgents)                     // 代理人列表
		agentGroup.GET("/:id", agentCtrl.GetAgent)                   // 代理人详情
		agentGroup.GET("/:id/properties", agentCtrl.GetAgentProperties) // 代理人房源列表
		agentGroup.POST("/:id/contact", agentCtrl.ContactAgent)      // 联系代理人
	}

	// ========== 代理公司路由（公开） ==========
	agencyGroup := v1.Group("/agencies")
	{
		agencyGroup.GET("/search", agencyCtrl.SearchAgencies)           // 搜索代理公司（需在 :id 前）
		agencyGroup.GET("", agencyCtrl.ListAgencies)                    // 代理公司列表
		agencyGroup.GET("/:id", agencyCtrl.GetAgency)                   // 代理公司详情
		agencyGroup.GET("/:id/properties", agencyCtrl.GetAgencyProperties) // 代理公司房源列表
		agencyGroup.POST("/:id/contact", agencyCtrl.ContactAgency)      // 联系代理公司
	}

	// ========== 地区路由（公开） ==========
	districtGroup := v1.Group("/districts")
	{
		districtGroup.GET("", districtCtrl.ListDistricts)                  // 地区列表
		districtGroup.GET("/:id", districtCtrl.GetDistrict)                // 地区详情
		districtGroup.GET("/:id/properties", districtCtrl.GetDistrictProperties) // 地区房源
		districtGroup.GET("/:id/estates", districtCtrl.GetDistrictEstates)       // 地区屋苑
		districtGroup.GET("/:id/statistics", districtCtrl.GetDistrictStatistics) // 地区统计
	}

	// ========== 设施路由 ==========
	facilityGroup := v1.Group("/facilities")
	{
		facilityGroup.GET("", facilityCtrl.ListFacilities)                // 设施列表（公开）
		facilityGroup.GET("/:id", facilityCtrl.GetFacility)              // 设施详情（公开）
		facilityGroup.POST("", middlewares.JWTAuth(), facilityCtrl.CreateFacility)   // 创建设施（需认证）
		facilityGroup.PUT("/:id", middlewares.JWTAuth(), facilityCtrl.UpdateFacility) // 更新设施（需认证）
		facilityGroup.DELETE("/:id", middlewares.JWTAuth(), facilityCtrl.DeleteFacility) // 删除设施（需认证）
	}

	// ========== 搜索路由（公开） ==========
	searchGroup := v1.Group("/search")
	{
		searchGroup.GET("", searchCtrl.GlobalSearch)                    // 全局搜索
		searchGroup.GET("/properties", searchCtrl.SearchProperties)     // 搜索房产
		searchGroup.GET("/estates", searchCtrl.SearchEstates)           // 搜索屋苑
		searchGroup.GET("/agents", searchCtrl.SearchAgents)             // 搜索代理人
		searchGroup.GET("/suggestions", searchCtrl.GetSearchSuggestions) // 搜索建议
		searchGroup.GET("/history", searchCtrl.GetSearchHistory)        // 搜索历史（可选认证）
	}

	// ========== 统计分析路由（公开） ==========
	statisticsGroup := v1.Group("/statistics")
	{
		statisticsGroup.GET("/overview", statisticsCtrl.GetOverviewStatistics)       // 总览统计
		statisticsGroup.GET("/properties", statisticsCtrl.GetPropertyStatistics)     // 房产统计
		statisticsGroup.GET("/transactions", statisticsCtrl.GetTransactionStatistics) // 成交统计
		statisticsGroup.GET("/users", statisticsCtrl.GetUserStatistics)              // 用户统计
	}
}
