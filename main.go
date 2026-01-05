package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/clutchtechnology/hk_ajoliving_app_go/controllers"
	"github.com/clutchtechnology/hk_ajoliving_app_go/databases"
	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"github.com/clutchtechnology/hk_ajoliving_app_go/routes"
	"github.com/clutchtechnology/hk_ajoliving_app_go/services"
	"github.com/clutchtechnology/hk_ajoliving_app_go/tools"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// 初始化日志
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// 设置 Gin 模式
	serverMode := getEnv("SERVER_MODE", "debug")
	gin.SetMode(serverMode)

	// 初始化数据库
	db, err := initDatabase(logger)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// 初始化 JWT 管理器
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	jwtExpireHours := getEnvAsInt("JWT_EXPIRE_HOURS", 24)
	jwtRefreshExpireHours := getEnvAsInt("JWT_REFRESH_EXPIRE_HOURS", 168)
	jwtManager := tools.NewJWTManager(jwtSecret, jwtExpireHours, jwtRefreshExpireHours)

	// 初始化仓库层
	// 初始化仓库层
	userRepo := databases.NewUserRepository(db)
	propertyRepo := databases.NewPropertyRepo(db)
	newPropertyRepo := databases.NewNewPropertyRepository(db)
	servicedApartmentRepo := databases.NewServicedApartmentRepository(db)
	estateRepo := databases.NewEstateRepository(db)
	valuationRepo := databases.NewValuationRepository(db)
	furnitureRepo := databases.NewFurnitureRepository(db)
	cartRepo := databases.NewCartRepository(db)
	mortgageRepo := databases.NewMortgageRepository(db)
	newsRepo := databases.NewNewsRepository(db)
	agentRepo := databases.NewAgentRepository(db)
	agencyRepo := databases.NewAgencyRepository(db)
	schoolRepo := databases.NewSchoolRepository(db)
	priceIndexRepo := databases.NewPriceIndexRepository(db)
	facilityRepo := databases.NewFacilityRepository(db)
	searchRepo := databases.NewSearchRepository(db)
	statisticsRepo := databases.NewStatisticsRepository(db)

	// 初始化服务层
	authService := services.NewAuthService(userRepo, jwtManager)
	userService := services.NewUserService(userRepo, propertyRepo)
	propertyService := services.NewPropertyService(propertyRepo)
	newPropertyService := services.NewNewPropertyService(newPropertyRepo, logger)
	servicedApartmentService := services.NewServicedApartmentService(servicedApartmentRepo, logger)
	estateService := services.NewEstateService(estateRepo, logger)
	valuationService := services.NewValuationService(valuationRepo, logger)
	furnitureService := services.NewFurnitureService(furnitureRepo)
	cartService := services.NewCartService(cartRepo, furnitureRepo)
	mortgageService := services.NewMortgageService(mortgageRepo, logger)
	newsService := services.NewNewsService(newsRepo, logger)
	agentService := services.NewAgentService(agentRepo, logger)
	agencyService := services.NewAgencyService(agencyRepo, logger)
	schoolService := services.NewSchoolService(schoolRepo, propertyRepo, estateRepo, logger)
	priceIndexService := services.NewPriceIndexService(priceIndexRepo, logger)
	facilityService := services.NewFacilityService(facilityRepo, logger)
	searchService := services.NewSearchService(searchRepo, logger)
	statisticsService := services.NewStatisticsService(statisticsRepo, logger)
	configService := services.NewConfigService(db, logger)
	// 初始化控制器层
	baseHandler := controllers.NewBaseHandler()
	authHandler := controllers.NewAuthHandler(authService)
	userHandler := controllers.NewUserHandler(userService)
	propertyHandler := controllers.NewPropertyHandler(propertyService)
	newPropertyHandler := controllers.NewNewPropertyHandler(newPropertyService, logger)
	servicedApartmentHandler := controllers.NewServicedApartmentHandler(servicedApartmentService)
	estateHandler := controllers.NewEstateHandler(estateService)
	valuationHandler := controllers.NewValuationHandler(valuationService)
	furnitureHandler := controllers.NewFurnitureHandler(furnitureService)
	cartHandler := controllers.NewCartHandler(cartService)
	mortgageHandler := controllers.NewMortgageHandler(mortgageService)
	newsHandler := controllers.NewNewsHandler(newsService)
	schoolHandler := controllers.NewSchoolHandler(schoolService)
	agentHandler := controllers.NewAgentHandler(agentService)
	agencyHandler := controllers.NewAgencyHandler(agencyService)
	priceIndexHandler := controllers.NewPriceIndexHandler(priceIndexService)
	facilityHandler := controllers.NewFacilityHandler(facilityService)
	searchHandler := controllers.NewSearchHandler(searchService)
	statisticsHandler := controllers.NewStatisticsHandler(baseHandler, statisticsService)
	configHandler := controllers.NewConfigHandler(baseHandler, configService)
	configHandler := controllers.NewConfigHandler(baseHandler, configService)

	// 设置路由
	r := routes.SetupRouter(baseHandler, authHandler, userHandler, propertyHandler, newPropertyHandler, servicedApartmentHandler, estateHandler, valuationHandler, furnitureHandler, cartHandler, mortgageHandler, newsHandler, schoolHandler, agentHandler, agencyHandler, priceIndexHandler, facilityHandler, searchHandler, statisticsHandler, configHandler, jwtManager, logger)

	// 启动服务器
	serverHost := getEnv("SERVER_HOST", "0.0.0.0")
	serverPort := getEnv("SERVER_PORT", "8080")
	addr := fmt.Sprintf("%s:%s", serverHost, serverPort)
	logger.Info("Starting server", zap.String("address", addr))

	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
func initDatabase(logger *zap.Logger) (*gorm.DB, error) {
	// 从环境变量读取数据库配置
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "ajoliving")
	dbPassword := getEnv("DB_PASSWORD", "secret")
	dbName := getEnv("DB_NAME", "ajoliving_db")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池
	maxIdleConns := getEnvAsInt("DB_MAX_IDLE_CONNS", 10)
	maxOpenConns := getEnvAsInt("DB_MAX_OPEN_CONNS", 100)
	connMaxLifetime := getEnvAsInt("DB_CONN_MAX_LIFETIME", 3600)

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Property{},
		&models.PropertyImage{},
		&models.District{},
		&models.Estate{},
		&models.Facility{},
		&models.Furniture{},
		&models.CartItem{},
		&models.NewProperty{},
		&models.ServicedApartment{},
		&models.Agent{},
		&models.AgentContact{},
		&models.AgencyDetail{},
		&models.MortgageRate{},
		&models.News{},
		&models.School{},
		&models.PriceIndex{},
		&models.SearchHistory{},
	)
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt 获取整型环境变量
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}	&model.MortgageRate{},
		&model.MortgageApplication{},
		&model.NewsCategory{},
		&model.News{},
		&model.SchoolNet{},
		&model.School{},
		&model.PriceIndex{},
		&model.SearchHistory{},
	)
}
