package user

import (
	"errors"
	"net/http"
	"order/app/configs"
	custerrors "order/app/pkg/custErrors"
	"order/app/pkg/middleware"
	"order/app/pkg/request"
	"order/app/pkg/response"
)

type UserHandler struct {
	*UserHandlerDep
	RespUser    User
	RespAuth    ResponseAuth
	RespConfirm ResponseConfirm
	RespTUser   TempUser
}
type UserHandlerDep struct {
	Service *UserService
	*configs.Config
}

func NewUserHandler(router *http.ServeMux, dep *UserHandlerDep) {
	user := &UserHandler{
		UserHandlerDep: dep,
	}
	router.HandleFunc("POST /user/login", user.Login())
	router.HandleFunc("POST /user/regist", user.Register())
	router.Handle("POST /user/auth/{method}", middleware.IsTempUser(user.Auth(), dep.Config))
	router.Handle("POST /user/auth", middleware.IsTempUser(user.Confirm(), dep.Config))
}

func (h *UserHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, errReq := request.RequestHandler[RequestUserRegist](r)
		if errReq != nil {
			h.RespTUser.Error = errReq.Error()
			response.RespJs(w, h.RespTUser, http.StatusBadRequest)
			return
		}
		userReg, errRegister := h.Service.Register(body)
		if errRegister != nil {
			h.RespTUser.Error = errRegister.Error()
			if errors.Is(errRegister, ErrReg) {
				response.RespJs(w, h.RespTUser, http.StatusBadRequest)
			} else {
				response.RespJs(w, h.RespTUser, http.StatusInternalServerError)
			}
			return
		}
		response.RespJs(w, userReg, http.StatusCreated)
	}
}
func (h *UserHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, errReq := request.RequestHandler[RequestUserLogin](r)
		if errReq != nil || (body.Email == "" && body.Phone == "") || (body.Email != "" && body.Phone != "") {
			h.RespTUser.Error = custerrors.ErrInvalidData.Error()
			response.RespJs(w, h.RespTUser, http.StatusBadRequest)
			return
		}
		userLog, errLog := h.Service.Login(body)
		if errLog != nil {
			h.RespTUser.Error = errLog.Error()
			if errors.Is(errLog, ErrSecurity) {
				response.RespJs(w, h.RespTUser, http.StatusInternalServerError)
			} else {
				response.RespJs(w, h.RespTUser, http.StatusUnauthorized)
			}
			return
		}
		response.RespJs(w, userLog, http.StatusOK)
	}
}
func (h *UserHandler) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.PathValue("method")
		tempEmail, ok := r.Context().Value(middleware.KeyTempUser).(string)
		if !ok {
			h.RespAuth.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(w, h.RespAuth, http.StatusUnauthorized)
			return
		}
		if method == "" {
			h.RespAuth.Error = ErrMissing.Error()
			response.RespJs(w, h.RespAuth, http.StatusBadRequest)
			return
		}
		authRes, errAuth := h.Service.Auth(method, tempEmail)
		if errAuth != nil {
			h.RespAuth.Error = errAuth.Error()
			if errors.Is(errAuth, ErrMethod) || errors.Is(errAuth, custerrors.ErrInvalidToken) {
				response.RespJs(w, h.RespAuth, http.StatusUnauthorized)
			} else {
				response.RespJs(w, h.RespAuth, http.StatusInternalServerError)
			}
			return
		}
		response.RespJs(w, authRes, http.StatusOK)
	}
}
func (h *UserHandler) Confirm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tempEmail, ok := r.Context().Value(middleware.KeyTempUser).(string)
		if !ok {
			h.RespAuth.Error = custerrors.ErrInvalidToken.Error()
			response.RespJs(w, h.RespAuth, http.StatusUnauthorized)
			return
		}
		body, errReq := request.RequestHandler[RequestConfirm](r)
		if errReq != nil {
			h.RespConfirm.Error = errReq.Error()
			response.RespJs(w, h.RespConfirm, http.StatusBadRequest)
			return
		}
		sessId := r.Header.Get("Authorization")
		action := r.URL.Query().Get("action")
		if (sessId == "") || (action == "") {
			h.RespConfirm.Error = ErrMissing.Error()
			response.RespJs(w, h.RespConfirm, http.StatusBadRequest)
			return
		}
		resJwt, errConfirm := h.Service.Confirm(body, sessId, tempEmail, action)
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
		h.RespConfirm.Jwt = resJwt
		h.RespConfirm.Error = ""
		response.RespJs(w, h.RespConfirm, http.StatusCreated)
	}
}
