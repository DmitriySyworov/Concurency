package product

import "errors"


var (
	ErrCreateProduct = errors.New("Failed to save product")
	ErrNotFoundProduct = errors.New("no product found for the specified parameter")
	ErrNotUpdateProduct = errors.New("Failed to update product record")
	ErrNotDeleteProduct = errors.New("Failed to delete product entry")
	ErrNotCategoryProduct = errors.New("There are no products in this category")
	ErrRecordCategory = errors.New("Failed to record for this category")
)