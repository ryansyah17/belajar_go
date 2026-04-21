package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
    engine *gin.Engine
    db     *gorm.DB
}

func New(db *gorm.DB) *Router {
    r := gin.Default()

    // Global middleware bisa ditambahkan di sini
    // r.Use(middleware.CORS())
    // r.Use(middleware.RateLimiter())

    return &Router{engine: r, db: db}
}

func (r *Router) Setup() {
    api := r.engine.Group("/api/v1")

    // Health check
    api.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "success": true,
            "message": "Warehouse API is running",
            "version": "1.0.0",
        })
    })

    // Auth routes akan didaftarkan di sini (Step 4)
    // r.setupAuthRoutes(api)

    // Product routes (Step 7)
    // r.setupProductRoutes(api)
}

func (r *Router) Run(port string) error {
    return r.engine.Run(":" + port)
}