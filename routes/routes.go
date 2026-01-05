package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"github.com/clutchtechnology/hk_ajoliving_app_go/controllers"
	"github.com/clutchtechnology/hk_ajoliving_app_go/middlewares"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
)

// SetupRouter 设置路由
func SetupRouter(
	baseHandler *controllers.BaseHandler,
	authHandler *controllers.AuthHandler,
	userHandler *controllers.UserHandler,
	propertyHandler *controllers.PropertyHandler,
	newPropertyHandler *controllers.NewPropertyHandler,
	servicedApartmentHandler *controllers.ServicedApartmentHandler,
	estateHandler *controllers.EstateHandler,
	valuationHandler *controllers.ValuationHandler,
	furnitureHandler *controllers.FurnitureHandler,
	cartHandler *controllers.CartHandler,
	mortgageHandler *controllers.MortgageHandler,
	newsHandler *controllers.NewsHandler,
	schoolHandler *controllers.SchoolHandler,
	agentHandler *controllers.AgentHandler,
	agencyHandler *controllers.AgencyHandler,
	priceIndexHandler *controllers.PriceIndexHandler,
	facilityHandler *controllers.FacilityHandler,
	searchHandler *controllers.SearchHandler,
	statisticsHandler *controllers.StatisticsHandler,
	configHandler *controllers.ConfigHandler,
	jwtManager *utils.JWTManager,
	logger *zap.Logger,
) *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS())

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 基础路由（无需认证）
		v1.GET("/health", baseHandler.HealthCheck)
		v1.GET("/version", baseHandler.Version)

		// 认证路由（无需认证）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
			auth.POST("/verify-code", authHandler.VerifyCode)
			
			// 需要认证的认证路由
			authProtected := auth.Group("")
			authProtected.Use(middleware.JWTAuth(jwtManager))
			{
				authProtected.POST("/logout", authHandler.Logout)
			}
		}

		// 用户模块（需要认证）
		users := v1.Group("/users")
		users.Use(middleware.JWTAuth(jwtManager))
		{
			users.GET("/me", userHandler.GetCurrentUser)
			users.PUT("/me", userHandler.UpdateCurrentUser)
			users.PUT("/me/password", userHandler.ChangePassword)
			users.GET("/me/listings", userHandler.GetMyListings)
			users.PUT("/me/settings", userHandler.UpdateSettings)
		}

		// 房产模块
		properties := v1.Group("/properties")
		{
			// 公开路由（无需认证）
			properties.GET("", propertyHandler.ListProperties)
			properties.GET("/featured", propertyHandler.GetFeaturedProperties)
			properties.GET("/hot", propertyHandler.GetHotProperties)
			properties.GET("/:id", propertyHandler.GetProperty)
			properties.GET("/:id/similar", propertyHandler.GetSimilarProperties)
			
			// 买房相关路由
			properties.GET("/buy", propertyHandler.ListBuyProperties)
			properties.GET("/buy/new", propertyHandler.ListNewProperties)
			properties.GET("/buy/secondhand", propertyHandler.ListSecondhandProperties)
			
			// 租房相关路由
			properties.GET("/rent", propertyHandler.ListRentProperties)
			properties.GET("/rent/short-term", propertyHandler.ListShortTermRent)
			properties.GET("/rent/long-term", propertyHandler.ListLongTermRent)

			// 需要认证的路由
			propertiesProtected := properties.Group("")
			propertiesProtected.Use(middleware.JWTAuth(jwtManager))
			{
				propertiesProtected.POST("", propertyHandler.CreateProperty)
				propertiesProtected.PUT("/:id", propertyHandler.UpdateProperty)
				propertiesProtected.DELETE("/:id", propertyHandler.DeleteProperty)
			}
		}

		// 新楼盘模块
		newProperties := v1.Group("/new-properties")
		{
			newProperties.GET("", newPropertyHandler.ListNewDevelopments)
			newProperties.GET("/featured", newPropertyHandler.GetFeaturedNewDevelopments)
			newProperties.GET("/:id", newPropertyHandler.GetNewDevelopment)
			newProperties.GET("/:id/units", newPropertyHandler.GetDevelopmentUnits)
		}

		// 服务式住宅模块
		servicedApartments := v1.Group("/serviced-apartments")
		{
			// 公开路由（无需认证）
			servicedApartments.GET("", servicedApartmentHandler.ListServicedApartments)
			servicedApartments.GET("/featured", servicedApartmentHandler.GetFeaturedApartments)
			servicedApartments.GET("/:id", servicedApartmentHandler.GetServicedApartment)
			servicedApartments.GET("/:id/units", servicedApartmentHandler.GetApartmentUnits)

			// 需要认证的路由
			servicedApartmentsProtected := servicedApartments.Group("")
			servicedApartmentsProtected.Use(middleware.JWTAuth(jwtManager))
			{
				servicedApartmentsProtected.POST("", servicedApartmentHandler.CreateServicedApartment)
				servicedApartmentsProtected.PUT("/:id", servicedApartmentHandler.UpdateServicedApartment)
				servicedApartmentsProtected.DELETE("/:id", servicedApartmentHandler.DeleteServicedApartment)
			}
		}

		// 屋苑模块
		estates := v1.Group("/estates")
		{
			// 公开路由（无需认证）
			estates.GET("", estateHandler.ListEstates)
			estates.GET("/featured", estateHandler.GetFeaturedEstates)
			estates.GET("/:id", estateHandler.GetEstate)
			estates.GET("/:id/properties", estateHandler.GetEstateProperties)
			estates.GET("/:id/statistics", estateHandler.GetEstateStatistics)

			// 需要认证的路由
			estatesProtected := estates.Group("")
			estatesProtected.Use(middleware.JWTAuth(jwtManager))
			{
				estatesProtected.POST("", estateHandler.CreateEstate)
				estatesProtected.PUT("/:id", estateHandler.UpdateEstate)
				estatesProtected.DELETE("/:id", estateHandler.DeleteEstate)
			}
		}

		// 物业估价模块
		valuation := v1.Group("/valuation")
		{
			valuation.GET("", valuationHandler.ListValuations)
			valuation.GET("/:estateId", valuationHandler.GetEstateValuation)
			valuation.GET("/search", valuationHandler.SearchValuations)
			valuation.GET("/districts/:districtId", valuationHandler.GetDistrictValuations)
		}

		// 家具商城模块
		furniture := v1.Group("/furniture")
		{
			// 公开路由（无需认证）
			furniture.GET("", furnitureHandler.ListFurniture)
			furniture.GET("/categories", furnitureHandler.GetFurnitureCategories)
			furniture.GET("/featured", furnitureHandler.GetFeaturedFurniture)
			furniture.GET("/:id", furnitureHandler.GetFurniture)
			furniture.GET("/:id/images", furnitureHandler.GetFurnitureImages)

			// 需要认证的路由
			furnitureProtected := furniture.Group("")
			furnitureProtected.Use(middleware.JWTAuth(jwtManager))
			{
				furnitureProtected.POST("", furnitureHandler.CreateFurniture)
				furnitureProtected.PUT("/:id", furnitureHandler.UpdateFurniture)
				furnitureProtected.DELETE("/:id", furnitureHandler.DeleteFurniture)
				furnitureProtected.PUT("/:id/status", furnitureHandler.UpdateFurnitureStatus)
			}
		}

		// 购物车模块（需要认证）
		cart := v1.Group("/cart")
		cart.Use(middleware.JWTAuth(jwtManager))
		{
			cart.GET("", cartHandler.GetCart)
			cart.POST("/items", cartHandler.AddToCart)
			cart.PUT("/items/:id", cartHandler.UpdateCartItem)
			cart.DELETE("/items/:id", cartHandler.RemoveFromCart)
			cart.DELETE("", cartHandler.ClearCart)
		}

		// 按揭模块
		mortgage := v1.Group("/mortgage")
		{
			// 公开路由（无需认证）
			mortgage.POST("/calculate", mortgageHandler.CalculateMortgage)
			mortgage.GET("/rates", mortgageHandler.GetMortgageRates)
			mortgage.GET("/rates/bank/:bank_id", mortgageHandler.GetBankMortgageRate)
			mortgage.POST("/rates/compare", mortgageHandler.CompareMortgageRates)

			// 需要认证的路由
			mortgageProtected := mortgage.Group("")
			mortgageProtected.Use(middleware.JWTAuth(jwtManager))
			{
				mortgageProtected.POST("/apply", mortgageHandler.ApplyMortgage)
				mortgageProtected.GET("/applications", mortgageHandler.GetMortgageApplications)
				mortgageProtected.GET("/applications/:id", mortgageHandler.GetMortgageApplication)
			}
		}

		// 新闻资讯模块（所有路由公开，无需认证）
		news := v1.Group("/news")
		{
			news.GET("", newsHandler.ListNews)
			news.GET("/categories", newsHandler.GetNewsCategories)
			news.GET("/hot", newsHandler.GetHotNews)
			news.GET("/featured", newsHandler.GetFeaturedNews)
			news.GET("/latest", newsHandler.GetLatestNews)
			news.GET("/:id", newsHandler.GetNews)
			news.GET("/:id/related", newsHandler.GetRelatedNews)
		}

		// 校网模块（所有路由公开，无需认证）
		schoolNets := v1.Group("/school-nets")
		{
			schoolNets.GET("", schoolHandler.ListSchoolNets)
			schoolNets.GET("/search", schoolHandler.SearchSchoolNets)
			schoolNets.GET("/:id", schoolHandler.GetSchoolNet)
			schoolNets.GET("/:id/schools", schoolHandler.GetSchoolsInNet)
			schoolNets.GET("/:id/properties", schoolHandler.GetPropertiesInNet)
			schoolNets.GET("/:id/estates", schoolHandler.GetEstatesInNet)
		}

		// 学校模块（所有路由公开，无需认证）
		schools := v1.Group("/schools")
		{
			schools.GET("", schoolHandler.ListSchools)
			schools.GET("/search", schoolHandler.SearchSchools)
			schools.GET("/:id/school-net", schoolHandler.GetSchoolNetBySchoolID)
		}

		// 代理人模块（所有路由公开，无需认证）
		agents := v1.Group("/agents")
		{
			agents.GET("", agentHandler.ListAgents)
			agents.GET("/:id", agentHandler.GetAgent)
			agents.GET("/:id/properties", agentHandler.GetAgentProperties)
			agents.POST("/:id/contact", agentHandler.ContactAgent)
		}

		// 代理公司模块（所有路由公开，无需认证）
		agencies := v1.Group("/agencies")
		{
			agencies.GET("", agencyHandler.ListAgencies)
			agencies.GET("/search", agencyHandler.SearchAgencies)
			agencies.GET("/:id", agencyHandler.GetAgency)
			agencies.GET("/:id/properties", agencyHandler.GetAgencyProperties)
			agencies.POST("/:id/contact", agencyHandler.ContactAgency)
		}

		// 楼价指数模块
		priceIndex := v1.Group("/price-index")
		{
			// 公开路由（无需认证）
			priceIndex.GET("", priceIndexHandler.GetPriceIndex)
			priceIndex.GET("/latest", priceIndexHandler.GetLatestPriceIndex)
			priceIndex.GET("/districts/:districtId", priceIndexHandler.GetDistrictPriceIndex)
			priceIndex.GET("/estates/:estateId", priceIndexHandler.GetEstatePriceIndex)
			priceIndex.GET("/trends", priceIndexHandler.GetPriceTrends)
			priceIndex.GET("/compare", priceIndexHandler.ComparePriceIndex)
			priceIndex.GET("/export", priceIndexHandler.ExportPriceData)
			priceIndex.GET("/history", priceIndexHandler.GetPriceIndexHistory)

			// 需要认证的路由
			priceIndexProtected := priceIndex.Group("")
			priceIndexProtected.Use(middleware.JWTAuth(jwtManager))
			{
				priceIndexProtected.POST("", priceIndexHandler.CreatePriceIndex)
				priceIndexProtected.PUT("/:id", priceIndexHandler.UpdatePriceIndex)
			}
		}

		// 设施模块
		facilities := v1.Group("/facilities")
		{
			// 公开路由（无需认证）
			facilities.GET("", facilityHandler.ListFacilities)
			facilities.GET("/:id", facilityHandler.GetFacility)

			// 需要认证的路由
			facilitiesProtected := facilities.Group("")
			facilitiesProtected.Use(middleware.JWTAuth(jwtManager))
			{
				facilitiesProtected.POST("", facilityHandler.CreateFacility)
				facilitiesProtected.PUT("/:id", facilityHandler.UpdateFacility)
				facilitiesProtected.DELETE("/:id", facilityHandler.DeleteFacility)
			}
		}

		// 搜索模块
		search := v1.Group("/search")
		{
			// 公开路由（无需认证）
			search.GET("", searchHandler.GlobalSearch)
			search.GET("/properties", searchHandler.SearchProperties)
			search.GET("/estates", searchHandler.SearchEstates)
			search.GET("/agents", searchHandler.SearchAgents)
			search.GET("/suggestions", searchHandler.GetSearchSuggestions)

			// 需要认证的路由
			searchProtected := search.Group("")
			searchProtected.Use(middleware.JWTAuth(jwtManager))
			{
				searchProtected.GET("/history", searchHandler.GetSearchHistory)
			}
		}

		// 统计分析模块
		statistics := v1.Group("/statistics")
		{
			// 公开路由（无需认证）
			statistics.GET("/overview", statisticsHandler.GetOverviewStatistics)
			statistics.GET("/properties", statisticsHandler.GetPropertyStatistics)
			statistics.GET("/transactions", statisticsHandler.GetTransactionStatistics)
			statistics.GET("/users", statisticsHandler.GetUserStatistics)
		}

		// 系统配置模块
		config := v1.Group("/config")
		{
			// 公开路由（无需认证）
			config.GET("", configHandler.GetConfig)
			config.GET("/regions", configHandler.GetRegions)
			config.GET("/property-types", configHandler.GetPropertyTypes)
		}
	}

	return r
}

