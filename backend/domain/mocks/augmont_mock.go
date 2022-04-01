package mocks

import (
	"github.com/stretchr/testify/mock"

	"github.com/EQUISEED-WEALTH/pinch/backend/domain/models"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/utils"
)

type AugmontService struct {
	mock.Mock
}

func NewAugmontService() *AugmontService {
	return &AugmontService{}
}

func (m *AugmontService) CreateUser(info *utils.AugmontUserInfo, user *models.AugmontUser) error {
	args := m.Called(info, user)
	return args.Error(0)
}
func (m *AugmontService) GetUserInfo(user string) (*utils.AugmontUserInfo, error) {
	args := m.Called(user)
	return nil, args.Error(0)
}

func (m *AugmontService) UpdateUser(info *utils.AugmontUserInfo) error {
	args := m.Called(info)
	return args.Error(0)
}
