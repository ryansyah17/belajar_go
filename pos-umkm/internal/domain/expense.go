package domain

import (
	"time"

	"gorm.io/gorm"
)

type Expense struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Category    string         `json:"category" gorm:"not null"` // operasional, bahan baku, dll
	Amount      float64        `json:"amount" gorm:"not null"`
	Description string         `json:"description"`
	Date        time.Time      `json:"date" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type ExpenseRepository interface {
	FindAll(params ExpenseFilter) ([]Expense, int64, error)
	FindByID(id uint) (*Expense, error)
	SumByDateRange(start, end time.Time) (float64, error)
	Create(expense *Expense) error
	Update(expense *Expense) error
	Delete(id uint) error
}

type ExpenseUsecase interface {
	GetAll(params ExpenseFilter) ([]Expense, int64, error)
	GetByID(id uint) (*Expense, error)
	Create(userID uint, req *CreateExpenseRequest) (*Expense, error)
	Update(id uint, req *UpdateExpenseRequest) (*Expense, error)
	Delete(id uint) error
}

type ExpenseFilter struct {
	PaginationParams
	Category  string    `form:"category"`
	StartDate time.Time `form:"start_date"`
	EndDate   time.Time `form:"end_date"`
}

type CreateExpenseRequest struct {
	Name        string    `json:"name" validate:"required"`
	Category    string    `json:"category" validate:"required"`
	Amount      float64   `json:"amount" validate:"required,gt=0"`
	Description string    `json:"description"`
	Date        time.Time `json:"date" validate:"required"`
}

type UpdateExpenseRequest struct {
	Name        string    `json:"name" validate:"omitempty"`
	Category    string    `json:"category" validate:"omitempty"`
	Amount      float64   `json:"amount" validate:"omitempty,gt=0"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}