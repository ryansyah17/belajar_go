package domain

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;uniqueIndex;size:100"`
	Description string         `json:"description" gorm:"size:500"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

type CategoryRepository interface {
	FindAll(params PaginationParams) ([]Category, int64, error)
	FindByID(id uint) (*Category, error)
	FindByName(name string) (*Category, error)
	Create(category *Category) error
	Update(category *Category) error
	Delete(id uint) error
}

type CategoryUsecase interface {
	GetAll(params PaginationParams) ([]Category, int64, error)
	GetByID(id uint) (*Category, error)
	Create(req *CreateCategoryRequest) (*Category, error)
	Update(id uint, req *UpdateCategoryRequest) (*Category, error)
	Delete(id uint) error
}

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"max=500"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=100"`
	Description string `json:"description" validate:"omitempty,max=500"`
}