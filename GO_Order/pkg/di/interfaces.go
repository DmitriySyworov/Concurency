package di

import "order/app/internal/common"

type IProductRepo interface {
	GetByAllHash(string) (*common.Product, error)
}

type IUserRepo interface {
	GetByIdUser(int) error
}
