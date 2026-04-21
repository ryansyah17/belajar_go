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

	// Auto migrate
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

	// ── Inisialisasi JWT Util ──────────────────────────────────────
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpire, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES_HOUR"))
	if jwtExpire == 0 {
		jwtExpire = 24
	}
	jwtUtil := utils.NewJWTUtil(jwtSecret, jwtExpire)

	// ── Dependency Injection ───────────────────────────────────────
	// Layer order: Repository → Usecase → Handler

	// Auth
	userRepo    := repository.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepo, jwtUtil)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Handlers
	handlers := &handler.Handlers{
		Auth: handler.NewAuthHandler(authUsecase),
		User: handler.NewUserHandler(userUsecase),
	}

	// ── Setup Gin ─────────────────────────────────────────────────
	cfg := config.GetAppConfig()
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "POS UMKM API running"})
	})

	handler.SetupRouter(r, handlers, jwtUtil)

	log.Printf("Server running on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}