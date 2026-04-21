package usecase

import (
	"errors"

	"pos-umkm/internal/domain"
)

type categoryUsecase struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryUsecase(categoryRepo domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{categoryRepo: categoryRepo}
}

func (u *categoryUsecase) GetAll(params domain.PaginationParams) ([]domain.Category, int64, error) {
	params.SetDefaults()
	return u.categoryRepo.FindAll(params)
}

func (u *categoryUsecase) GetByID(id uint) (*domain.Category, error) {
	return u.categoryRepo.FindByID(id)
}

func (u *categoryUsecase) Create(req *domain.CreateCategoryRequest) (*domain.Category, error) {
	// Cek nama kategori sudah ada
	existing, err := u.categoryRepo.FindByName(req.Name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("nama kategori sudah digunakan")
	}

	category := &domain.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := u.categoryRepo.Create(category); err != nil {
		return nil, errors.New("gagal membuat kategori")
	}
	return category, nil
}

func (u *categoryUsecase) Update(id uint, req *domain.UpdateCategoryRequest) (*domain.Category, error) {
	category, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" && req.Name != category.Name {
		existing, _ := u.categoryRepo.FindByName(req.Name)
		if existing != nil && existing.ID != id {
			return nil, errors.New("nama kategori sudah digunakan")
		}
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}

	if err := u.categoryRepo.Update(category); err != nil {
		return nil, errors.New("gagal mengupdate kategori")
	}
	return category, nil
}

func (u *categoryUsecase) Delete(id uint) error {
	_, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}
	return u.categoryRepo.Delete(id)
}