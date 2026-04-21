package middleware

import (
	"pos-umkm/pkg/response"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware — cek apakah role user ada di daftar role yang diizinkan
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(CtxRole)
		if !exists {
			response.Forbidden(c, "Role tidak ditemukan")
			c.Abort()
			return
		}

		userRole := role.(string)
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Anda tidak memiliki akses ke resource ini")
		c.Abort()
	}
}