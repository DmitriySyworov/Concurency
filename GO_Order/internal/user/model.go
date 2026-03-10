package user

import (
	"order/app/internal/common"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint           `gorm:"primaryKey"`
	Name     string         `json:"name" gorm:"not null"`
	Email    string         `json:"email" gorm:"not null;uniqueIndex:idx_email"`
	Phone    string         `json:"phone" gorm:"not null;unique"`
	Password string         `json:"password" gorm:"not null" `
	Orders   []common.Order `gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	UserId   int            `json:"user_id" gorm:"not null;uniqueIndex:user_idx"`
	Jwt      string         `json:"jwt" gorm:"-" `
	Error    string         `json:"error" gorm:"-" `
}

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
	UserId    int       `json:"user_id" gorm:"type:bigint;not null;uniqueIndex:user_idx"`
	Jwt       string    `json:"jwt" gorm:"-"`
	Error     string    `json:"error" gorm:"-"`
}
