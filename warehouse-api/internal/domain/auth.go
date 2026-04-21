package domain

// DTO untuk request — ini yang diterima dari body HTTP
type RegisterRequest struct{
	Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Role     Role   `json:"role" validate:"omitempty,oneof=admin manager staff"`
}

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

// DTO untuk response — data user yang aman dikirim ke client
type UserResponse struct {
    ID        uint   `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    Role      Role   `json:"role"`
    IsActive  bool   `json:"is_active"`
}

type LoginResponse struct {
    Token string       `json:"token"`
    User  UserResponse `json:"user"`
}

// Helper: konversi dari domain.User ke UserResponse (tanpa password)
func ToUserResponse(u *User) UserResponse {
    return UserResponse{
        ID:       u.ID,
        Name:     u.Name,
        Email:    u.Email,
        Role:     u.Role,
        IsActive: u.IsActive,
    }
}