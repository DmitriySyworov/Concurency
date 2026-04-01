package auth

import (
	"order/app/pkg/db"
	"time"
)

type RepositoryAuth struct {
	*db.Db
}

func NewRepositoryAuth(database *db.Db) *RepositoryAuth {
	return &RepositoryAuth{
		Db: database,
	}
}

func (r *RepositoryAuth) CreateTempUser(tempUser *TempUser) error {
	result := r.DB.Create(&tempUser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *RepositoryAuth) CreateSession(sess *Session) error {
	result := r.DB.Create(&sess)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *RepositoryAuth) GetTempUser(email string) (*TempUser, error) {
	var tempUser TempUser
	res := r.DB.Where("email = ?", email).First(&tempUser)
	if res.Error != nil {
		return nil, res.Error
	}
	return &tempUser, nil
}

func (r *RepositoryAuth) GetSession(sessId string) (Session, error) {
	var sess Session
	resultGet := r.DB.Where("session_id = ?", sessId).First(&sess)
	if resultGet.Error != nil {
		return sess, resultGet.Error
	}
	r.DB.Where("session_id = ?", sessId).Delete(&Session{})
	return sess, nil
}

func (r *RepositoryAuth) DeleteTempDB() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			r.DB.Where("expires_at < ?", time.Now()).Delete(&Session{})
			r.DB.Where("expires_at < ?", time.Now()).Delete(&TempUser{})
		}
	}
}
