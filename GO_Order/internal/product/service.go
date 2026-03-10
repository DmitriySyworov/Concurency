package product

import (
	"order/app/internal/common"
	generaterand "order/app/pkg/generateRand"
)

type ProductService struct {
	Repo *ProductRepository
}

func NewProductService(repo *ProductRepository) *ProductService {
	return &ProductService{
		Repo: repo,
	}
}
func (s *ProductService) CreateProduct(body *RequestProductCreate, idUser int) (*common.Product, error) {
	var resHash string
	for {
		resHash = generaterand.RandStr(8)
		_, notFound := s.Repo.GetByHash(resHash, idUser)
		if notFound != nil {
			break
		}
	}
	product := &common.Product{
		Name:        body.Name,
		Description: body.Description,
		Category:    body.Category,
		Hash:        resHash,
		UserId:      idUser,
		Images:      body.Images,
	}
	errCreate := s.Repo.CreateProduct(product)
	if errCreate != nil {
		return nil, ErrCreateProduct
	}
	return product, nil
}

func (s *ProductService) UpdateProduct(body *RequestProductUpdate, hash string, idUser int) (*common.Product, error) {
	record, errGet := s.Repo.GetByHash(hash, idUser)
	if errGet != nil {
		return nil, ErrNotFoundProduct
	}
	record.Name = body.Name
	record.Images = body.Images
	record.Category = body.Category
	record.Description = body.Description
	resProd, errUpdate := s.Repo.UpdateProduct(record, idUser)
	if errUpdate != nil {
		return nil, ErrNotUpdateProduct
	}
	return resProd, nil
}
func (s *ProductService) DeleteProduct(hash string, idUser int) error {
	_, errGet := s.Repo.GetByHash(hash, idUser)
	if errGet != nil {
		return ErrNotFoundProduct
	}
	errDel := s.DeleteProduct(hash, idUser)
	if errDel != nil {
		return ErrNotDeleteProduct
	}
	return nil
}
func (s *ProductService) AllProduct(category string, idUser int) (*ResponseSliceProduct, error) {
	sliceProduct, errAll := s.Repo.GetAllProduct(category, idUser)
	if errAll != nil {
		return nil, ErrRecordCategory
	}
	if len(sliceProduct.CategoryProduct) < 1 {
		return nil, ErrNotCategoryProduct
	}
	return sliceProduct, nil
}
