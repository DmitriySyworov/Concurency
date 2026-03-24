package user

import (
	"fmt"
	"net/smtp"
	"order/app/configs"
	generaterand "order/app/pkg/generateRand"
	jwts "order/app/pkg/jwt"
	"strings"
	"time"

	"github.com/jordan-wright/email"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	*UserRepository
	*configs.Config
}

func NewUserService(repo *UserRepository, conf *configs.Config) *UserService {
	return &UserService{
		UserRepository: repo,
		Config:         conf,
	}
}

func (s *UserService) Register(body *RequestUserRegister) (*TempUser, error) {
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
		errId := s.GetByIdUser(idUser)
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
	tempUser.Jwt = jwtToken
	return tempUser, nil
}
func (s *UserService) Login(body *RequestUserLogin) (*TempUser, error) {
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
	errTemppUser := s.CreateTempUser(tempUser)
	if errTemppUser != nil {
		return nil, errTemppUser
	}
	jwtToken, errJwt := jwts.NewJWT(s.Secret).CreateTemporaryJWT(&jwts.DataJWt{Email: tempUser.Email})
	if errJwt != nil {
		return nil, errJwt
	}
	tempUser.Jwt = jwtToken
	return tempUser, nil
}
func (s *UserService) Auth(method, tempEmail string) (*ResponseAuth, error) {
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
		tempUser, errGet := s.GetTempUser(tempEmail)
		if errGet != nil {
			return nil, ErrUserNotFound
		}
		return &ResponseAuth{
			SessionId: sessID,
			Message:   fmt.Sprintf("we sent an email with a password to the specified phone number: %s", tempUser.Phone) + fmt.Sprint("    Это имитация отправки через телефон кода:", tempPass),
		}, nil
	case "email":
		errSess := s.CreateSession(session)
		if errSess != nil {
			return nil, ErrSecurity
		}
		errSend := send(s.Config.VerifyEmail, tempEmail, tempPass)
		if errSend != nil {
			return nil, errSend
		}
		return &ResponseAuth{
			SessionId: sessID,
			Message:   fmt.Sprintf("we have sent an email with a password to the specified email: %s", tempEmail),
		}, nil
	default:
		return nil, ErrMethod
	}
}
func (s *UserService) Confirm(body *RequestConfirm, session, tempEmail, action string) (string, error) {
	switch action {
	case "login":
		_, jwtT, errCompare := s.compareSession(body, session, tempEmail)
		if errCompare != nil {
			return "", errCompare
		}
		return jwtT, nil
	case "register":
		tempUser, jwtT, errCompare := s.compareSession(body, session, tempEmail)
		if errCompare != nil {
			return "", errCompare
		}
		_, errRegist := s.GetByEmailOrPhone(tempUser.Email, tempUser.Phone)
		if errRegist == nil {
			return "", ErrReg
		}
		errCreate := s.CreateUser(&User{
			Email:    tempUser.Email,
			Phone:    tempUser.Phone,
			Name:     tempUser.Name,
			Password: tempUser.Password,
			UserId:   tempUser.UserId,
		})
		if errCreate != nil {
			return "", ErrCreateUser
		}
		return jwtT, nil
	default:
		return "", ErrAuth
	}
}
func send(conf *configs.VerifyEmail, emailUser string, sessPassword int) error {
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
func (s *UserService) compareSession(body *RequestConfirm, session, tempEmail string) (*TempUser, string, error) {
	sess, errSess := s.GetSession(strings.TrimPrefix(session, "Bearer "))
	if errSess != nil {
		return nil, "", ErrSess
	}
	if sess.TempPassword != body.TempPassword {
		return nil, "", ErrIncorrectCode
	}
	tempUser, errGet := s.GetTempUser(tempEmail)
	if errGet != nil {
		return nil, "", ErrUserNotFound
	}
	jwtToken, errToken := jwts.NewJWT(s.Secret).CreateJWT(&jwts.DataJWt{IdUser: float64(tempUser.UserId)})
	if errToken != nil {
		return nil, "", errToken
	}
	return tempUser, jwtToken, nil
}
