package main

import (
	"log"
	"os"
	"strconv"

	"pos-umkm/internal/domain"
	"pos-umkm/internal/handler"
	"pos-umkm/internal/repository"
	"pos-umkm/internal/usecase"
	"pos-umkm/pkg/config"
	"pos-umkm/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	db := config.InitDB()

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

	// JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpire, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES_HOUR"))
	if jwtExpire == 0 {
		jwtExpire = 24
	}
	jwtUtil := utils.NewJWTUtil(jwtSecret, jwtExpire)

	// ── Repositories ──────────────────────────────────────────────
	userRepo     := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo  := repository.NewProductRepository(db)

	// ── Usecases ──────────────────────────────────────────────────
	authUsecase     := usecase.NewAuthUsecase(userRepo, jwtUtil)
	userUsecase     := usecase.NewUserUsecase(userRepo)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	productUsecase  := usecase.NewProductUsecase(productRepo, categoryRepo)

	// ── Handlers ──────────────────────────────────────────────────
	handlers := &handler.Handlers{
		Auth:     handler.NewAuthHandler(authUsecase),
		User:     handler.NewUserHandler(userUsecase),
		Category: handler.NewCategoryHandler(categoryUsecase),
		Product:  handler.NewProductHandler(productUsecase),
	}

	// ── Gin ───────────────────────────────────────────────────────
	cfg := config.GetAppConfig()
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	handler.SetupRouter(r, handlers, jwtUtil)

	log.Printf("Server running on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}