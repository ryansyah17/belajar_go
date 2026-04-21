package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    AppPort string
    AppEnv  string
    DB      DBConfig
}

type DBConfig struct {
	Host string
	Port string
	User string
	Password string
	Name string
}

// Load membaca .env dan mengembalikan struct Config
func Load() *Config{
	// Hanya load .env jika bukan production
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, reading from environment")
		}
	}

	return &Config{
        AppPort: getEnv("APP_PORT", "8080"),
        AppEnv:  getEnv("APP_ENV", "development"),
        DB: DBConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "3306"),
            User:     getEnv("DB_USER", "root"),
            Password: getEnv("DB_PASSWORD", ""),
            Name:     getEnv("DB_NAME", "warehouse_db"),
        },
    }
}

// getEnv mengambil nilai env, fallback ke defaultVal jika kosong
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// DSN menghasilkan string koneksi MySQL untuk GORM
func (c *Config) DSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        c.DB.User,
        c.DB.Password,
        c.DB.Host,
        c.DB.Port,
        c.DB.Name,
    )
}