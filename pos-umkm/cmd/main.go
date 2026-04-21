package main

import (
	"log"

	"pos-umkm/internal/domain"
	"pos-umkm/pkg/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Initialize database
	db := config.InitDB()

	// Auto migrate semua tabel — GORM akan buat/update tabel otomatis
	err := db.AutoMigrate(
		&domain.User{},
		&domain.Category{},
		&domain.Product{},
		&domain.Expense{},
		&domain.Transaction{},
		&domain.TransactionItem{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Database migration completed")

	// Setup Gin
	cfg := config.GetAppConfig()
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "POS UMKM API is running",
		})
	})

	// TODO: Register routes (akan kita isi di phase berikutnya)

	log.Printf("Server running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}