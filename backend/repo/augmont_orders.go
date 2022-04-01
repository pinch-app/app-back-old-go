package repo

import (
	"gorm.io/gorm"

	"github.com/EQUISEED-WEALTH/pinch/backend/domain/interfaces"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/models"
)

type augmontOrdersRepo struct {
	db *gorm.DB
}

// NewAugmontOrdersRepo creates a new Augmont orders repo
func NewAugmontOrderRepo(db *gorm.DB) interfaces.AugmontOrderRepo {
	return &augmontOrdersRepo{
		db: db,
	}
}

// ---- BuyOrders Repo ----

func (r augmontOrdersRepo) CreateBuy(order *models.AugmontBuyOrder) error {
	return r.db.Create(order).Error
}

func (r *augmontOrdersRepo) FindBuy(order *models.AugmontBuyOrder) (*models.AugmontBuyOrder, error) {
	var newOrder models.AugmontBuyOrder
	err := r.db.
		Where(order).
		First(&newOrder).
		Error
	if err != nil {
		return nil, err
	}
	return &newOrder, err
}
func (r *augmontOrdersRepo) FindBuys(order *models.AugmontBuyOrder) ([]*models.AugmontBuyOrder, error) {
	var orders []*models.AugmontBuyOrder
	err := r.db.
		Where(order).
		Find(orders).
		Error
	if err != nil {
		return nil, err
	}
	return orders, err
}
func (r *augmontOrdersRepo) FindAllBuys() ([]*models.AugmontBuyOrder, error) {
	var orders []*models.AugmontBuyOrder
	err := r.db.
		Find(orders).
		Error
	if err != nil {
		return nil, err
	}
	return orders, err
}

func (r *augmontOrdersRepo) CreateSell(order *models.AugmontSellOrder) error {
	return r.db.Create(order).Error
}

func (r *augmontOrdersRepo) FindSell(order *models.AugmontSellOrder) (*models.AugmontSellOrder, error) {
	var newOrder models.AugmontSellOrder
	err := r.db.
		Where(order).
		First(&newOrder).
		Error
	if err != nil {
		return nil, err
	}
	return &newOrder, err
}

func (r *augmontOrdersRepo) FindSells(order *models.AugmontSellOrder) ([]*models.AugmontSellOrder, error) {
	var orders []*models.AugmontSellOrder
	err := r.db.
		Where(order).
		Find(orders).
		Error
	if err != nil {
		return nil, err
	}
	return orders, err
}

func (r *augmontOrdersRepo) FindAllSells() ([]*models.AugmontSellOrder, error) {
	var orders []*models.AugmontSellOrder
	err := r.db.
		Find(orders).
		Error
	if err != nil {
		return nil, err
	}
	return orders, err
}

func (r *augmontOrdersRepo) CreateRedeem(order *models.AugmontRedeemOrder) error {
	return r.db.Create(order).Error
}

func (r *augmontOrdersRepo) FindRedeem(order *models.AugmontRedeemOrder) (*models.AugmontRedeemOrder, error) {
	var newOrder models.AugmontRedeemOrder
	err := r.db.
		Where(order).
		First(&newOrder).
		Error
	if err != nil {
		return nil, err
	}
	return &newOrder, err
}

func (r *augmontOrdersRepo) FindRedeems(order *models.AugmontRedeemOrder) ([]*models.AugmontRedeemOrder, error) {
	var orders []*models.AugmontRedeemOrder
	err := r.db.
		Where(order).
		Find(orders).
		Error
	if err != nil {
		return nil, err
	}
	return orders, err
}

func (r *augmontOrdersRepo) FindAllRedeems() ([]*models.AugmontRedeemOrder, error) {
	var orders []*models.AugmontRedeemOrder
	err := r.db.
		Find(orders).
		Error
	if err != nil {
		return nil, err
	}
	return orders, err
}
