package domain

import (
	"time"

	"gorm.io/gorm"
)

type PaymentMethod string

const (
	PaymentCash   PaymentMethod = "cash"
	PaymentTransfer PaymentMethod = "transfer"
	PaymentQRIS   PaymentMethod = "qris"
)

type Transaction struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserID        uint           `json:"user_id"`
	InvoiceNumber string         `json:"invoice_number" gorm:"uniqueIndex;not null"`
	TotalAmount   float64        `json:"total_amount"`
	DiscountAmount float64       `json:"discount_amount" gorm:"default:0"`
	FinalAmount   float64        `json:"final_amount"`
	PaymentMethod PaymentMethod  `json:"payment_method" gorm:"type:enum('cash','transfer','qris')"`
	PaidAmount    float64        `json:"paid_amount"`
	ChangeAmount  float64        `json:"change_amount"`
	Notes         string         `json:"notes"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	User  User              `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Items []TransactionItem `json:"items" gorm:"foreignKey:TransactionID"`
}

type TransactionItem struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	TransactionID uint    `json:"transaction_id"`
	ProductID     uint    `json:"product_id"`
	ProductName   string  `json:"product_name"` // snapshot nama saat transaksi
	Price         float64 `json:"price"`         // snapshot harga saat transaksi
	Quantity      int     `json:"quantity"`
	Subtotal      float64 `json:"subtotal"`

	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

type TransactionRepository interface {
	FindAll(params TransactionFilter) ([]Transaction, int64, error)
	FindByID(id uint) (*Transaction, error)
	FindByInvoice(invoice string) (*Transaction, error)
	Create(tx *Transaction) error
	SumRevenueByDateRange(start, end time.Time) (float64, error)
	GetDailySummary(date time.Time) (*DailySummary, error)
}

type TransactionUsecase interface {
	GetAll(params TransactionFilter) ([]Transaction, int64, error)
	GetByID(id uint) (*Transaction, error)
	Create(userID uint, req *CreateTransactionRequest) (*Transaction, error)
}

type TransactionFilter struct {
	PaginationParams
	StartDate     time.Time `form:"start_date"`
	EndDate       time.Time `form:"end_date"`
	PaymentMethod string    `form:"payment_method"`
}

type CreateTransactionRequest struct {
	Items          []TransactionItemRequest `json:"items" validate:"required,min=1,dive"`
	DiscountAmount float64                  `json:"discount_amount" validate:"min=0"`
	PaymentMethod  PaymentMethod            `json:"payment_method" validate:"required,oneof=cash transfer qris"`
	PaidAmount     float64                  `json:"paid_amount" validate:"required,gt=0"`
	Notes          string                   `json:"notes"`
}

type TransactionItemRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,min=1"`
}

type DailySummary struct {
	Date         time.Time `json:"date"`
	TotalRevenue float64   `json:"total_revenue"`
	TotalExpense float64   `json:"total_expense"`
	NetProfit    float64   `json:"net_profit"`
	TotalOrders  int64     `json:"total_orders"`
}