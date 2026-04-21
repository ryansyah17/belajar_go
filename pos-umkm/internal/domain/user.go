package domain

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleCashier Role = "cashier"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null;size:100"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:150"`
	Password  string         `json:"-" gorm:"not null;size:255"`
	Role      Role           `json:"role" gorm:"type:enum('owner','admin','cashier');default:'cashier'"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserRepository — interface yang WAJIB diimplementasi repository layer
type UserRepository interface {
	FindAll(params PaginationParams) ([]User, int64, error)
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

// AuthUsecase — interface untuk auth business logic
type AuthUsecase interface {
	Register(req *RegisterRequest) (*User, error)
	Login(req *LoginRequest) (*LoginResponse, error)
}

// UserUsecase — interface untuk user management
type UserUsecase interface {
	GetAll(params PaginationParams) ([]User, int64, error)
	GetByID(id uint) (*User, error)
	Create(req *CreateUserRequest) (*User, error)
	Update(id uint, req *UpdateUserRequest) (*User, error)
	Delete(id uint) error
}

// --- Request & Response DTOs ---

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     Role   `json:"role" validate:"required,oneof=owner admin cashier"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     Role   `json:"role" validate:"required,oneof=owner admin cashier"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=100"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"omitempty,min=6"`
	Role     Role   `json:"role" validate:"omitempty,oneof=owner admin cashier"`
	IsActive *bool  `json:"is_active"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}