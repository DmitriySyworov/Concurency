package product

import (
	"net/http"
	"order/app/configs"
	custerrors "order/app/pkg/custErrors"
	"order/app/pkg/middleware"
	"order/app/pkg/request"
	"order/app/pkg/response"
)

type ProductHandler struct {
	Product
	ProductHandlerDep
}
type ProductHandlerDep struct {
	*ProductService
	*configs.Config
}

func NewProductHandler(router *http.ServeMux, setting ProductHandlerDep) {
	prod := &ProductHandler{
		ProductHandlerDep: setting,
	}
	router.Handle("POST /product", middleware.IsAuth(prod.HandlerCreateProduct(), setting.Config))
	router.Handle("PATCH /product/{hash}", middleware.IsAuth(prod.HandlerUpdateProduct(), setting.Config))
	router.Handle("GET /product/{hash}", middleware.IsAuth(prod.HadlerGetProduct(), setting.Config))
	router.Handle("GET /product", middleware.IsAuth(prod.HandlerAllProduct(), setting.Config))
	router.Handle("DELETE /product/{hash}", middleware.IsAuth(prod.HandlerDeleteProduct(), setting.Config))
}
func (h *ProductHandler) HandlerCreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idUser, ok := r.Context().Value(middleware.KeyIDUser).(string)
		if !ok {
			h.Product.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(w, h.Product, http.StatusUnauthorized)
			return
		}

		body, errReq := request.RequestHandler[RequestProductCreate](w, r)
		if errReq != nil {
			h.Product.Error = errReq.Error()
			response.RespJs(w, h.Product, http.StatusBadRequest)
			return
		}
		product, errCreate := h.ServiceCreate(body, idUser)
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
		idUser, ok := r.Context().Value(middleware.KeyIDUser).(string)
		if !ok {
			h.Product.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(w, h.Product, http.StatusUnauthorized)
			return
		}

		body, errReq := request.RequestHandler[RequestProductUpdate](w, r)
		if errReq != nil {
			h.Product.Error = errReq.Error()
			response.RespJs(w, h.Product, http.StatusBadRequest)
			return
		}
		hash := r.PathValue("hash")
		resProd, errUpdate := h.ServiceUpdate(body, hash, idUser)
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
		idUser, ok := r.Context().Value(middleware.KeyIDUser).(string)
		if !ok {
			h.Product.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(w, h.Product, http.StatusUnauthorized)
			return
		}

		hash := r.PathValue("hash")
		record, errGet := h.GetByHash(hash, idUser)
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
		idUser, ok := r.Context().Value(middleware.KeyIDUser).(string)
		if !ok {
			h.Product.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(w, h.Product, http.StatusUnauthorized)
			return
		}

		hash := r.PathValue("hash")
		errDel := h.ServiceDelete(hash, idUser)
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
		idUser, ok := r.Context().Value(middleware.KeyIDUser).(string)
		if !ok {
			h.Product.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(w, h.Product, http.StatusUnauthorized)
			return
		}

		category := r.URL.Query().Get("category")
		categoruProducts, errCategory := h.ServiceAllProduct(category, idUser)
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
