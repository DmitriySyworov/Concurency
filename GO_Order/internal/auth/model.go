package auth

import (
	"time"
)

type Session struct {
	SessionId    string    `gorm:"not null;unique"`
	TempPassword int       `gorm:"not null"`
	ExpiresAt    time.Time `gorm:"not null"`
}
type TempUser struct {
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"not null"`
	Phone     string    `json:"phone" gorm:"not null"`
	Password  string    `json:"password" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	UserId    int       `json:"user_id" gorm:"not null;uniqueIndex:user_idx"`
	Jwt       string    `json:"jwt" gorm:"-"`
	Error     string    `json:"error" gorm:"-"`
}
