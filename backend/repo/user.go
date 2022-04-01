package repo

import (
	"gorm.io/gorm"

	"github.com/EQUISEED-WEALTH/pinch/backend/domain/interfaces"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/models"
)

type userRepo struct {
	db *gorm.DB
}

// NewUserRepo creates a new UserRepo
func NewUserRepo(db *gorm.DB) interfaces.UserRepo {
	// Migrate User Model
	db.AutoMigrate(&models.User{})

	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) Update(user *models.User) error {
	return r.db.
		Where(models.User{
			ID: user.ID,
		}).
		Updates(user).Error
}

func (r *userRepo) Delete(user *models.User) error {
	return r.db.Where(user).Delete(&models.User{}).Error
}

func (r *userRepo) FindOne(user *models.User) (*models.User, error) {
	var u models.User
	err := r.db.Where(user).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil

}

func (r *userRepo) FindMany(user *models.User) ([]*models.User, error) {
	var (
		users []*models.User
		err   error
	)

	if user != nil {
		err = r.db.
			Where(user).
			Find(&users).Error
	} else {
		err = r.db.
			Find(&users).Error
	}

	if err != nil {
		return nil, err
	}
	return users, nil
}
