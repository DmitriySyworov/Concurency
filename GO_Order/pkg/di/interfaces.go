package di

import "order/app/internal/common"

type IProductRepo interface {
	GetByAllHash(string) (*common.Product, error)
}
