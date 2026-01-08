package databases

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/clutchtechnology/hk_ajoliving_app_go/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() error {
	// 从环境变量读取数据库配置
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "ajoliving")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "ajoliving_db")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// 构建 DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Hong_Kong",
		host, user, password, dbname, port, sslmode,
	)

	// GORM 配置
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().In(time.FixedZone("HKT", 8*60*60))
		},
	}

	// 连接数据库
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取底层数据库连接
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("✅ Database connected successfully")

	// 自动迁移（开发环境）
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	return nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate() error {
	log.Println("🔄 Running database auto migration...")

	// 按依赖顺序迁移表：先创建被引用的表，再创建引用它们的表
	err := DB.AutoMigrate(
		// 基础表（无外键依赖）
		&models.User{},
		&models.District{},
		&models.FurnitureCategory{},
		&models.Facility{},
		
		// 有外键依赖的表
		&models.Estate{},
		&models.EstateImage{},
		&models.EstateFacility{},
		&models.Property{},
		&models.PropertyImage{},
		&models.NewProperty{},
		&models.NewPropertyImage{},
		&models.NewPropertyLayout{},
		&models.ServicedApartment{},
		&models.ServicedApartmentUnit{},
		&models.ServicedApartmentImage{},
		&models.Furniture{},
		&models.FurnitureImage{},
		&models.CartItem{},
		&models.SchoolNet{},
		&models.School{},
		&models.Agent{},
		&models.AgentServiceArea{},
		&models.AgentContact{},
		&models.AgencyDetail{},
		&models.AgencyContact{},
		&models.SearchHistory{},
	)

	if err != nil {
		return err
	}

	log.Println("✅ Database auto migration completed")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
