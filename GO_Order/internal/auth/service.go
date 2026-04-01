package auth

import (
	"fmt"
	"net/smtp"
	"order/app/configs"
	"order/app/internal/common"
	custerrors "order/app/pkg/custErrors"
	"order/app/pkg/di"
	generaterand "order/app/pkg/generateRand"
	jwts "order/app/pkg/jwt"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"golang.org/x/crypto/bcrypt"
)

type ServiceAuth struct {
	*RepositoryAuth
	*ServiceAuthDep
}
type ServiceAuthDep struct {
	di.IUserRepo
	*configs.Config
}

func NewAuthService(repo *RepositoryAuth, dep *ServiceAuthDep) *ServiceAuth {
	return &ServiceAuth{
		RepositoryAuth: repo,
		ServiceAuthDep: dep,
	}
}

func (s *ServiceAuth) Register(body *RequestUserRegister, method string) (*ResponseAuth, error) {
	_, errReg := s.GetByEmailOrPhone(body.Email, body.Phone)
	if errReg == nil {
		return nil, ErrReg
	}
	hashPass, errPass := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if errPass != nil {
		return nil, ErrSecurity
	}
	var idUser int
	for {
		idUser = generaterand.RandNumber(12)
		_, errId := s.GetByIdUser(idUser)
		if errId != nil {
			break
		}
	}
	tempUser := &TempUser{
		Name:      body.Name,
		Email:     body.Email,
		Phone:     body.Phone,
		Password:  string(hashPass),
		UserId:    idUser,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	errTempUser := s.CreateTempUser(tempUser)
	if errTempUser != nil {
		return nil, errTempUser
	}
	jwtToken, errJwt := jwts.NewJWT(s.Secret).CreateTemporaryJWT(&jwts.DataJWt{Email: tempUser.Email})
	if errJwt != nil {
		return nil, errJwt
	}
	respAuth, errHelper := s.helperAuth(method, tempUser.Email, tempUser.Phone, s.Config)
	if errHelper != nil {
		return nil, errHelper
	}
	respAuth.Jwt = jwtToken
	return respAuth, nil
}
func (s *ServiceAuth) Login(body *RequestUserLogin, method string) (*ResponseAuth, error) {
	user, errGet := s.GetByEmailOrPhone(body.Email, body.Phone)
	if errGet != nil {
		return nil, ErrWrongData
	}
	errPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if errPass != nil {
		return nil, ErrWrongData
	}
	tempUser := &TempUser{
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Password:  user.Password,
		UserId:    user.UserId,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	errTempUser := s.CreateTempUser(tempUser)
	if errTempUser != nil {
		return nil, errTempUser
	}
	jwtToken, errJwt := jwts.NewJWT(s.Secret).CreateTemporaryJWT(&jwts.DataJWt{Email: tempUser.Email})
	if errJwt != nil {
		return nil, errJwt
	}
	tempUser.Jwt = jwtToken
	respAuth, errHelper := s.helperAuth(method, tempUser.Email, tempUser.Phone, s.Config)
	if errHelper != nil {
		return nil, errHelper
	}
	respAuth.Jwt = jwtToken
	return respAuth, nil
}
func (s *ServiceAuth) Restore(body *RequestRestore, method string) (*ResponseAuth, error) {
	user, errGet := s.GetByDeleteUser(body.Email, body.Phone)
	if errGet != nil {
		return nil, custerrors.ErrUserNotFound
	}
	tempUser := &TempUser{
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Password:  user.Password,
		UserId:    user.UserId,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	errTempUser := s.CreateTempUser(tempUser)
	if errTempUser != nil {
		return nil, errTempUser
	}
	jwtToken, errJwt := jwts.NewJWT(s.Secret).CreateTemporaryJWT(&jwts.DataJWt{Email: tempUser.Email})
	if errJwt != nil {
		return nil, errJwt
	}
	tempUser.Jwt = jwtToken
	respAuth, errHelper := s.helperAuth(method, tempUser.Email, tempUser.Phone, s.Config)
	if errHelper != nil {
		return nil, errHelper
	}
	respAuth.Jwt = jwtToken
	return respAuth, nil
}
func (s *ServiceAuth) helperAuth(method, email, phone string, conf *configs.Config) (*ResponseAuth, error) {
	sessID := generaterand.RandStr(10)
	tempPass := generaterand.RandNumber(6)
	session := &Session{
		SessionId:    sessID,
		TempPassword: tempPass,
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	}
	switch method {
	case "phone":
		errSess := s.CreateSession(session)
		if errSess != nil {
			return nil, ErrSecurity
		}
		return &ResponseAuth{
			SessionId: sessID,
			Message:   fmt.Sprintf("we sent an email with a password to the specified phone number: %s", phone) + fmt.Sprint("    Это имитация отправки через телефон кода:", tempPass),
		}, nil
	case "email":
		errSess := s.CreateSession(session)
		if errSess != nil {
			return nil, ErrSecurity
		}
		errSend := send(conf.VerifyEmail, email, tempPass)
		if errSend != nil {
			return nil, errSend
		}
		return &ResponseAuth{
			SessionId: sessID,
			Message:   fmt.Sprintf("we have sent an email with a password to the specified email: %s", email),
		}, nil
	default:
		return nil, ErrMethodAuth
	}
}
func (s *ServiceAuth) Confirm(body *RequestConfirm, session, tempEmail, action string) (*ResponseConfirm, error) {
	switch action {
	case "login":
		_, jwtT, errCompare := s.compareSession(body, session, tempEmail)
		if errCompare != nil {
			return nil, errCompare
		}
		return &ResponseConfirm{Jwt: jwtT}, nil
	case "register":
		tempUser, jwtT, errCompare := s.compareSession(body, session, tempEmail)
		if errCompare != nil {
			return nil, errCompare
		}
		_, errRegister := s.GetByEmailOrPhone(tempUser.Email, tempUser.Phone)
		if errRegister == nil {
			return nil, ErrReg
		}
		errCreate := s.CreateUser(&common.User{
			Email:    tempUser.Email,
			Phone:    tempUser.Phone,
			Name:     tempUser.Name,
			Password: tempUser.Password,
			UserId:   tempUser.UserId,
		})
		if errCreate != nil {
			return nil, ErrCreateUser
		}
		return &ResponseConfirm{Jwt: jwtT}, nil
	case "restore":
		tempUser, jwtT, errCompare := s.compareSession(body, session, tempEmail)
		if errCompare != nil {
			return nil, errCompare
		}
		_, errGet := s.GetByDeleteUser(tempUser.Email, tempUser.Phone)
		if errGet != nil {
			return nil, custerrors.ErrUserNotFound
		}
		fmt.Println(tempUser)
		errRestore := s.RestoreUser(tempUser.UserId)
		if errRestore != nil {
			return nil, ErrRestoreUser
		}
		return &ResponseConfirm{Jwt: jwtT}, nil
	default:
		return nil, ErrMethodAuth
	}
}
func send(conf *configs.VerifyEmail, emailUser string, sessPassword int) error {
	e := email.NewEmail()
	e.From = conf.EmailApi
	e.To = []string{emailUser}
	e.Subject = "Verification message"
	e.Text = []byte(fmt.Sprint(sessPassword))
	e.HTML = []byte(fmt.Sprintf(`Your code: %d Enter this code in the application or website to gain access to your account. If you did not requestJs a code, simply ignore this email.`, sessPassword))
	errSend := e.Send(conf.AddressHost, smtp.PlainAuth("", conf.EmailApi, conf.PasswordApi, conf.Address))
	if errSend != nil {
		return ErrSendEmail
	}
	return nil
}
func (s *ServiceAuth) compareSession(body *RequestConfirm, session, tempEmail string) (*TempUser, string, error) {
	sess, errSess := s.GetSession(strings.TrimPrefix(session, "Bearer "))
	if errSess != nil {
		return nil, "", ErrSess
	}
	if sess.TempPassword != body.TempPassword {
		return nil, "", ErrIncorrectCode
	}
	tempUser, errGet := s.GetTempUser(tempEmail)
	if errGet != nil {
		return nil, "", custerrors.ErrUserNotFound
	}
	jwtToken, errToken := jwts.NewJWT(s.Secret).CreateJWT(&jwts.DataJWt{IdUser: float64(tempUser.UserId)})
	if errToken != nil {
		return nil, "", errToken
	}
	return tempUser, jwtToken, nil
}
