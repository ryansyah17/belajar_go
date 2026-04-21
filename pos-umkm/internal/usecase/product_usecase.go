package usecase

import (
	"errors"

	"pos-umkm/internal/domain"
	"pos-umkm/pkg/utils"
)

type productUsecase struct {
	productRepo  domain.ProductRepository
	categoryRepo domain.CategoryRepository
}

func NewProductUsecase(
	productRepo domain.ProductRepository,
	categoryRepo domain.CategoryRepository,
) domain.ProductUsecase {
	return &productUsecase{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (u *productUsecase) GetAll(params domain.ProductFilter) ([]domain.Product, int64, error) {
	params.SetDefaults()
	return u.productRepo.FindAll(params)
}

func (u *productUsecase) GetByID(id uint) (*domain.Product, error) {
	return u.productRepo.FindByID(id)
}

func (u *productUsecase) GetLowStock() ([]domain.Product, error) {
	return u.productRepo.FindLowStock()
}

func (u *productUsecase) Create(req *domain.CreateProductRequest) (*domain.Product, error) {
	// Validasi kategori exist
	_, err := u.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}

	product := &domain.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		MinStock:    req.MinStock,
		Unit:        req.Unit,
		IsActive:    true,
	}

	if err := u.productRepo.Create(product); err != nil {
		return nil, errors.New("gagal membuat produk")
	}

	// Reload dengan relasi Category
	return u.productRepo.FindByID(product.ID)
}

func (u *productUsecase) Update(id uint, req *domain.UpdateProductRequest) (*domain.Product, error) {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.CategoryID != 0 {
		_, err := u.categoryRepo.FindByID(req.CategoryID)
		if err != nil {
			return nil, errors.New("kategori tidak ditemukan")
		}
		product.CategoryID = req.CategoryID
	}
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.MinStock > 0 {
		product.MinStock = req.MinStock
	}
	if req.Unit != "" {
		product.Unit = req.Unit
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := u.productRepo.Update(product); err != nil {
		return nil, errors.New("gagal mengupdate produk")
	}

	return u.productRepo.FindByID(product.ID)
}

func (u *productUsecase) UpdateImage(id uint, imagePath string) error {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Hapus gambar lama kalau ada
	if product.ImageURL != "" {
		utils.DeleteFile(product.ImageURL)
	}

	product.ImageURL = imagePath
	return u.productRepo.Update(product)
}

func (u *productUsecase) Delete(id uint) error {
	_, err := u.productRepo.FindByID(id)
	if err != nil {
		return err
	}
	return u.productRepo.Delete(id)
}