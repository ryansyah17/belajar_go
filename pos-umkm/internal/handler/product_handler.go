package handler

import (
	"strconv"

	"pos-umkm/internal/domain"
	"pos-umkm/pkg/response"
	"pos-umkm/pkg/utils"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productUsecase domain.ProductUsecase
}

func NewProductHandler(productUsecase domain.ProductUsecase) *ProductHandler {
	return &ProductHandler{productUsecase: productUsecase}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	var params domain.ProductFilter
	c.ShouldBindQuery(&params)
	params.SetDefaults()

	products, total, err := h.productUsecase.GetAll(params)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data produk")
		return
	}

	response.SuccessWithPagination(c, "Data produk", products, response.Meta{
		Page:      params.Page,
		Limit:     params.Limit,
		Total:     total,
		TotalPage: domain.CalculateTotalPage(total, params.Limit),
	})
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	product, err := h.productUsecase.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, "Detail produk", product)
}

func (h *ProductHandler) GetLowStock(c *gin.Context) {
	products, err := h.productUsecase.GetLowStock()
	if err != nil {
		response.InternalError(c, "Gagal mengambil data stok")
		return
	}
	response.Success(c, "Produk dengan stok menipis", products)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req domain.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid", err.Error())
		return
	}

	if errs := utils.ValidateStruct(&req); len(errs) > 0 {
		response.BadRequest(c, "Validasi gagal", errs)
		return
	}

	product, err := h.productUsecase.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Created(c, "Produk berhasil dibuat", product)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var req domain.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid", err.Error())
		return
	}

	product, err := h.productUsecase.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Produk berhasil diupdate", product)
}

// UploadImage — endpoint khusus upload gambar produk
// menggunakan multipart/form-data bukan JSON
func (h *ProductHandler) UploadImage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	// Ambil file dari form
	file, err := c.FormFile("image")
	if err != nil {
		response.BadRequest(c, "File gambar tidak ditemukan, gunakan field 'image'", nil)
		return
	}

	// Proses upload
	result, err := utils.UploadImage(file, "products")
	if err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	// Simpan file ke disk
	if err := c.SaveUploadedFile(file, result.FilePath); err != nil {
		response.InternalError(c, "Gagal menyimpan file")
		return
	}

	// Update URL gambar di database
	if err := h.productUsecase.UpdateImage(uint(id), result.URL); err != nil {
		// Kalau gagal simpan ke DB, hapus file yang sudah terupload
		utils.DeleteFile(result.URL)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Gambar berhasil diupload", gin.H{
		"image_url": result.URL,
	})
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "ID tidak valid", nil)
		return
	}

	if err := h.productUsecase.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error(), nil)
		return
	}

	response.Success(c, "Produk berhasil dihapus", nil)
}