package product

import (
	"order/app/pkg/db"

	"gorm.io/gorm/clause"
)

type ProductRepository struct {
	Database *db.Db
}

func NewProductRepository(database *db.Db) *ProductRepository {
	return &ProductRepository{
		Database: database,
	}
}
func (repo *ProductRepository) Create(product *Product) error {
	result := repo.Database.DB.Create(product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (repo *ProductRepository) Update(product *Product, idUser string) (*Product, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Where("id_user = ?", idUser).Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}
func (repo *ProductRepository) GetByHash(hash, idUser string) (*Product, error) {
	var product Product
	result := repo.Database.DB.Where("hash = ? AND id_user = ?", hash, idUser).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}
func (repo *ProductRepository) Delete(hash, idUser string) error {
	result := repo.Database.DB.Where("hash = ? AND id_user = ?", hash, idUser).Delete(&Product{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (repo *ProductRepository) GetAll(category, idUser string) (*ResponseSliceProduct, error) {
	var sliceProduct ResponseSliceProduct
	result := repo.Database.DB.Where("category = ? AND id_user = ?", category, idUser).Find(&sliceProduct.CategoryProduct)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sliceProduct, nil
}
