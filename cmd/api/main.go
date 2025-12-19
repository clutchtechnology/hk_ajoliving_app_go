package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/config"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/handler"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/model"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/pkg/utils"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/repository"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/router"
	"github.com/clutchtechnology/hk_ajoliving_app_go/internal/service"
)

func main() {
	// 初始化日志
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	db, err := initDatabase(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// 初始化 JWT 管理器
	jwtManager := utils.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.ExpireHours,
		cfg.JWT.RefreshExpireHours,
	)

	// 初始化仓库层
	userRepo := repository.NewUserRepository(db)
	propertyRepo := repository.NewPropertyRepository(db)
	newPropertyRepo := repository.NewNewPropertyRepository(db)
	servicedApartmentRepo := repository.NewServicedApartmentRepository(db)
	estateRepo := repository.NewEstateRepository(db)
	valuationRepo := repository.NewValuationRepository(db)
	furnitureRepo := repository.NewFurnitureRepository(db)
	cartRepo := repository.NewCartRepository(db)
	mortgageRepo := repository.NewMortgageRepository(db)
	newsRepo := repository.NewNewsRepository(db)
	agentRepo := repository.NewAgentRepository(db)
	agencyRepo := repository.NewAgencyRepository(db)
	schoolRepo := repository.NewSchoolRepository(db)
	priceIndexRepo := repository.NewPriceIndexRepository(db)
	facilityRepo := repository.NewFacilityRepository(db)
	searchRepo := repository.NewSearchRepository(db)
	statisticsRepo := repository.NewStatisticsRepository(db)

	// 初始化服务层
	authService := service.NewAuthService(userRepo, jwtManager)
	userService := service.NewUserService(userRepo, propertyRepo)
	propertyService := service.NewPropertyService(propertyRepo)
	newPropertyService := service.NewNewPropertyService(newPropertyRepo, logger)
	servicedApartmentService := service.NewServicedApartmentService(servicedApartmentRepo, logger)
	estateService := service.NewEstateService(estateRepo, logger)
	valuationService := service.NewValuationService(valuationRepo, logger)
	furnitureService := service.NewFurnitureService(furnitureRepo)
	cartService := service.NewCartService(cartRepo, furnitureRepo)
	mortgageService := service.NewMortgageService(mortgageRepo, logger)
	newsService := service.NewNewsService(newsRepo, logger)
	agentService := service.NewAgentService(agentRepo, logger)
	agencyService := service.NewAgencyService(agencyRepo, logger)
	schoolService := service.NewSchoolService(schoolRepo, propertyRepo, estateRepo, logger)
	priceIndexService := service.NewPriceIndexService(priceIndexRepo, logger)
	facilityService := service.NewFacilityService(facilityRepo, logger)
	searchService := service.NewSearchService(searchRepo, logger)
	statisticsService := service.NewStatisticsService(statisticsRepo, logger)
	configService := service.NewConfigService(db, logger)

	// 初始化处理器层
	baseHandler := handler.NewBaseHandler()
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	propertyHandler := handler.NewPropertyHandler(propertyService)
	newPropertyHandler := handler.NewNewPropertyHandler(newPropertyService, logger)
	servicedApartmentHandler := handler.NewServicedApartmentHandler(servicedApartmentService)
	estateHandler := handler.NewEstateHandler(estateService)
	valuationHandler := handler.NewValuationHandler(valuationService)
	furnitureHandler := handler.NewFurnitureHandler(furnitureService)
	cartHandler := handler.NewCartHandler(cartService)
	mortgageHandler := handler.NewMortgageHandler(mortgageService)
	newsHandler := handler.NewNewsHandler(newsService)
	schoolHandler := handler.NewSchoolHandler(schoolService)
	agentHandler := handler.NewAgentHandler(agentService)
	agencyHandler := handler.NewAgencyHandler(agencyService)
	priceIndexHandler := handler.NewPriceIndexHandler(priceIndexService)
	facilityHandler := handler.NewFacilityHandler(facilityService)
	searchHandler := handler.NewSearchHandler(searchService)
	statisticsHandler := handler.NewStatisticsHandler(baseHandler, statisticsService)
	configHandler := handler.NewConfigHandler(baseHandler, configService)

	// 设置路由
	r := router.SetupRouter(baseHandler, authHandler, userHandler, propertyHandler, newPropertyHandler, servicedApartmentHandler, estateHandler, valuationHandler, furnitureHandler, cartHandler, mortgageHandler, newsHandler, schoolHandler, agentHandler, agencyHandler, priceIndexHandler, facilityHandler, searchHandler, statisticsHandler, configHandler, jwtManager, logger)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Starting server", zap.String("address", addr))

	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func initDatabase(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
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
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// 自动迁移
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	logger.Info("Database connected successfully")
	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Property{},
		&model.PropertyImage{},
		&model.District{},
		&model.Estate{},
		&model.EstateImage{},
		&model.Facility{},
		&model.Furniture{},
		&model.FurnitureCategory{},
		&model.FurnitureImage{},
		&model.CartItem{},
		&model.NewProperty{},
		&model.NewPropertyImage{},
		&model.NewPropertyLayout{},
		&model.ServicedApartment{},
		&model.ServicedApartmentUnit{},
		&model.ServicedApartmentImage{},
		&model.Agent{},
		&model.AgentServiceArea{},
		&model.AgentContactRequest{},
		&model.AgencyDetail{},
		&model.Bank{},
		&model.MortgageRate{},
		&model.MortgageApplication{},
		&model.NewsCategory{},
		&model.News{},
		&model.SchoolNet{},
		&model.School{},
		&model.PriceIndex{},
		&model.SearchHistory{},
	)
}
