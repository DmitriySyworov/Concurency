package user

import (
	"fmt"
	generaterand "order/app/pkg/generateRand"
	jwts "order/app/pkg/jwt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	*UserRepository
	Secret []byte
}

func NewUserService(repo *UserRepository, secret []byte) *UserService {
	return &UserService{
		UserRepository: repo,
		Secret:         secret,
	}
}

func (s *UserService) Register(body *RequestUserRegist) (*User, error) {
	hashPass, errPass := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if errPass != nil {
		return nil, ErrSecurity
	}
	user := &User{
		Name:     body.Name,
		Email:    body.Email,
		Phone:    body.Phone,
		Password: string(hashPass),
	}
	jwtToken, errJwt := jwts.NewJWT(s.Secret).FastJWT(user)
	if errJwt != nil {
		return nil, errJwt
	}
	user.Jwt = jwtToken
	return user, nil
}
func (s *UserService) Login(body *RequestUserLogin) (*User, error) {
	user, errGet := s.GetUser(body.Email, body.Phone)
	if errGet != nil {
		return nil, ErrWrongData
	}
	errPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if errPass != nil {
		return nil, ErrWrongData
	}
	jwtToken, errJwt := jwts.NewJWT(s.Secret).FastJWT(user)
	if errJwt != nil {
		return nil, errJwt
	}
	user.Jwt = jwtToken
	return user, nil
}
func (s *UserService) Auth(method, userToken string) (*ResponseAuth, error) {
	sessID := generaterand.RandStr(10)
	tempPass := generaterand.RandStr(6)
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
			Message:   fmt.Sprintf("we sent an email with a password to the specified phone number: %s", user.Phone),
		}, nil
	case "email":
		errSess := s.Session(session)
		if errSess != nil {
			return nil, ErrSecurity
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
	jwtToken, errToken := jwts.NewJWT(s.Secret).CreateJWT(user.Email, user.Phone)
	if errToken != nil {
		return nil, "", errToken
	}
	return user, jwtToken, nil
}
func (s *UserService) assearchUser(userToken string) (*User, error) {
	usToken, errJwt := jwts.NewJWT(s.Secret).DecodeJWT(strings.TrimPrefix(userToken, "Bearer "))
	if errJwt != nil {
		return nil, errJwt
	}
	data := usToken["user"]
	var user *User
	switch u := data.(type) {
	case map[string]any:
		user = &User{
			Email:    u["email"].(string),
			Phone:    u["phone"].(string),
			Password: u["password"].(string),
			Name:     u["name"].(string),
		}
	}
	return user, nil
}
