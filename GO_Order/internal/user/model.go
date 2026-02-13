package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct{
	gorm.Model
	Name string`json:"name" gorm:"not null"`
	Email string`json:"email" gorm:"not null;uniqueIndex:idx_email"`
	Phone string`json:"phone" gorm:"not null;unique"`
	Password string`json:"password" gorm:"not null"`
	Jwt string `json:"jwt" gorm:"-"`
	Error string`json:"error" gorm:"-"`
}

type Session struct{
	SessionId  string`gorm:"not null;unique"`
	TempPassword string`gorm:"not null"`
	ExpiresAt time.Time`gorm:"not null"`
}