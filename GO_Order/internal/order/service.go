package order

import (
	"order/app/internal/common"
	custerrors "order/app/pkg/custErrors"
	"order/app/pkg/di"
	"strconv"

	"github.com/google/uuid"
)

type ServiceOrder struct {
	*RepositoryOrder
	*ServiceOrderDep
}
type ServiceOrderDep struct {
	di.IProductRepo
	di.IUserRepo
}

func NewServiceOrder(repo *RepositoryOrder, dep *ServiceOrderDep) *ServiceOrder {
	return &ServiceOrder{
		RepositoryOrder: repo,
		ServiceOrderDep: dep,
	}
}
func (s *ServiceOrder) CreateOrder(productsHash []string, idUser int) (*common.Order, error) {
	_, errId := s.ServiceOrderDep.GetByIdUser(idUser)
	if errId != nil {
		return nil, custerrors.ErrUserDontExist
	}
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
	id := uuid.New()
	order := &common.Order{
		Products: correctProd,
		UserId:   idUser,
		OrderId:  id.String(),
	}
	errCreate := s.RepositoryOrder.CreateOrder(order)
	if errCreate != nil {
		return nil, ErrCreateOrder
	}
	order.Error = ErrHashes(incorrectHash).Error()
	return order, nil
}
func (s *ServiceOrder) GetOrder(idUser int, idOrder string) (*common.Order, error) {
	order, errGet := s.RepositoryOrder.GetOrder(idOrder, idUser)
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
