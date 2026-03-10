package order

import (
	"order/app/internal/common"
	"order/app/pkg/di"
	"strconv"
)

type ServiceOrder struct {
	*RepositoryOrder
	*ServiceOrderDep
}
type ServiceOrderDep struct {
	di.IProductRepo
}

func NewServiceOrder(repo *RepositoryOrder, dep *ServiceOrderDep) *ServiceOrder {
	return &ServiceOrder{
		RepositoryOrder: repo,
		ServiceOrderDep: dep,
	}
}
func (s *ServiceOrder) CreateOrder(productsHash []string, idUser int) (*common.Order, error) {
	var incorrectHash string
	var correctProd []common.Product
	for _, hash := range productsHash {
		prod, errHash := s.IProductRepo.GetByAllHash(hash)
		if errHash != nil {
			incorrectHash += hash + " "
		} else {
			correctProd = append(correctProd, *prod)
		}
	}
	if len(correctProd) == 0 {
		return nil, ErrHashes(incorrectHash)
	}
	order := &common.Order{
		Products: correctProd,
		UserId:   idUser,
	}
	errCreate := s.RepositoryOrder.CreateOrder(order)
	if errCreate != nil {
		return nil, ErrCreateOrder
	}
	order.Error = ErrHashes(incorrectHash).Error()
	return order, nil
}
func (s *ServiceOrder) GetOrder(idUser int, idOrder string) (*common.Order, error) {
	id, errId := strconv.Atoi(idOrder)
	if errId != nil {
		return nil, ErrIndexOrder
	}
	order, errGet := s.RepositoryOrder.GetOrder(idUser, uint(id))
	if errGet != nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}
func (s *ServiceOrder) GetAllOrders(idUser int, limitStr, offsetStr string) (*AllOrdersResponse, error) {
	var allOrders AllOrdersResponse
	if limitStr == "" && offsetStr == "" {
		Orders, errGet := s.RepositoryOrder.GetAllOrders(idUser)
		if errGet != nil {
			return nil, ErrOrderNotFound
		}
		allOrders.MyOrders = Orders
		return &allOrders, nil
	}

	limit, errLimit := strconv.Atoi(limitStr)
	offset, errOffset := strconv.Atoi(offsetStr)
	if errLimit != nil || errOffset != nil {
		return nil, ErrParams
	}
	Orders, errOrders := s.RepositoryOrder.GetAllOrdersWithParams(idUser, limit, offset)
	if errOrders != nil {
		return nil, ErrOrderNotFound
	}
	allOrders.MyOrders = Orders
	allOrders.Limit = limit
	allOrders.Offset = offset
	return &allOrders, nil
}
