package product

import (
	generaterand "order/app/pkg/generateRand"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string`json:"name" gorm:"not null"`
	Description string`json:"description" gorm:"not null"`
	Images      pq.StringArray`json:"image" gorm:"type:text[];not null"`
	Category string`json:"category" gorm:"not null"`
	Hash string`json:"hash" gorm:"not null;uniqueIndex:idx_hash"`
	Error string`json:"error" gorm:"-"`
}
func NewProduct(name, description, category, hash string, image []string)*Product{
return &Product{
	Name: name,
	Description: description,
	Images: image,
	Category: category,
	Hash: hash,
}
}
func UniqueCheckHash(getHash func(string)(*Product, error))string{
for {
	hash := generaterand.RandStr(8)
	_, notFound := getHash(hash)
	if notFound != nil{
		return hash 
	}
}
}