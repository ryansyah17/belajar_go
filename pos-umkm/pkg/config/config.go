package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AppConfig struct {
	Port           string
	Env            string
	JWTSecret      string
	JWTExpiresHour int
	UploadPath     string
	MaxFileSize    int64
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found using system env")
	}
}

func GetAppConfig() AppConfig{
	return  AppConfig{
		Port: getEnv("APP_PORT", "8080"),
		Env: getEnv("APP_ENV", "development"),
		JWTSecret: getEnv("JWT_SECRET", "secret"),
		UploadPath: getEnv("UPLOAD_PATH", "./storage/uploads"),
		MaxFileSize: 5 * 1024 * 1024, // 5MB default
	}
}

func InitDB() *gorm.DB{
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "3306"),
		getEnv("DB_NAME", "pos_umkm"),
	)

	gormLogger := logger.Default.LogMode(logger.Silent)
	if getEnv("APP_ENV", "development") == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Connection pool settings (best practice)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Database connected successfully")
	return db
}

func getEnv(key, defaultValue string) string{
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}