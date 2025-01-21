package service

import (
	"errors"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(r *repository.UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) CreateUser(username, password, role string) (int, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	u := repository.User{
		Username:     username,
		PasswordHash: string(h),
		Role:         role,
	}
	id, err := s.repo.Create(u)
	return id, err
}

func (s *UserService) CheckPassword(username, password string) (repository.User, error) {
	u, err := s.repo.GetByUsername(username)
	if err != nil {
		return repository.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return repository.User{}, errors.New("invalid password")
	}
	return u, nil
}