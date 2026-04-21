package handler

import (
	"pos-umkm/internal/domain"
	"pos-umkm/pkg/response"
	"pos-umkm/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid", err.Error())
		return
	}

	// Validasi input
	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		response.BadRequest(c, "Validasi gagal", errs)
		return
	}

	user, err := h.authUsecase.Register(&req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Registrasi berhasil", user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid", err.Error())
		return
	}

	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		response.BadRequest(c, "Validasi gagal", errs)
		return
	}

	result, err := h.authUsecase.Login(&req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Login berhasil", result)
}

func (h *AuthHandler) Me(c *gin.Context) {
	// Ambil data dari context yang sudah di-inject middleware
	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	response.Success(c, "Data user", gin.H{
		"user_id": userID,
		"email":   email,
		"role":    role,
	})
}