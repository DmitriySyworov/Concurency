package order

import (
	"order/app/internal/common"
	"order/app/pkg/db"
)

type RepositoryOrder struct {
	*db.Db
}

func NewRepositoryOrder(db *db.Db) *RepositoryOrder {
	return &RepositoryOrder{
		Db: db,
	}
}
func (repo *RepositoryOrder) CreateOrder(order *common.Order) error {
	res := repo.Db.Create(&order)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (repo *RepositoryOrder) GetOrder(idUser int, id uint) (*common.Order, error) {
	var order common.Order
	res := repo.Db.
		Model(&common.Order{}).
		Preload("Products").
		Where("id = ? AND user_id = ?", id, idUser).
		First(&order)
	if res.Error != nil {
		return nil, res.Error
	}
	return &order, nil
}
func (repo *RepositoryOrder) GetAllOrders(idUser int) ([]common.Order, error) {
	var orders []common.Order
	res := repo.Db.
		Model(&common.Order{}).
		Preload("Products").
		Where("user_id = ?", idUser).
		Find(&orders)
	if res.Error != nil || len(orders) == 0 {
		return nil, ErrOrderNotFound
	}
	return orders, nil
}
func (repo *RepositoryOrder) GetAllOrdersWithParams(idUser, limit, offset int) ([]common.Order, error) {
	var orders []common.Order
	res := repo.Db.
		Model(&common.Order{}).
		Preload("Products").
		Where("user_id = ?", idUser).
		Limit(limit).
		Offset(offset).
		Find(&orders)
	if res.Error != nil || len(orders) == 0 {
		return nil, ErrOrderNotFound
	}
	return orders, nil
}
