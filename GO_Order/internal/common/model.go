package common

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID       uint      `gorm:"primaryKey"`
	Products []Product `gorm:"many2many:order_products"`
	UserId   int       `json:"user_id" gorm:"uniqueIndex:user_idx;not null"`
	OrderId  string    `json:"order_id" gorm:"not null"`
	Error    string    `json:"error" gorm:"-"`
}

type Product struct {
	gorm.Model
	ID          uint           `gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	Images      pq.StringArray `json:"image" gorm:"type:text[];not null"`
	Category    string         `json:"category" gorm:"not null"`
	Hash        string         `json:"hash" gorm:"not null;uniqueIndex"`
	UserId      int            `json:"user_id" gorm:"not null;uniqueIndex:user_idx"`
	Orders      []Order        `gorm:"many2many:order_products"`
	Error       string         `json:"error" gorm:"-"`
}

type User struct {
	gorm.Model
	ID       uint    `gorm:"primaryKey"`
	Name     string  `json:"name" gorm:"not null"`
	Email    string  `json:"email" gorm:"not null;uniqueIndex:idx_email"`
	Phone    string  `json:"phone" gorm:"not null;unique"`
	Password string  `json:"password" gorm:"not null" `
	Orders   []Order `gorm:"foreignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	UserId   int     `json:"user_id" gorm:"not null;uniqueIndex:user_idx"`
	Jwt      string  `json:"jwt" gorm:"-" `
	Error    string  `json:"error" gorm:"-" `
}
