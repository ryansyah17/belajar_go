package usecase

import (
	"errors"

	"pos-umkm/internal/domain"
	"pos-umkm/pkg/utils"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) GetAll(params domain.PaginationParams) ([]domain.User, int64, error) {
	params.SetDefaults()
	return u.userRepo.FindAll(params)
}

func (u *userUsecase) GetByID(id uint) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *userUsecase) Create(req *domain.CreateUserRequest) (*domain.User, error) {
	existing, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
		Role:     req.Role,
		IsActive: true,
	}
	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Update(id uint, req *domain.UpdateUserRequest) (*domain.User, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// cek email baru tidak dipakai user lain
		existing, _ := u.userRepo.FindByEmail(req.Email)
		if existing != nil && existing.ID != id {
			return nil, errors.New("email sudah dipakai user lain")
		}
		user.Email = req.Email
	}
	if req.Password != "" {
		hashed, err := utils.HashPassword(req.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashed
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Delete(id uint) error {
	_, err := u.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	return u.userRepo.Delete(id)
}