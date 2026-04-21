package middleware

import (
	"strings"

	"pos-umkm/pkg/response"
	"pos-umkm/pkg/utils"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserID = "user_id"
	CtxEmail  = "email"
	CtxRole   = "role"
)

// AuthMiddleware — validasi JWT token dari header Authorization
func AuthMiddleware(jwtUtil *utils.JWTUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Token tidak ditemukan")
			c.Abort()
			return
		}

		// Format header: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Format token tidak valid")
			c.Abort()
			return
		}

		claims, err := jwtUtil.ValidateToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "Token tidak valid atau sudah expired")
			c.Abort()
			return
		}

		// Simpan claims ke context — bisa diakses handler dengan c.GetUint("user_id")
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxEmail, claims.Email)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}