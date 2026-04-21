package handler

import (
	"strconv"

	"pos-umkm/internal/domain"
	"pos-umkm/pkg/response"
	"pos-umkm/pkg/utils"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryUsecase domain.CategoryUsecase
}

func NewCategoryHandler(categoryUsecase domain.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{categoryUsecase: categoryUsecase}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	var params domain.PaginationParams
	c.ShouldBindQuery(&params)
	params.SetDefaults()

	categories, total, err := h.categoryUsecase.GetAll(params)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data kategori")
		return
	}

	response.SuccessWithPagination(c, "Data kategori", categories, response.Meta{
		Page:      params.Page,
		Limit:     params.Limit,
		Total:     total,
		TotalPage: domain.CalculateTotalPage(total, params.Limit),
	})
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	category, err := h.categoryUsecase.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Detail kategori", category)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req domain.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid", err.Error())
		return
	}

	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		response.BadRequest(c, "Validasi gagal", errs)
		return
	}

	category, err := h.categoryUsecase.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Kategori berhasil dibuat", category)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var req domain.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid", err.Error())
		return
	}

	category, err := h.categoryUsecase.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Kategori berhasil diupdate", category)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	if err := h.categoryUsecase.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Kategori berhasil dihapus", nil)
}