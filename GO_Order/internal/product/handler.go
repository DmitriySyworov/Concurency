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
	*ProductService
}

func NewProductHandler(router *http.ServeMux, setting ProductHandlerDep) {
	prod := &ProductHandler{
		ProductHandlerDep: setting,
	}
	router.HandleFunc("POST /product", prod.HandlerCreateProduct())
	router.HandleFunc("PATCH /product/{hash}", prod.HandlerUpdateProduct())
	router.HandleFunc("GET /product/{hash}", prod.HadlerGetProduct())
	router.HandleFunc("GET /product", prod.HandlerAllProduct())
	router.HandleFunc("DELETE /product/{hash}", prod.HandlerDeleteProduct())
}
func (h *ProductHandler) HandlerCreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, errReq := request.RequestHandler[RequestProductCreate](w, r)
		if errReq != nil {
			h.Product.Error = errReq.Error()
			response.RespJs(w, h.Product, http.StatusBadRequest)
			return
		}
		product, errCreate := h.ServiceCreate(body)
		if errCreate != nil {
			h.Product.Error = errCreate.Error()
			response.RespJs(w, h.Product, http.StatusInternalServerError)
			return
		}
		response.RespJs(w, product, http.StatusCreated)
	}
}
func (h *ProductHandler) HandlerUpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, errReq := request.RequestHandler[RequestProductUpdate](w, r)
		if errReq != nil {
			h.Product.Error = errReq.Error()
			response.RespJs(w, h.Product, http.StatusBadRequest)
			return
		}
		hash := r.PathValue("hash")
		resProd, errUpdate := h.ServiceUpdate(body, hash)
		if errUpdate != nil {
			h.Product.Error = errUpdate.Error()
			if errUpdate == ErrNotFoundProduct {
				response.RespJs(w, h.Product, http.StatusNotFound)
			} else {
				response.RespJs(w, h.Product, http.StatusInternalServerError)
			}
			return
		}
		response.RespJs(w, resProd, http.StatusCreated)
	}
}
func (h *ProductHandler) HadlerGetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		record, errGet := h.GetByHash(hash)
		if errGet != nil {
			h.Product.Error = ErrNotFoundProduct.Error()
			response.RespJs(w, h.Product, http.StatusNotFound)
			return
		}
		response.RespJs(w, record, http.StatusOK)
	}
}
func (h *ProductHandler) HandlerDeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		errDel := h.ServiceDelete(hash)
		if errDel != nil {
			h.Product.Error = errDel.Error()
			if errDel == ErrNotFoundProduct {
				response.RespJs(w, h.Product, http.StatusNotFound)
			} else {
				response.RespJs(w, h.Product, http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
func (h *ProductHandler) HandlerAllProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := r.URL.Query().Get("category")
		categoruProducts, errCategory := h.ServiceAllProduct(category)
		if errCategory != nil {
			h.Product.Error = errCategory.Error()
			if errCategory == ErrNotCategoryProduct {
				response.RespJs(w, h.Product, http.StatusNotFound)
			} else {
				response.RespJs(w, h.Product, http.StatusInternalServerError)
			}
			return
		}
		response.RespJs(w, categoruProducts, http.StatusOK)
	}
}
