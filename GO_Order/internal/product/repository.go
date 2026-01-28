package product

import (
	"order/app/pkg/db"

	"gorm.io/gorm/clause"
)

type ProductRepository struct{
	Database *db.Db
}

func NewProductRepository(database *db.Db)*ProductRepository{
	return &ProductRepository{
		Database: database,
	}
}
func (repo *ProductRepository)Create(product *Product)error{
result := repo.Database.DB.Create(product)
if result.Error != nil {
	return result.Error
}
return nil
}
func (repo *ProductRepository)Update(product *Product)(*Product, error){
	result := repo.Database.DB.Clauses(clause.Returning{}).Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}
return product, nil
}
func (repo *ProductRepository)GetByHash(hash string)(*Product, error){
	var product Product
	result := repo.Database.DB.Where("hash = ?", hash).First(&product)
	if result.Error != nil {
	return nil, result.Error
	}
return &product, nil
}
func (repo *ProductRepository)Delete(hash string)error{
result := repo.Database.DB.Where("hash = ?", hash).Delete(&Product{})
if result.Error != nil {
	return result.Error
}
return nil
}
func (repo *ProductRepository)GetAll(category string)([]Product, error){
	var sliceProduct []Product
	result := repo.Database.DB.Where("category = ?", category).Find(&sliceProduct)
	if result.Error != nil {
		return nil, result.Error
	}
return sliceProduct, nil
}