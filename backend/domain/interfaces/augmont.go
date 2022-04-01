package interfaces

import (
	"time"

	"github.com/EQUISEED-WEALTH/pinch/backend/domain/models"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/utils"
)

// Augmont User, Bank, Address Table CRUD Interface
type AugmontUserRepo interface {
	CreateUser(*models.AugmontUser) error
	UpdateUser(*models.AugmontUser) error
	FindUser(*models.AugmontUser) (*models.AugmontUser, error)
	FindUsers(*models.AugmontUser) ([]*models.AugmontUser, error)
	FindAllUsers() ([]*models.AugmontUser, error)

	CreateBank(*models.AugmontUserBank) error
	DeleteBank(*models.AugmontUserBank) error
	FindBank(*models.AugmontUserBank) (*models.AugmontUserBank, error)
	FindBanks(*models.AugmontUserBank) ([]*models.AugmontUserBank, error)
	FindAllBanks() ([]*models.AugmontUserBank, error)

	CreateAddress(*models.AugmontUserAddress) error
	DeleteAddress(*models.AugmontUserAddress) error
	FindAddress(*models.AugmontUserAddress) (*models.AugmontUserAddress, error)
	FindAddresses(*models.AugmontUserAddress) ([]*models.AugmontUserAddress, error)
	FindAllAddress() ([]*models.AugmontUserAddress, error)
}

// Augmont Order Interface form Buy, Sell & Redeem
type AugmontOrderRepo interface {
	CreateBuy(*models.AugmontBuyOrder) error
	FindBuy(*models.AugmontBuyOrder) (*models.AugmontBuyOrder, error)
	FindBuys(*models.AugmontBuyOrder) ([]*models.AugmontBuyOrder, error)
	FindAllBuys() ([]*models.AugmontBuyOrder, error)

	CreateSell(*models.AugmontSellOrder) error
	FindSell(*models.AugmontSellOrder) (*models.AugmontSellOrder, error)
	FindSells(*models.AugmontSellOrder) ([]*models.AugmontSellOrder, error)
	FindAllSells() ([]*models.AugmontSellOrder, error)

	CreateRedeem(*models.AugmontRedeemOrder) error
	FindRedeem(*models.AugmontRedeemOrder) (*models.AugmontRedeemOrder, error)
	FindRedeems(*models.AugmontRedeemOrder) ([]*models.AugmontRedeemOrder, error)
	FindAllRedeems() ([]*models.AugmontRedeemOrder, error)
}

// Services offered by Augmont
type AugmontService interface {
	// Create customer account using mobile number and unique Id
	CreateUser(info *utils.AugmontUserInfo, user *models.AugmontUser) error
	GetUserInfo(uniqueID string) (*utils.AugmontUserInfo, error)
	UpdateUser(userInfo *utils.AugmontUserInfo) error

	CreateUserBank(user *models.AugmontUser, bankInfo *utils.AugmontUserBankInfo) error
	GetUserBanks(user *models.AugmontUser) ([]*utils.AugmontUserBankInfo, error)
	UpdateUserBank(user *models.AugmontUser, bankInfo *utils.AugmontUserBankInfo) error
	DeleteUserBank(user *models.AugmontUser, bankInfo *utils.AugmontUserBankInfo) error

	CreateUserAddress(user *models.AugmontUser, addressInfo *utils.AugmontUserAddressInfo) error
	GetUserAddresses(user *models.AugmontUser) ([]*utils.AugmontUserAddressInfo, error)
	DeleteUserAddress(user *models.AugmontUser, addressInfo *utils.AugmontUserAddressInfo) error

	PostUserKyc(
		name, pan, dob string,
		user *models.AugmontUser,
		file *utils.File,
	) (utils.Any, error)

	UpdateUserKycStatus(
		user *models.AugmontUser,
	) error

	Buy(
		user *models.AugmontUser,
		buyInfo *utils.AugmontBugInfo,
	) (utils.Any, error)

	BuyInfo(
		userUniqueID,
		tnxID string,
	) (utils.Any, error)
	BuyList(userUniqueID string) (utils.Any, error)

	Sell(
		*models.AugmontUser,
		*utils.AugmontSellInfo,
	) (utils.Any, error)

	SellInfo(
		userUniqueID,
		tnxID string,
	) (utils.Any, error)

	SellList(userUniqueID string) (utils.Any, error)

	Redeem(
		*models.AugmontUser,
		*utils.AugmontRedeemInfo,
	) (utils.Any, error)

	RedeemInfo(
		userUniqueID,
		tnxID string,
	) (utils.Any, error)

	RedeemList(userUniqueID string) (utils.Any, error)
}

// InMemory Augmont Repo
type AugmontInMemRepo interface {
	// GetToken returns augmont auth token
	// Return empty string without error
	// if token not found or expired
	GetToken() (string, error)

	// Set Token with expiry time
	SetToken(token string, expireAt time.Time) error
}

// Augmont Authentication Service
type AugmontAuthService interface {
	AuthToken() (string, error)
}
