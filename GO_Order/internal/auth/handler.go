package auth

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

type HandlerAuth struct {
	*HandlerAuthDep
	RespUser    common.User
	RespAuth    ResponseAuth
	RespConfirm ResponseConfirm
}
type HandlerAuthDep struct {
	Service *ServiceAuth
	*configs.Config
}

func NewUserHandler(router *http.ServeMux, dep *HandlerAuthDep) {
	user := &HandlerAuth{
		HandlerAuthDep: dep,
	}
	router.HandleFunc("POST /auth/login/{method}", user.Login())
	router.HandleFunc("POST /auth/register/{method}", user.Register())
	router.HandleFunc("POST /auth/restore/{method}", user.Restore())
	router.Handle("POST /auth/confirm", middleware.IsTempUser(user.Confirm(), dep.Config))
}

func (h *HandlerAuth) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, errReq := requestJs.RequestHandler[RequestUserRegister](r)
		if errReq != nil {
			h.RespAuth.Error = errReq.Error()
			response.RespJs(w, h.RespAuth, http.StatusBadRequest)
			return
		}
		method := r.PathValue("method")
		if method != "phone" && method != "email" {
			h.RespAuth.Error = ErrMethodAuth.Error()
			response.RespJs(w, h.RespAuth, http.StatusBadRequest)
			return
		}
		respAuth, errRegister := h.Service.Register(body, method)
		if errRegister != nil {
			h.RespAuth.Error = errRegister.Error()
			if errors.Is(errRegister, ErrReg) {
				response.RespJs(w, h.RespAuth, http.StatusBadRequest)
			} else {
				response.RespJs(w, h.RespAuth, http.StatusInternalServerError)
			}
			return
		}
		response.RespJs(w, respAuth, http.StatusCreated)
	}
}
func (h *HandlerAuth) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, errReq := requestJs.RequestHandler[RequestUserLogin](r)
		if errReq != nil || (body.Email == "" && body.Phone == "") || (body.Email != "" && body.Phone != "") {
			h.RespAuth.Error = custerrors.ErrInvalidData.Error()
			response.RespJs(w, h.RespAuth, http.StatusBadRequest)
			return
		}
		method := r.PathValue("method")
		if (method != "phone" && method != "email") || (method == "phone" && body.Email != "") || (method == "email" && body.Phone != "") {
			h.RespAuth.Error = ErrMethodAuth.Error()
			response.RespJs(w, h.RespAuth, http.StatusBadRequest)
			return
		}
		respJwt, errLog := h.Service.Login(body, method)
		if errLog != nil {
			h.RespAuth.Error = errLog.Error()
			if errors.Is(errLog, ErrSecurity) || errors.Is(errLog, ErrSendEmail) {
				response.RespJs(w, h.RespAuth, http.StatusInternalServerError)
			} else {
				response.RespJs(w, h.RespAuth, http.StatusUnauthorized)
			}
			return
		}
		response.RespJs(w, respJwt, http.StatusOK)
	}
}
func (h *HandlerAuth) Restore() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, errBody := requestJs.RequestHandler[RequestRestore](request)
		if errBody != nil {
			h.RespAuth.Error = custerrors.ErrInvalidData.Error()
			response.RespJs(writer, h.RespAuth, http.StatusBadRequest)
			return
		}
		method := request.PathValue("method")
		if (method != "phone" && method != "email") || (method == "phone" && body.Email != "") || (method == "email" && body.Phone != "") {
			h.RespAuth.Error = ErrMethodAuth.Error()
			response.RespJs(writer, h.RespAuth, http.StatusBadRequest)
			return
		}
		respRestore, errRestore := h.Service.Restore(body, method)
		if errRestore != nil {
			h.RespAuth.Error = errRestore.Error()
			if errors.Is(errRestore, ErrSecurity) || errors.Is(errRestore, ErrSendEmail) {
				response.RespJs(writer, h.RespAuth, http.StatusInternalServerError)
			} else {
				response.RespJs(writer, h.RespAuth, http.StatusUnauthorized)
			}
			return
		}
		response.RespJs(writer, respRestore, http.StatusOK)
	}
}
func (h *HandlerAuth) Confirm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tempEmail, ok := r.Context().Value(middleware.KeyTempUser).(string)
		if !ok {
			h.RespConfirm.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(w, h.RespConfirm, http.StatusUnauthorized)
			return
		}
		body, errReq := requestJs.RequestHandler[RequestConfirm](r)
		if errReq != nil {
			h.RespConfirm.Error = errReq.Error()
			response.RespJs(w, h.RespConfirm, http.StatusBadRequest)
			return
		}
		sessId := r.Header.Get("Authorization")
		action := r.URL.Query().Get("action")
		if sessId == "" || (action != "register" && action != "login" && action != "restore") {
			h.RespConfirm.Error = ErrMissing.Error()
			response.RespJs(w, h.RespConfirm, http.StatusBadRequest)
			return
		}
		respConfirm, errConfirm := h.Service.Confirm(body, sessId, tempEmail, action)
		if errConfirm != nil {
			h.RespConfirm.Error = errConfirm.Error()
			h.RespConfirm.Jwt = ""
			if errors.Is(errConfirm, ErrSecurity) || errors.Is(errConfirm, ErrCreateUser) {
				response.RespJs(w, h.RespConfirm, http.StatusInternalServerError)
			} else {
				response.RespJs(w, h.RespConfirm, http.StatusUnauthorized)
			}
			return
		}
		h.RespConfirm.Error = ""
		response.RespJs(w, respConfirm, http.StatusCreated)
	}
}
