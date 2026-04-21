package domain

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CategoryID  uint           `json:"category_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	Stock       int            `json:"stock" gorm:"default:0"`
	MinStock    int            `json:"min_stock" gorm:"default:5"` // untuk low stock alert
	Unit        string         `json:"unit" gorm:"default:'pcs'"`
	ImageURL    string         `json:"image_url"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Category Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

type ProductRepository interface {
	FindAll(params ProductFilter) ([]Product, int64, error)
	FindByID(id uint) (*Product, error)
	FindLowStock() ([]Product, error)
	Create(product *Product) error
	Update(product *Product) error
	UpdateStock(id uint, quantity int) error
	Delete(id uint) error
}

type ProductUsecase interface {
	GetAll(params ProductFilter) ([]Product, int64, error)
	GetByID(id uint) (*Product, error)
	GetLowStock() ([]Product, error)
	Create(req *CreateProductRequest) (*Product, error)
	Update(id uint, req *UpdateProductRequest) (*Product, error)
	UpdateImage(id uint, imagePath string) error
	Delete(id uint) error
}

type ProductFilter struct {
	PaginationParams
	CategoryID uint   `form:"category_id"`
	Search     string `form:"search"`
	IsActive   *bool  `form:"is_active"`
}

type CreateProductRequest struct {
	CategoryID  uint    `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required,min=2,max=200"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"min=0"`
	MinStock    int     `json:"min_stock" validate:"min=0"`
	Unit        string  `json:"unit" validate:"required"`
}

type UpdateProductRequest struct {
	CategoryID  uint    `json:"category_id"`
	Name        string  `json:"name" validate:"omitempty,min=2,max=200"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"omitempty,gt=0"`
	MinStock    int     `json:"min_stock" validate:"omitempty,min=0"`
	Unit        string  `json:"unit"`
	IsActive    *bool   `json:"is_active"`
}