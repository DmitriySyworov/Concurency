package user

import (
	"order/app/internal/common"
	custerrors "order/app/pkg/custErrors"
)

type ServiceUser struct {
	*RepositoryUser
}

func NewServiceUser(repo *RepositoryUser) *ServiceUser {
	return &ServiceUser{
		RepositoryUser: repo,
	}
}

func (su *ServiceUser) UpdateUser(body *RequestUpdateUser, idUser int) (*common.User, error) {
	user, errGet := su.RepositoryUser.GetByIdUser(idUser)
	if errGet != nil {
		return nil, custerrors.ErrUserDontExist
	}
	user.Name = body.Name
	user.Email = body.Email
	user.Password = body.Phone
	errUpdate := su.RepositoryUser.UpdateUser(user)
	if errUpdate != nil {
		return nil, ErrUpdateUser
	}
	return user, nil
}
func (su *ServiceUser) DeleteUser(idUser int) error {
	_, errGet := su.GetByIdUser(idUser)
	if errGet != nil {
		return custerrors.ErrUserDontExist
	}
	errDel := su.RepositoryUser.DeleteUser(idUser)
	if errDel != nil {
		return ErrDeleteUser
	}
	return nil
}
