package di

import "order/app/internal/common"

type IProductRepo interface {
	GetByAllHash(string) (*common.Product, error)
}

type IUserRepo interface {
	GetByIdUser(int) (*common.User, error)
	CreateUser(*common.User) error
	GetByEmailOrPhone(string, string) (*common.User, error)
	GetByDeleteUser(string, string) (*common.User, error)
	RestoreUser(int) error
}
