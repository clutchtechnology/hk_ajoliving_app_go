package main

import (
	"log"
	"os"

	"github.com/clutchtechnology/hk_ajoliving_app_go/controllers"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/middlewares"
	"github.com/clutchtechnology/hk_ajoliving_app_go/routes"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables")
	}

	// 初始化数据库
	if err := databases.InitDB(); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer databases.CloseDB()

	// 初始化仓储层
	userRepo := databases.NewUserRepo(databases.DB)
	propertyRepo := databases.NewPropertyRepo(databases.DB)
	newDevelopmentRepo := databases.NewNewDevelopmentRepo(databases.DB)
	servicedApartmentRepo := databases.NewServicedApartmentRepo(databases.DB)
	estateRepo := databases.NewEstateRepo(databases.DB)
	valuationRepo := databases.NewValuationRepo(databases.DB)
	furnitureRepo := databases.NewFurnitureRepo(databases.DB)
	cartRepo := databases.NewCartRepo(databases.DB)
	schoolNetRepo := databases.NewSchoolNetRepo(databases.DB)
	schoolRepo := databases.NewSchoolRepo(databases.DB)
	agentRepo := databases.NewAgentRepo(databases.DB)
	agencyRepo := databases.NewAgencyRepo(databases.DB)
	districtRepo := databases.NewDistrictRepo(databases.DB)
	facilityRepo := databases.NewFacilityRepo(databases.DB)
	searchRepo := databases.NewSearchRepo(databases.DB)
	statisticsRepo := databases.NewStatisticsRepo(databases.DB)

	// 初始化服务层
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo, propertyRepo)
	propertyService := services.NewPropertyService(propertyRepo)
	newDevelopmentService := services.NewNewDevelopmentService(newDevelopmentRepo)
	servicedApartmentService := services.NewServicedApartmentService(servicedApartmentRepo)
	estateService := services.NewEstateService(estateRepo)
	valuationService := services.NewValuationService(valuationRepo)
	furnitureService := services.NewFurnitureService(furnitureRepo)
	cartService := services.NewCartService(cartRepo, furnitureRepo)
	schoolNetService := services.NewSchoolNetService(schoolNetRepo)
	schoolService := services.NewSchoolService(schoolRepo)
	agentService := services.NewAgentService(agentRepo)
	agencyService := services.NewAgencyService(agencyRepo)
	districtService := services.NewDistrictService(districtRepo)
	facilityService := services.NewFacilityService(facilityRepo)
	searchService := services.NewSearchService(searchRepo)
	statisticsService := services.NewStatisticsService(statisticsRepo)

	// 初始化控制器层
	healthCtrl := controllers.NewHealthController()
	authCtrl := controllers.NewAuthController(authService)
	userCtrl := controllers.NewUserController(userService)
	propertyCtrl := controllers.NewPropertyController(propertyService)
	newDevelopmentCtrl := controllers.NewNewDevelopmentController(newDevelopmentService)
	servicedApartmentCtrl := controllers.NewServicedApartmentController(servicedApartmentService)
	estateCtrl := controllers.NewEstateController(estateService)
	valuationCtrl := controllers.NewValuationController(valuationService)
	furnitureCtrl := controllers.NewFurnitureController(furnitureService)
	cartCtrl := controllers.NewCartController(cartService)
	schoolNetCtrl := controllers.NewSchoolNetController(schoolNetService)
	schoolCtrl := controllers.NewSchoolController(schoolService)
	agentCtrl := controllers.NewAgentController(agentService)
	agencyCtrl := controllers.NewAgencyController(agencyService)
	districtCtrl := controllers.NewDistrictController(districtService)
	facilityCtrl := controllers.NewFacilityController(facilityService)
	searchCtrl := controllers.NewSearchController(searchService)
	statisticsCtrl := controllers.NewStatisticsController(statisticsService)

	// 设置 Gin 模式
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	// 创建 Gin 引擎
	r := gin.Default()

	// 应用 CORS 中间件
	r.Use(middlewares.CORS())

	// 设置路由
	routes.SetupRoutes(r, healthCtrl, authCtrl, userCtrl, propertyCtrl, newDevelopmentCtrl, servicedApartmentCtrl, estateCtrl, valuationCtrl, furnitureCtrl, cartCtrl, schoolNetCtrl, schoolCtrl, agentCtrl, agencyCtrl, districtCtrl, facilityCtrl, searchCtrl, statisticsCtrl)

	// 启动服务器
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
