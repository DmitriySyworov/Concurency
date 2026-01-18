package verify

import (
	"fmt"
	"net/http"
	"net/smtp"
	"verif/app/configs"
	"verif/app/internal/payload"
	"verif/app/internal/service/file"
	requesthandl "verif/app/pkg/requestHandl"
	responsejs "verif/app/pkg/responseJS"

	"github.com/jordan-wright/email"
)

type VerifyHandler struct {
	Conf VerifyDep
}
type VerifyDep struct {
	configs.Configs
}

func NewVeryHandler(router *http.ServeMux, conf VerifyDep) {
	verify := &VerifyHandler{
		Conf: conf,
	}
	router.HandleFunc("POST /send", verify.Send())
	router.HandleFunc("GET /verify/{hash}", verify.Verify())
}
func (vr *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ver, errVer := requesthandl.RequestHandle[payload.Verification](w, r)
		if errVer != nil {
			responsejs.RespJs(w, payload.NewResponseSend(errVer), http.StatusBadRequest)
			return
		}
		resVerify := payload.NewVerification(ver.Email)
		e := email.NewEmail()
		e.From = vr.Conf.EmailApi
		e.To = []string{resVerify.Email}
		e.Subject = "Verification message"
		e.Text = []byte(resVerify.Hash)
		e.HTML = []byte(fmt.Sprintf(`<a href="http://localhost:8081/verify/%s"> Кликните, чтобы подтвердить вход в аккаунт</a>`, resVerify.Hash))
		errSend := e.Send(vr.Conf.AddressHost, smtp.PlainAuth("", vr.Conf.EmailApi, vr.Conf.Password, vr.Conf.Address))
		if errSend != nil {
			responsejs.RespJs(w, payload.NewResponseSend(errSend), http.StatusInternalServerError)
			return
		}
		errCreate := file.CreateJsFile(resVerify)
		if errCreate != nil {
			responsejs.RespJs(w, payload.NewResponseSend(errCreate), http.StatusInternalServerError)
		}
		responsejs.RespJs(w, payload.NewResponseSend(nil), http.StatusCreated)
	}
}
func (vr *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		errhash := file.CheckHash(hash)
		if errhash != nil{
			responsejs.RespJs(w, payload.NewResponseVerify(errhash), http.StatusBadRequest)
			return
		}
		responsejs.RespJs(w, payload.NewResponseVerify(nil), http.StatusOK)
	}
}
