package user

import (
	"fmt"
	"order/app/configs"
	generaterand "order/app/pkg/generateRand"
	jwts "order/app/pkg/jwt"
	"strings"
	"time"

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

func (s *UserService) Register(body *RequestUserRegist) (*TempUser, error) {
	_, errReg := s.GetByEmailOrPhone(body.Email, body.Phone)
	if errReg == nil {
		return nil, ErrReg
	}
	hashPass, errPass := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if errPass != nil {
		return nil, ErrSecurity
	}
	var idUser string
	for {
		idUser = generaterand.RandNumberStr(12)
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
		IdUser:    idUser,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	errTemppUser := s.CreateTempUser(tempUser)
	if errTemppUser != nil {
		return nil, errTemppUser
	}
	jwtToken, errJwt := jwts.NewJWT(s.Secret).TemporaryJWT(&jwts.DataJWt{Email: tempUser.Email, Phone: tempUser.Phone, IdUser: tempUser.IdUser})
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
		IdUser:    user.IdUser,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	errTemppUser := s.CreateTempUser(tempUser)
	if errTemppUser != nil {
		return nil, errTemppUser
	}
	jwtToken, errJwt := jwts.NewJWT(s.Secret).TemporaryJWT(&jwts.DataJWt{Email: tempUser.Email, Phone: tempUser.Phone, IdUser: tempUser.IdUser})
	if errJwt != nil {
		return nil, errJwt
	}
	tempUser.Jwt = jwtToken
	return tempUser, nil
}
func (s *UserService) Auth(method, userToken string) (*ResponseAuth, error) {
	sessID := generaterand.RandStr(10)
	tempPass := generaterand.RandNumberStr(9)
	session := &Session{
		SessionId:    sessID,
		TempPassword: tempPass,
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	}
	user, errAssearch := s.assearchUser(userToken)
	if errAssearch != nil {
		return nil, errAssearch
	}
	switch method {
	case "phone":
		errSess := s.Session(session)
		if errSess != nil {
			return nil, ErrSecurity
		}
		return &ResponseAuth{
			SessionId: sessID,
			Message:   fmt.Sprintf("we sent an email with a password to the specified phone number: %s", user.Phone) + fmt.Sprint("    Это имитация отправки через телефон кода:", tempPass),
		}, nil
	case "email":
		errSess := s.Session(session)
		if errSess != nil {
			return nil, ErrSecurity
		}
		errSend := Send(s.Config.VerifyEmail, user.Email, tempPass)
		if errSend != nil {
			return nil, errSend
		}
		return &ResponseAuth{
			SessionId: sessID,
			Message:   fmt.Sprintf("we have sent an email with a password to the specified email: %s", user.Email),
		}, nil
	default:
		return nil, ErrMethod
	}
}
func (s *UserService) Confirm(body *RequestConfirm, session, userToken, action string) (string, error) {
	switch action {
	case "login":
		_, jwtT, errCompare := s.compareSession(body, session, userToken)
		if errCompare != nil {
			return "", errCompare
		}
		return jwtT, nil
	case "register":
		user, jwtT, errCompare := s.compareSession(body, session, userToken)
		if errCompare != nil {
			return "", errCompare
		}
		_, errRegist := s.GetByEmailOrPhone(user.Email, user.Phone)
		if errRegist == nil {
			return "", ErrReg
		}
		errCreate := s.CreateUser(user)
		if errCreate != nil {
			return "", ErrCreateUser
		}
		return jwtT, nil
	default:
		return "", ErrAuth
	}
}
func (s *UserService) compareSession(body *RequestConfirm, session, userToken string) (*User, string, error) {
	sess, errSess := s.GetSession(strings.TrimPrefix(session, "Bearer "))
	if errSess != nil {
		return nil, "", ErrSess
	}
	if sess.TempPassword != body.TempPassword {
		return nil, "", ErrIncorrectCode
	}
	user, errAssearch := s.assearchUser(userToken)
	if errAssearch != nil {
		return nil, "", errAssearch
	}
	jwtToken, errToken := jwts.NewJWT(s.Secret).CreateJWT(&jwts.DataJWt{IdUser: user.IdUser})
	if errToken != nil {
		return nil, "", errToken
	}
	return user, jwtToken, nil
}
func (s *UserService) assearchUser(userToken string) (*User, error) {
	usJwt, errJwt := jwts.NewJWT(s.Secret).DecodeJWT(strings.TrimPrefix(userToken, "Bearer "))
	if errJwt != nil {
		return nil, errJwt
	}
	var userJwt TempJWTUser
	userJwt = TempJWTUser{
		Email: usJwt.Email,
		Phone: usJwt.Phone,
	}
	user, errGet := s.GetTempUser(userJwt.Email, userJwt.Phone)
	if errGet != nil {
		return nil, errGet
	}
	resUser := &User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Phone:    user.Phone,
		IdUser:   user.IdUser,
	}
	return resUser, nil
}
