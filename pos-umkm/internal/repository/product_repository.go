package repository

import (
	"errors"

	"pos-umkm/internal/domain"

	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(params domain.ProductFilter) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	query := r.db.Model(&domain.Product{}).Preload("Category")

	// Filter dinamis
	if params.CategoryID != 0 {
		query = query.Where("category_id = ?", params.CategoryID)
	}
	if params.Search != "" {
		query = query.Where("name LIKE ?", "%"+params.Search+"%")
	}
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	query.Count(&total)

	err := query.
		Offset(params.GetOffset()).
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) FindByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Category").First(&product, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("produk tidak ditemukan")
	}
	return &product, err
}

func (r *productRepository) FindLowStock() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Preload("Category").
		Where("stock <= min_stock AND is_active = true").
		Order("stock ASC").
		Find(&products).Error
	return products, err
}

func (r *productRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) UpdateStock(id uint, quantity int) error {
	// Atomic update — aman dari race condition
	result := r.db.Model(&domain.Product{}).
		Where("id = ? AND stock + ? >= 0", id, quantity).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity))

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("stok tidak mencukupi")
	}
	return nil
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Product{}, id).Error
}