package user

import (
	"fmt"
	"net/smtp"
	"order/app/configs"

	"github.com/jordan-wright/email"
)

func Send(conf *configs.VerifyEmail, emailUser string, sessPassword int) error {
	e := email.NewEmail()
	e.From = conf.EmailApi
	e.To = []string{emailUser}
	e.Subject = "Verification message"
	e.Text = []byte(fmt.Sprint(sessPassword))
	e.HTML = []byte(fmt.Sprintf(`Your code: %d Enter this code in the application or website to gain access to your account. If you did not request a code, simply ignore this email.`, sessPassword))
	errSend := e.Send(conf.AddressHost, smtp.PlainAuth("", conf.EmailApi, conf.PasswordApi, conf.Address))
	if errSend != nil {
		return ErrSendEmail
	}
	return nil
}
