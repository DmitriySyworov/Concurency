package user

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

type HandlerUser struct {
	common.User
	RequestUpdateUser
	*HandlerUserDep
}
type HandlerUserDep struct {
	*ServiceUser
	*configs.Config
}

func NewHandlerUser(router *http.ServeMux, dep *HandlerUserDep) {
	user := &HandlerUser{
		HandlerUserDep: dep,
	}
	router.Handle("GET /users/my", middleware.IsAuthID(user.GetUser(), dep.Config))
	router.Handle("PATCH /users/my", middleware.IsAuthID(user.UpdateUser(), dep.Config))
	router.Handle("DELETE /users/my", middleware.IsAuthID(user.DeleteUser(), dep.Config))
}

func (hl *HandlerUser) GetUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idUser, ok := request.Context().Value(middleware.KeyIDUser).(float64)
		if !ok {
			hl.User.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(writer, hl.User, http.StatusUnauthorized)
			return
		}
		user, errGet := hl.GetByIdUser(int(idUser))
		if errGet != nil {
			hl.User.Error = custerrors.ErrUserNotFound.Error()
			response.RespJs(writer, hl.User, http.StatusNotFound)
			return
		}
		response.RespJs(writer, user, http.StatusOK)
	}
}

func (hl *HandlerUser) UpdateUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idUser, ok := request.Context().Value(middleware.KeyIDUser).(float64)
		if !ok {
			hl.User.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(writer, hl.User, http.StatusUnauthorized)
			return
		}
		body, errBody := requestJs.RequestHandler[RequestUpdateUser](request)
		if errBody != nil {
			hl.User.Error = custerrors.ErrInvalidData.Error()
			response.RespJs(writer, hl.User, http.StatusBadRequest)
			return
		}
		newUser, errUpdate := hl.ServiceUser.UpdateUser(body, int(idUser))
		if errUpdate != nil {
			hl.User.Error = errUpdate.Error()
			if errors.Is(errUpdate, custerrors.ErrUserDontExist) {
				response.RespJs(writer, hl.User, http.StatusBadRequest)
			} else {
				response.RespJs(writer, hl.User, http.StatusInternalServerError)
			}
			return
		}
		response.RespJs(writer, newUser, http.StatusCreated)
	}
}
func (hl *HandlerUser) DeleteUser() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		idUser, ok := request.Context().Value(middleware.KeyIDUser).(float64)
		if !ok {
			http.Error(writer, custerrors.ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}
		errDelete := hl.ServiceUser.DeleteUser(int(idUser))
		if errDelete != nil {
			if errors.Is(errDelete, custerrors.ErrUserDontExist) {
				http.Error(writer, errDelete.Error(), http.StatusBadRequest)
			} else {
				http.Error(writer, errDelete.Error(), http.StatusInternalServerError)
			}
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}
