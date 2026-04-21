package domain

import "time"

// Role adalah tipe custom string untuk membatasi nilai yang valid
type Role string

const (
	RoleAdmin Role = "admin"
	RoleManager Role = "manager"
	RoleStaff Role = "staff"
)

type User struct{
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // json:"-" → tidak pernah dikirim ke response
	Role      Role      `gorm:"type:enum('admin','manager','staff');default:'staff'" json:"role"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}