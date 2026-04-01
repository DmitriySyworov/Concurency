package order

import (
	"order/app/internal/common"
)

type CreateOrderRequest struct {
	ProductsHash []string `json:"products_hash"`
}
type AllOrdersResponse struct {
	MyOrders []common.Order `json:"my-orders"`
	MetaData `json:"metadata"`
	Error    string `json:"error"`
}
type MetaData struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
