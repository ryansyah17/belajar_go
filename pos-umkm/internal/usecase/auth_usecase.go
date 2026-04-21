package usecase

import (
	"errors"
	"os"
	"strconv"

	"pos-umkm/internal/domain"
	"pos-umkm/pkg/utils"
)

type authUsecase struct {
	userRepo domain.UserRepository
	jwtUtil  *utils.JWTUtil
}

func NewAuthUsecase(userRepo domain.UserRepository, jwtUtil *utils.JWTUtil) domain.AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

func (u *authUsecase) Register(req *domain.RegisterRequest) (*domain.User, error) {
	// Cek apakah email sudah dipakai
	existing, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
		IsActive: true,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, errors.New("gagal membuat user")
	}

	return user, nil
}

func (u *authUsecase) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// Cari user by email
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("email atau password salah")
	}

	// Cek apakah user aktif
	if !user.IsActive {
		return nil, errors.New("akun tidak aktif, hubungi administrator")
	}

	// Verifikasi password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("email atau password salah")
	}

	// Generate JWT
	expireHour, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES_HOUR"))
	if expireHour == 0 {
		expireHour = 24
	}

	token, err := u.jwtUtil.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	return &domain.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}