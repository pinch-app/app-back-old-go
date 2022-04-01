package service

import (
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/interfaces"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/models"
)

type userService struct {
	userRepo interfaces.UserRepo
}

// NewUserService creates a new UserService
func NewUserService(userRepo interfaces.UserRepo) interfaces.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Create(user *models.User) error {
	return s.userRepo.Create(user)
}

func (s *userService) Update(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *userService) Delete(user *models.User) error {
	return s.userRepo.Delete(user)
}

func (s *userService) FindOne(user *models.User) (*models.User, error) {
	return s.userRepo.FindOne(user)
}

func (s *userService) FindAll() ([]*models.User, error) {
	return s.userRepo.FindMany(nil)
}
