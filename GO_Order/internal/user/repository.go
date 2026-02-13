package user

import (
	"order/app/pkg/db"
	"time"
)

type UserRepository struct {
	*db.Db
}

func NewUserRepository(database *db.Db) *UserRepository {
	return &UserRepository{
		Db: database,
	}
}

func (r *UserRepository) CreateUser(user *User) error{
	result := r.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) GetUser(email, phone string) (*User, error) {
	var user User
	resEm := r.DB.Where("email = ?", email).First(&user)
	resPh := r.DB.Where("phone = ?", phone).First(&user)
	if resEm.Error == nil {
		return &user, nil
	}
	if resPh.Error == nil {
		return &user, nil
	}
	return nil, resEm.Error
}
func (r *UserRepository) Session(sess *Session) error {
	result := r.DB.Create(&sess)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *UserRepository) GetSession(sessId string) (Session, error) {
	var sess Session
	r.DB.Where("expires_at < ?", time.Now()).Delete(&Session{})
	resultGet := r.DB.Where("session_id = ?", sessId).First(&sess)
	if resultGet != nil {
		return sess, resultGet.Error
	}
	r.DB.Where("session_id = ?", sessId).Delete(&Session{})
	return sess, nil
}
