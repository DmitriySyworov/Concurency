package product

import (
	"net/http"
	"order/app/pkg/request"
	"order/app/pkg/response"
)

type ProductHandler struct {
	Product
	ProductHandlerDep
}
type ProductHandlerDep struct {
	*ProductRepository
}

func NewProductHandler(router *http.ServeMux, setting ProductHandlerDep) {
	prod := &ProductHandler{
		ProductHandlerDep: setting,
	}
	router.HandleFunc("POST /product", prod.CreateProduct())
	router.HandleFunc("PATCH /product/{hash}", prod.UpdateProduct())
	router.HandleFunc("GET /product/{hash}", prod.GetProduct())
	router.HandleFunc("GET /product/all/{category}", prod.AllProduct())
	router.HandleFunc("DELETE /product/{hash}", prod.DeleteProduct())
}
func (ph *ProductHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, errReq := request.RequestHandler[RequestProductCreate](w, r)
		if errReq != nil {
			ph.Product.Error = errReq.Error()
			response.RespJs(w, ph.Product, http.StatusBadRequest)
			return
		}
		unuqHash := UniqueCheckHash(ph.GetByHash)
		product := NewProduct(body.Name, body.Description, body.Category, unuqHash, body.Images)
		errCreate := ph.Create(product)
		if errCreate != nil {
			ph.Product.Error = errCreate.Error()
			response.RespJs(w, ph.Product, http.StatusInternalServerError)
			return
		}
		response.RespJs(w, product, http.StatusCreated)
	}
}
func (ph *ProductHandler) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		record, errGet := ph.GetByHash(hash)
		if errGet != nil {
			ph.Product.Error = errGet.Error()
			response.RespJs(w, ph.Product, http.StatusNotFound)
			return
		}
		body, errReq := request.RequestHandler[RequestProductCreate](w, r)
		if errReq != nil {
			ph.Product.Error = errReq.Error()
			response.RespJs(w, ph.Product, http.StatusBadRequest)
			return
		}
		record.Name = body.Name
		record.Images = body.Images
		record.Category = body.Category
		record.Description = body.Description
		resProd, errUpdate := ph.Update(record)
		if errUpdate != nil {
			ph.Product.Error = errUpdate.Error()
			response.RespJs(w, ph.Product, http.StatusBadRequest)
			return
		}
		response.RespJs(w, resProd, http.StatusCreated)
	}
}
func (ph *ProductHandler) GetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		record, errGet := ph.GetByHash(hash)
		if errGet != nil {
			ph.Product.Error = errGet.Error()
			response.RespJs(w, ph.Product, http.StatusNotFound)
			return
		}
		response.RespJs(w, record, http.StatusOK)
	}
}
func (ph *ProductHandler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		_, errGet := ph.GetByHash(hash)
		if errGet != nil {
			ph.Product.Error = errGet.Error()
			response.RespJs(w, ph.Product, http.StatusNotFound)
			return
		}
		errDel := ph.Delete(hash)
		if errDel != nil {
			ph.Product.Error = errDel.Error()
			response.RespJs(w, ph.Product, http.StatusInternalServerError)
			return
		}
		response.RespJs(w, nil, http.StatusOK)
	}
}
func (ph *ProductHandler) AllProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := r.PathValue("category")
		sliceProduct, errAll := ph.GetAll(category)
		if errAll != nil {
			ph.Product.Error = errAll.Error()
			response.RespJs(w, ph.Product, http.StatusInternalServerError)
			return
		}
		if len(sliceProduct) < 1 {
			ph.Product.Error = "There are no products in this category"
			response.RespJs(w, ph.Product, http.StatusNotFound)
			return
		}
		var payload ResponseSliceProduct
		payload.CategoryProduct = sliceProduct
		response.RespJs(w, payload, http.StatusOK)
	}
}

// func (ph *ProductHandler)uniquenessCheck(product *Product)error{
// sliceProd, errProd := ph.GetAll()
// if errProd != nil {
// 	return errProd
// }
// for _, productDb := range *sliceProd{
// 	if productDb.Name == product.Name{
// 		return errors.New("")
// 	}
// 	// if productDb.
// }
// return nil
// }
