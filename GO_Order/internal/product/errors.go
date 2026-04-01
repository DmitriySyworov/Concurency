package product

import "errors"

var (
	ErrCreateProduct      = errors.New("failed to save product")
	ErrNotFoundProduct    = errors.New("no product found for the specified parameter")
	ErrNotUpdateProduct   = errors.New("failed to update product record")
	ErrNotDeleteProduct   = errors.New("failed to delete product entry")
	ErrNotCategoryProduct = errors.New("there are no products in this category")
	ErrRecordCategory     = errors.New("failed to record for this category")
)
