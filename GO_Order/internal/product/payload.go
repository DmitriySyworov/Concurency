package product

type RequestProductCreate struct {
	Name        string`json:"name" validate:"required"`
	Description string`json:"description" validate:"required"`
	Images      []string`json:"image" validate:"required"`
	Category string`json:"category" validate:"required"`
}
type RequestProductUpdate struct {
	Name        string`json:"name"`
	Description string`json:"description"`
	Images      []string`json:"image"`
	Category string`json:"category"`
}
type ResponseSliceProduct struct{
	CategoryProduct []Product`json:"products by category"`
}