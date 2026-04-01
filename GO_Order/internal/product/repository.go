package product

import (
	"order/app/internal/common"
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
func (repo *ProductRepository) CreateProduct(product *common.Product) error {
	result := repo.Database.DB.Table("products").Create(&product)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (repo *ProductRepository) UpdateProduct(product *common.Product, idUser int) (*common.Product, error) {
	result := repo.Database.DB.Clauses(clause.Returning{}).Where("user_id = ? AND deleted_at is null", idUser).Updates(product)
	if result.Error != nil {
		return nil, result.Error
	}
	return product, nil
}
func (repo *ProductRepository) GetByHash(hash string, idUser int) (*common.Product, error) {
	var product common.Product
	result := repo.Database.DB.Where("hash = ? AND user_id = ? AND deleted_at is null", hash, idUser).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}
func (repo *ProductRepository) GetByAllHash(hash string) (*common.Product, error) {
	var product common.Product
	result := repo.Database.DB.Where("hash = ? AND deleted_at is null", hash).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}
func (repo *ProductRepository) DeleteProduct(hash string, idUser int) error {
	result := repo.Database.DB.Where("hash = ? AND user_id = ? AND deleted_at is null", hash, idUser).Delete(&common.Product{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (repo *ProductRepository) GetAllProduct(category string, idUser int) (*ResponseSliceProduct, error) {
	var sliceProduct ResponseSliceProduct
	result := repo.Database.DB.Where("category = ? AND user_id = ? AND deleted_at is null", category, idUser).Find(&sliceProduct.CategoryProduct)
	if result.Error != nil {
		return nil, result.Error
	}
	return &sliceProduct, nil
}
