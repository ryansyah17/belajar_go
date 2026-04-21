package repository

import (
	"errors"

	"pos-umkm/internal/domain"

	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll(params domain.PaginationParams) ([]domain.Category, int64, error) {
	var categories []domain.Category
	var total int64

	query := r.db.Model(&domain.Category{})
	query.Count(&total)

	err := query.
		Offset(params.GetOffset()).
		Limit(params.Limit).
		Order("name ASC").
		Find(&categories).Error

	return categories, total, err
}

func (r *categoryRepository) FindByID(id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("kategori tidak ditemukan")
	}
	return &category, err
}

func (r *categoryRepository) FindByName(name string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Where("name = ?", name).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}

func (r *categoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *domain.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uint) error {
	// Cek apakah masih ada produk di kategori ini
	var count int64
	r.db.Model(&domain.Product{}).Where("category_id = ? AND deleted_at IS NULL", id).Count(&count)
	if count > 0 {
		return errors.New("kategori tidak bisa dihapus karena masih memiliki produk")
	}
	return r.db.Delete(&domain.Category{}, id).Error
}