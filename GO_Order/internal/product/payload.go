package product

type RequestProductCreate struct {
	Name        string`json:"name" validate:"required"`
	Description string`json:"description" validate:"required"`
	Images      []string`json:"image" validate:"required"`
	Category string`json:"category" validate:"required"`
}
type ResponseSliceProduct struct{
	CategoryProduct []Product`json:"products by category"`
}