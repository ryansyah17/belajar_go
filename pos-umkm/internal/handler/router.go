package handler

import (
	"pos-umkm/pkg/middleware"
	"pos-umkm/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth *AuthHandler
	User *UserHandler
	// Phase berikutnya: Category, Product, dst
}

func SetupRouter(r *gin.Engine, h *Handlers, jwtUtil *utils.JWTUtil) {
	// Static files — untuk akses gambar upload
	r.Static("/storage", "./storage")

	api := r.Group("/api/v1")

	// ── Public routes (tidak perlu token) ──────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
	}

	// ── Protected routes (wajib token) ─────────────────────────────
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtUtil))
	{
		// Auth
		protected.GET("/auth/me", h.Auth.Me)

		// User management — hanya owner
		users := protected.Group("/users")
		users.Use(middleware.RoleMiddleware("owner"))
		{
			users.GET("", h.User.GetAll)
			users.GET("/:id", h.User.GetByID)
			users.POST("", h.User.Create)
			users.PUT("/:id", h.User.Update)
			users.DELETE("/:id", h.User.Delete)
		}

		// Phase berikutnya: tambah category, product, dst di sini
	}
}