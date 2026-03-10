package order

import (
	"errors"
	"net/http"
	"order/app/configs"
	"order/app/internal/common"
	custerrors "order/app/pkg/custErrors"
	"order/app/pkg/middleware"
	"order/app/pkg/request"
	"order/app/pkg/response"
)

type HandlerOrder struct {
	*HandlerOrderDep
	common.Order
	allOrders AllOrdersResponse
}
type HandlerOrderDep struct {
	*configs.Config
	Service *ServiceOrder
}

func NewHandlerOrder(router *http.ServeMux, dep *HandlerOrderDep) {
	order := &HandlerOrder{
		HandlerOrderDep: dep,
	}
	router.Handle("POST /order", middleware.IsAuthID(order.CreateOrder(), dep.Config))
	router.Handle("GET /order/{id}", middleware.IsAuthID(order.GetOrder(), dep.Config))
	router.Handle("GET /my-orders", middleware.IsAuthID(order.GetAllOrders(), dep.Config))
}
func (ho *HandlerOrder) CreateOrder() http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		userId, ok := r.Context().Value(middleware.KeyIDUser).(float64)
		if !ok {
			ho.Order.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(writer, ho.Order, http.StatusUnauthorized)
			return
		}
		body, errBody := request.RequestHandler[CreateOrderRequest](r)
		if errBody != nil || len(body.ProductsHash) == 0 {
			ho.Order.Error = errBody.Error()
			response.RespJs(writer, ho.Order, http.StatusBadRequest)
		}
		order, errOrder := ho.Service.CreateOrder(body.ProductsHash, int(userId))
		if errOrder != nil {
			ho.Order.Error = errOrder.Error()
			if errors.Is(errOrder, ErrCreateOrder) {
				response.RespJs(writer, ho.Order, http.StatusInternalServerError)
			} else {
				response.RespJs(writer, ho.Order, http.StatusBadRequest)
			}
			return
		}
		ho.Order.Error = ""
		response.RespJs(writer, order, http.StatusCreated)
	}
}
func (ho *HandlerOrder) GetOrder() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idUser, ok := request.Context().Value(middleware.KeyIDUser).(float64)

		if !ok {
			ho.Order.Error = custerrors.ErrResponse.Error()
			response.RespJs(writer, ho.Order, http.StatusUnauthorized)
			return
		}
		idStr := request.PathValue("id")
		order, errGet := ho.Service.GetOrder(int(idUser), idStr)
		if errGet != nil {
			ho.Order.Error = errGet.Error()
			if errors.Is(errGet, ErrIndexOrder) {
				response.RespJs(writer, ho.Order, http.StatusBadRequest)
			} else {
				response.RespJs(writer, ho.Order, http.StatusNotFound)
			}
			return
		}
		response.RespJs(writer, order, http.StatusOK)
	}
}
func (ho *HandlerOrder) GetAllOrders() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idUser, ok := request.Context().Value(middleware.KeyIDUser).(float64)
		if !ok {
			ho.Order.Error = custerrors.ErrResponse.Error()
			response.RespJs(writer, ho.allOrders, http.StatusUnauthorized)
			return
		}
		limitStr := request.URL.Query().Get("limit")
		offsetStr := request.URL.Query().Get("offset")
		allOrders, errOrders := ho.Service.GetAllOrders(int(idUser), limitStr, offsetStr)
		if errOrders != nil {
			ho.allOrders.Error = errOrders.Error()
			if errors.Is(errOrders, ErrOrderNotFound) {
				response.RespJs(writer, ho.allOrders, http.StatusNotFound)
			} else {
				response.RespJs(writer, ho.allOrders, http.StatusBadRequest)
			}
			return
		}
		response.RespJs(writer, allOrders, http.StatusOK)
	}
}
