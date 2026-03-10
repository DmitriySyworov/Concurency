package user

import (
	"context"
	"log"
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

func (r *UserRepository) CreateUser(user *User) error {
	result := r.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *UserRepository) GetByIdUser(idUser int) error {
	var user User
	res := r.Db.Where("user_id = ?", idUser).First(&user)
	return res.Error
}
func (r *UserRepository) GetByEmailOrPhone(email, phone string) (*User, error) {
	var user User
	res := r.DB.Where("email = ? OR phone = ?", email, phone).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}
func (r *UserRepository) CreateTempUser(tempUser *TempUser) error {
	result := r.DB.Create(&tempUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *UserRepository) CreateSession(sess *Session) error {
	result := r.DB.Create(&sess)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *UserRepository) GetTempUser(email string) (*TempUser, error) {
	var tempUser TempUser
	res := r.DB.Where("email = ?", email).First(&tempUser)
	if res.Error != nil {
		return nil, res.Error
	}
	return &tempUser, nil
}

func (r *UserRepository) GetSession(sessId string) (Session, error) {
	var sess Session
	resultGet := r.DB.Where("session_id = ?", sessId).First(&sess)
	if resultGet.Error != nil {
		return sess, resultGet.Error
	}
	r.DB.Where("session_id = ?", sessId).Delete(&Session{})
	return sess, nil
}

func (r *UserRepository) DeleteTempDB(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.DB.Where("expires_at < ?", time.Now()).Delete(&Session{})
			r.DB.Where("expires_at < ?", time.Now()).Delete(&TempUser{})
		case <-ctx.Done():
			log.Println(ctx.Err())
			return
		}
	}
}
