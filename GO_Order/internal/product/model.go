package product

import (
	generaterand "order/app/pkg/generateRand"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description" gorm:"not null"`
	Images      pq.StringArray `json:"image" gorm:"type:text[];not null"`
	Category    string         `json:"category" gorm:"not null"`
	Hash        string         `json:"hash" gorm:"not null;uniqueIndex:idx_hash"`
	IdUser      string         `json:"id_user" gorm:"not null;index:idx_user"`
	Error       string         `json:"error" gorm:"-"`
}

func NewProduct(name, description, category, hash, idUser string, image []string) *Product {
	return &Product{
		Name:        name,
		Description: description,
		Images:      image,
		Category:    category,
		Hash:        hash,
		IdUser:      idUser,
	}
}
func UniqueCheckHash(getHash func(string, string) (*Product, error), idUser string) string {
	for {
		hash := generaterand.RandStr(8)
		_, notFound := getHash(hash, idUser)
		if notFound != nil {
			return hash
		}
	}
}
