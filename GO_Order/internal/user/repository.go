package user

import (
	"fmt"
	"order/app/internal/common"
	custerrors "order/app/pkg/custErrors"
	"order/app/pkg/db"

	"gorm.io/gorm/clause"
)

type RepositoryUser struct {
	*db.Db
}

func NewUserRepository(db *db.Db) *RepositoryUser {
	return &RepositoryUser{
		Db: db,
	}
}
func (db *RepositoryUser) UpdateUser(user *common.User) error {
	res := db.Clauses(clause.Returning{}).
		Where("deleted_at is null AND user_id = ?", user.UserId).
		Updates(user)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (db *RepositoryUser) DeleteUser(userId int) error {
	res := db.Where("user_id = ?", userId).Delete(&common.User{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (db *RepositoryUser) CreateUser(user *common.User) error {
	result := db.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (db *RepositoryUser) GetByIdUser(idUser int) (*common.User, error) {
	var user common.User
	res := db.Db.Where("user_id = ? AND deleted_at is null", idUser).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
func (db *RepositoryUser) GetByEmailOrPhone(email, phone string) (*common.User, error) {
	var user common.User
	res := db.DB.Where("deleted_at is null AND email = ? OR phone = ?", email, phone).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
func (db *RepositoryUser) GetByDeleteUser(email, phone string) (*common.User, error) {
	var user common.User
	res := db.DB.Unscoped().Where("email = ? OR phone = ?", email, phone).First(&user)
	fmt.Println(user.DeletedAt)
	if res.Error != nil {
		return nil, res.Error
	}
	if !user.DeletedAt.Valid {
		return nil, custerrors.ErrUserNotFound
	}
	return &user, nil
}
func (db *RepositoryUser) RestoreUser(userId int) error {
	res := db.Model(&common.User{}).Unscoped().Where("user_id = ?", userId).Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
