package handler

import (
	"strconv"

	"pos-umkm/internal/domain"
	"pos-umkm/pkg/response"
	"pos-umkm/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(userUsecase domain.UserUsecase) *UserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	var params domain.PaginationParams
	c.ShouldBindQuery(&params)
	params.SetDefaults()

	users, total, err := h.userUsecase.GetAll(params)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data user")
		return
	}

	response.SuccessWithPagination(c, "Data user", users, response.Meta{
		Page:      params.Page,
		Limit:     params.Limit,
		Total:     total,
		TotalPage: domain.CalculateTotalPage(total, params.Limit),
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	user, err := h.userUsecase.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Detail user", user)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid", err.Error())
		return
	}

	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		response.BadRequest(c, "Validasi gagal", errs)
		return
	}

	user, err := h.userUsecase.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "User berhasil dibuat", user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid", err.Error())
		return
	}

	user, err := h.userUsecase.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "User berhasil diupdate", user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	if err := h.userUsecase.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "User berhasil dihapus", nil)
}