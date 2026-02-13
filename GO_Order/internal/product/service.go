package product

type ProductService struct {
	*ProductRepository
}

func NewProductService(repo *ProductRepository) *ProductService {
	return &ProductService{
		ProductRepository: repo,
	}
}
func (s *ProductService) ServiceCreate(body *RequestProductCreate) (*Product, error) {
	unuqHash := UniqueCheckHash(s.GetByHash)
	product := NewProduct(body.Name, body.Description, body.Category, unuqHash, body.Images)
	errCreate := s.Create(product)
	if errCreate != nil {
		return nil, ErrCreateProduct
	}
	return product, nil
}
func (s *ProductService) ServiceUpdate(body *RequestProductUpdate, hash string) (*Product, error) {
	record, errGet := s.GetByHash(hash)
	if errGet != nil {
		return nil, ErrNotFoundProduct
	}
	record.Name = body.Name
	record.Images = body.Images
	record.Category = body.Category
	record.Description = body.Description
	resProd, errUpdate := s.Update(record)
	if errUpdate != nil {
		return nil, ErrNotUpdateProduct
	}
	return resProd, nil
}
func (s *ProductService) ServiceDelete(hash string) error {
	_, errGet := s.GetByHash(hash)
	if errGet != nil {
		return ErrNotFoundProduct
	}
	errDel := s.Delete(hash)
	if errDel != nil {
		return ErrNotDeleteProduct
	}
	return nil
}
func (s *ProductService)ServiceAllProduct(category string) (*ResponseSliceProduct, error){
		sliceProduct, errAll := s.GetAll(category)
		if errAll != nil {
			return nil, ErrRecordCategory
		}
		if len(sliceProduct.CategoryProduct) < 1 {
		return nil, ErrNotCategoryProduct
		}
		return sliceProduct, nil
}