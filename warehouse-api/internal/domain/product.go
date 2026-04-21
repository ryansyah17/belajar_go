package domain

import "time"

type Product struct{
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Name        string    `gorm:"type:varchar(150);not null" json:"name"`
    SKU         string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"sku"`
    Description string    `gorm:"type:text" json:"description"`
    Stock       int       `gorm:"default:0" json:"stock"`
    MinStock    int       `gorm:"default:0" json:"min_stock"` // alert jika stok di bawah ini
    Unit        string    `gorm:"type:varchar(20)" json:"unit"` // pcs, kg, liter, dll
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}