package interfaces

import "github.com/EQUISEED-WEALTH/pinch/backend/domain/models"

// UserRepo interface, which is used to interact with the user repository
type UserRepo interface {
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(user *models.User) error
	FindOne(user *models.User) (*models.User, error)
	FindMany(user *models.User) ([]*models.User, error)
}

// UserService interface, which is used to interact with the repo and controller
type UserService interface {
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(user *models.User) error
	FindOne(user *models.User) (*models.User, error)
	FindAll() ([]*models.User, error)
}
