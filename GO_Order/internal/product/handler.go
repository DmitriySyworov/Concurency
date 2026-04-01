package product

import (
	"errors"
	"net/http"
	"order/app/configs"
	"order/app/internal/common"
	custerrors "order/app/pkg/custErrors"
	"order/app/pkg/middleware"
	"order/app/pkg/requestJs"
	"order/app/pkg/response"
)

type HandlerProduct struct {
	common.Product
	HandlerProductDep
}
type HandlerProductDep struct {
	Service *ProductService
	*configs.Config
}

func NewHandlerProduct(router *http.ServeMux, setting HandlerProductDep) {
	prod := &HandlerProduct{
		HandlerProductDep: setting,
	}
	router.Handle("POST /users/products", middleware.IsAuthID(prod.CreateProduct(), setting.Config))
	router.Handle("PATCH /users/products/{hash}", middleware.IsAuthID(prod.UpdateProduct(), setting.Config))
	router.Handle("GET /users/products/{hash}", middleware.IsAuthID(prod.GetProduct(), setting.Config))
	router.Handle("GET /users/products", middleware.IsAuthID(prod.AllProduct(), setting.Config))
	router.Handle("DELETE /users/products/{hash}", middleware.IsAuthID(prod.DeleteProduct(), setting.Config))
}
func (h *HandlerProduct) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idUser, ok := r.Context().Value(middleware.KeyIDUser).(float64)
		if !ok {
			h.Product.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(w, h.Product, http.StatusUnauthorized)
			return
		}
		body, errReq := requestJs.RequestHandler[RequestProductCreate](r)
		if errReq != nil {
			h.Product.Error = errReq.Error()
			response.RespJs(w, h.Product, http.StatusBadRequest)
			return
		}
		product, errCreate := h.Service.CreateProduct(body, int(idUser))
		if errCreate != nil {
			h.Product.Error = errCreate.Error()
			if errors.Is(errCreate, custerrors.ErrUserDontExist) {
				response.RespJs(w, h.Product, http.StatusUnauthorized)
			} else {
				response.RespJs(w, h.Product, http.StatusInternalServerError)
			}
			return
		}
		response.RespJs(w, product, http.StatusCreated)
	}
}
func (h *HandlerProduct) UpdateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idUser, ok := r.Context().Value(middleware.KeyIDUser).(float64)
		if !ok {
			h.Product.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(w, h.Product, http.StatusUnauthorized)
			return
		}

		body, errReq := requestJs.RequestHandler[RequestProductUpdate](r)
		if errReq != nil {
			h.Product.Error = errReq.Error()
			response.RespJs(w, h.Product, http.StatusBadRequest)
			return
		}
		hash := r.PathValue("hash")
		resProd, errUpdate := h.Service.UpdateProduct(body, hash, int(idUser))
		if errUpdate != nil {
			h.Product.Error = errUpdate.Error()
			if errors.Is(errUpdate, ErrNotFoundProduct) {
				response.RespJs(w, h.Product, http.StatusNotFound)
			} else {
				response.RespJs(w, h.Product, http.StatusInternalServerError)
			}
			return
		}
		response.RespJs(w, resProd, http.StatusCreated)
	}
}
func (h *HandlerProduct) GetProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idUser, ok := r.Context().Value(middleware.KeyIDUser).(float64)
		if !ok {
			h.Product.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(w, h.Product, http.StatusUnauthorized)
			return
		}

		hash := r.PathValue("hash")
		record, errGet := h.Service.Repo.GetByHash(hash, int(idUser))
		if errGet != nil {
			h.Product.Error = ErrNotFoundProduct.Error()
			response.RespJs(w, h.Product, http.StatusNotFound)
			return
		}
		response.RespJs(w, record, http.StatusOK)
	}
}
func (h *HandlerProduct) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idUser, ok := r.Context().Value(middleware.KeyIDUser).(float64)
		if !ok {
			h.Product.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(w, h.Product, http.StatusUnauthorized)
			return
		}

		hash := r.PathValue("hash")
		errDel := h.Service.DeleteProduct(hash, int(idUser))
		if errDel != nil {
			h.Product.Error = errDel.Error()
			if errors.Is(errDel, ErrNotFoundProduct) {
				response.RespJs(w, h.Product, http.StatusNotFound)
			} else {
				response.RespJs(w, h.Product, http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
func (h *HandlerProduct) AllProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idUser, ok := r.Context().Value(middleware.KeyIDUser).(float64)
		if !ok {
			h.Product.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(w, h.Product, http.StatusUnauthorized)
			return
		}

		category := r.URL.Query().Get("category")
		categoryProducts, errCategory := h.Service.AllProduct(category, int(idUser))
		if errCategory != nil {
			h.Product.Error = errCategory.Error()
			if errors.Is(errCategory, ErrNotCategoryProduct) {
				response.RespJs(w, h.Product, http.StatusNotFound)
			} else {
				response.RespJs(w, h.Product, http.StatusInternalServerError)
			}
			return
		}
		response.RespJs(w, categoryProducts, http.StatusOK)
	}
}
