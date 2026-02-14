package user

import (
	"net/http"
	custerrors "order/app/pkg/custErrors"
	"order/app/pkg/request"
	"order/app/pkg/response"
)

type UserHandler struct {
	*UserHandlerDep
	RespUser    User
	RespAuth    ResponseAuth
	RespConfirm ResponseConfirm
}
type UserHandlerDep struct {
	*UserService
}

func NewUserHandler(router *http.ServeMux, dep *UserHandlerDep) {
	user := &UserHandler{
		UserHandlerDep: dep,
	}
	//Важные примечания по эндпоинтам!
	//user/login - здесь пользователь отправляет данные для логина, name. password, email  или phone. Если такой пользователь имеется в базе мы его вытаскиваем и конвертируем в  jwt-токен отправляя пользователю
	router.HandleFunc("POST /user/login", user.HandlerLogin())
	//user/regist - здесь мы создаем нового пользователя, затем как же как при логине сформированного пользовтеля превращаем в  jwt и отправляем
	router.HandleFunc("POST /user/regist", user.HandlerRegister())
	//user/auth/{method} - здесь клиент передает пустое тело, но в заголовке тот самый отправленный ранее  jwt, мы его читаем, извлекая оттуда и почту и телефон,затем  согласно выбранному методу отплавляем верификационный пароль на почту или телефон 
	router.HandleFunc("POST /user/auth/{method}", user.HandlerAuth())
	///user/auth здесь в заголовке Autorization мы указываем идентификацию сессии, в теле пароль присланный на посту или телефон, а в  X-User-Token сформированный в /user/login  или /user/regist jwt-токен
	router.HandleFunc("POST /user/auth", user.HandlerConfirimation())
}

func (h *UserHandler) HandlerRegister() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, errReq := request.RequestHandler[RequestUserRegist](w, r)
		if errReq != nil {
			h.RespUser.Error = errReq.Error()
			response.RespJs(w, h.RespUser, http.StatusBadRequest)
			return
		}
		userReg, errRegister := h.Register(body)
		if errRegister != nil {
			h.RespUser.Error = errRegister.Error()
			response.RespJs(w, h.RespUser, http.StatusInternalServerError)
			return
		}
		response.RespJs(w, userReg, http.StatusCreated)
	}
}
func (h *UserHandler) HandlerLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, errReq := request.RequestHandler[RequestUserLogin](w, r)
		if errReq != nil || (body.Email == "" && body.Phone == "") || (body.Email != "" && body.Phone != "") {
			h.RespUser.Error = custerrors.ErrInvalidData.Error()
			response.RespJs(w, h.RespUser, http.StatusBadRequest)
			return
		}
		userLog, errLog := h.Login(body)
		if errLog != nil {
			h.RespUser.Error = errLog.Error()
			if errLog  == ErrSecurity {
				response.RespJs(w, h.RespUser, http.StatusInternalServerError)
			} else {
				response.RespJs(w, h.RespUser, http.StatusUnauthorized)
			}
			return
		}
		response.RespJs(w, userLog, http.StatusOK)
	}
}
func (h *UserHandler) HandlerAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		method := r.PathValue("method")
		tokenUser := r.Header.Get("X-User-Token")
		if method == "" || tokenUser == ""{
			h.RespAuth.Error = ErrMissing.Error()
			response.RespJs(w, h.RespAuth, http.StatusBadRequest)
			return
		}
		authRes, errAuth := h.Auth(method, tokenUser)
		if errAuth != nil {
			h.RespAuth.Error = errAuth.Error()
			if errAuth == ErrMethod || errAuth == custerrors.ErrInvalidToken{
			response.RespJs(w, h.RespAuth, http.StatusUnauthorized)
			} else {
				response.RespJs(w, h.RespAuth, http.StatusInternalServerError)
			}
			return
		}
		response.RespJs(w, authRes, http.StatusOK)
	}
}
func (h *UserHandler) HandlerConfirimation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, errReq := request.RequestHandler[RequestConfirm](w, r)
		if errReq != nil {
			h.RespConfirm.Error = errReq.Error()
			response.RespJs(w, h.RespConfirm, http.StatusBadRequest)
			return
		}
		sessId := r.Header.Get("Authorization")
		tokenUser := r.Header.Get("X-User-Token")
		action := r.URL.Query().Get("action")
		if (sessId == "")|| (tokenUser == "") || (action == ""){
			h.RespConfirm.Error = ErrMissing.Error()
			response.RespJs(w, h.RespConfirm, http.StatusBadRequest)
			return
		}
		resJwt, errConfirm := h.Confirm(body, sessId, tokenUser, action)
		if errConfirm != nil {
			h.RespConfirm.Error = errConfirm.Error()
			h.RespConfirm.Jwt = ""
			if errConfirm == ErrSecurity || errConfirm == ErrCreateUser{
			response.RespJs(w, h.RespConfirm, http.StatusInternalServerError)
			} else {
				response.RespJs(w, h.RespConfirm, http.StatusUnauthorized)
			}
			return
		}
		h.RespConfirm.Jwt = resJwt
		response.RespJs(w, h.RespConfirm, http.StatusCreated)
	}
}
