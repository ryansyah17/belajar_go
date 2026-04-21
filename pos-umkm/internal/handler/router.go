package handler

import (
	"pos-umkm/pkg/middleware"
	"pos-umkm/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Auth     *AuthHandler
	User     *UserHandler
	Category *CategoryHandler
	Product  *ProductHandler
}

func SetupRouter(r *gin.Engine, h *Handlers, jwtUtil *utils.JWTUtil) {
	r.Static("/storage", "./storage")

	api := r.Group("/api/v1")

	// ── Public ────────────────────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Auth.Register)
		auth.POST("/login", h.Auth.Login)
	}

	// ── Protected ─────────────────────────────────────────────────
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtUtil))
	{
		protected.GET("/auth/me", h.Auth.Me)

		// Users — owner only
		users := protected.Group("/users")
		users.Use(middleware.RoleMiddleware("owner"))
		{
			users.GET("", h.User.GetAll)
			users.GET("/:id", h.User.GetByID)
			users.POST("", h.User.Create)
			users.PUT("/:id", h.User.Update)
			users.DELETE("/:id", h.User.Delete)
		}

		// Categories — owner & admin bisa kelola, cashier read-only
		categories := protected.Group("/categories")
		{
			categories.GET("", h.Category.GetAll)
			categories.GET("/:id", h.Category.GetByID)
			// Write operations hanya owner & admin
			categories.POST("", middleware.RoleMiddleware("owner", "admin"), h.Category.Create)
			categories.PUT("/:id", middleware.RoleMiddleware("owner", "admin"), h.Category.Update)
			categories.DELETE("/:id", middleware.RoleMiddleware("owner", "admin"), h.Category.Delete)
		}

		// Products — owner & admin kelola, cashier read-only
		products := protected.Group("/products")
		{
			products.GET("", h.Product.GetAll)
			products.GET("/low-stock", h.Product.GetLowStock)
			products.GET("/:id", h.Product.GetByID)
			products.POST("", middleware.RoleMiddleware("owner", "admin"), h.Product.Create)
			products.PUT("/:id", middleware.RoleMiddleware("owner", "admin"), h.Product.Update)
			products.POST("/:id/image", middleware.RoleMiddleware("owner", "admin"), h.Product.UploadImage)
			products.DELETE("/:id", middleware.RoleMiddleware("owner", "admin"), h.Product.Delete)
		}
	}
}