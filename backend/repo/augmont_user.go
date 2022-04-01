package repo

import (
	"gorm.io/gorm"

	"github.com/EQUISEED-WEALTH/pinch/backend/domain/interfaces"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/models"
)

type augmontUserRepo struct {
	db *gorm.DB
}

// NewAugmontUserRepo returns a new instance of AugmontUserRepo
func NewAugmontUserRepo(db *gorm.DB) interfaces.AugmontUserRepo {

	// Migrate all Augmont related tables
	db.AutoMigrate(
		&models.AugmontUser{},

		&models.AugmontUserBank{},
		&models.AugmontUserAddress{},

		&models.AugmontBuyOrder{},
		&models.AugmontSellOrder{},
		&models.AugmontRedeemOrder{},
	)

	return &augmontUserRepo{
		db: db,
	}
}

// ---- Augmont User ----

func (r *augmontUserRepo) CreateUser(user *models.AugmontUser) error {
	return r.db.Create(user).Error
}

func (r *augmontUserRepo) UpdateUser(user *models.AugmontUser) error {
	return r.db.Model(&models.AugmontUser{ID: user.ID}).Updates(user).Error
}

func (r *augmontUserRepo) FindUser(user *models.AugmontUser) (*models.AugmontUser, error) {
	var userFound models.AugmontUser
	err := r.db.
		Where(user).
		First(&userFound).
		Error
	if err != nil {
		return nil, err
	}
	return &userFound, nil
}

func (r *augmontUserRepo) FindUsers(user *models.AugmontUser) ([]*models.AugmontUser, error) {
	var users []*models.AugmontUser
	err := r.db.
		Where(user).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *augmontUserRepo) FindAllUsers() ([]*models.AugmontUser, error) {
	var users []*models.AugmontUser
	err := r.db.
		Find(&users).
		Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// ---- Augmont User Bank ----

func (r *augmontUserRepo) CreateBank(bank *models.AugmontUserBank) error {
	return r.db.Create(bank).Error
}

func (r *augmontUserRepo) DeleteBank(bank *models.AugmontUserBank) error {
	return r.db.
		Where(bank).
		Limit(1).
		Delete(models.AugmontUserBank{}).
		Error
}

func (r *augmontUserRepo) FindBank(bank *models.AugmontUserBank) (*models.AugmontUserBank, error) {
	var userBank models.AugmontUserBank
	err := r.db.
		Where(bank).
		First(&userBank).
		Error
	if err != nil {
		return nil, err
	}
	return &userBank, nil
}

func (r *augmontUserRepo) FindBanks(bank *models.AugmontUserBank) ([]*models.AugmontUserBank, error) {
	var userBanks []*models.AugmontUserBank
	err := r.db.
		Where(bank).
		Find(&userBanks).
		Error
	if err != nil {
		return nil, err
	}
	return userBanks, nil
}

func (r *augmontUserRepo) FindAllBanks() ([]*models.AugmontUserBank, error) {
	var userBanks []*models.AugmontUserBank
	err := r.db.
		First(&userBanks).
		Error
	if err != nil {
		return nil, err
	}
	return userBanks, nil
}

// ---- Augmont User Address ----

func (r *augmontUserRepo) CreateAddress(address *models.AugmontUserAddress) error {
	return r.db.Create(address).Error
}

func (r *augmontUserRepo) DeleteAddress(address *models.AugmontUserAddress) error {
	return r.db.
		Where(address).
		Limit(1).
		Delete(models.AugmontUserAddress{}).
		Error
}

func (r *augmontUserRepo) FindAddress(address *models.AugmontUserAddress) (*models.AugmontUserAddress, error) {
	var userAddress models.AugmontUserAddress
	err := r.db.
		Where(address).
		Find(&userAddress).
		Error
	if err != nil {
		return nil, err
	}
	return &userAddress, nil
}

func (r *augmontUserRepo) FindAddresses(address *models.AugmontUserAddress) ([]*models.AugmontUserAddress, error) {
	var userAddress []*models.AugmontUserAddress
	err := r.db.
		Where(address).
		Find(&userAddress).
		Error
	if err != nil {
		return nil, err
	}
	return userAddress, nil
}

func (r *augmontUserRepo) FindAllAddress() ([]*models.AugmontUserAddress, error) {
	var userAddress []*models.AugmontUserAddress
	err := r.db.
		Find(&userAddress).
		Error
	if err != nil {
		return nil, err
	}
	return userAddress, nil
}
