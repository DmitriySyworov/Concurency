package verify

import (
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/smtp"
	"verif/app/configs"

	"github.com/jordan-wright/email"
)

type VerifyHandler struct {
	Email string`json:"email"`
	Conf VerifyDep
}
type VerifyDep struct{
	configs.Configs
}
func NewVeryHandler(router *http.ServeMux, email string, conf VerifyDep) {
	verify := &VerifyHandler{
		Email: email,
		Conf: conf,
	}
	router.HandleFunc("POST /send", verify.Send())
	router.HandleFunc("GET /verify/{hash}", verify.Verify())
}
func (vr *VerifyHandler) Send() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e := email.NewEmail()
		e.To = []string{vr.Email}
		e.Subject = "Verification message"
		e.Text = generateLetter()
		errSend := e.Send(vr.Conf.AddressHost, smtp.PlainAuth("", vr.Conf.EmailApi, vr.Conf.Password, vr.Conf.Address))
		if errSend != nil{
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(500)
			fmt.Fprint(w, errSend)
			return 
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(201)
		fmt.Fprintln(w, "The letter was sent successfully")
	}
}
func generateLetter() []byte {
	i := 0
	resLetter := make([]byte, 9)
	for 9 > i {
		randomNum := rand.IntN(58)
		if randomNum > 47 {
			resLetter = append(resLetter, byte(randomNum))
			i++
		}
	}
	return resLetter
}
func (vr *VerifyHandler) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}