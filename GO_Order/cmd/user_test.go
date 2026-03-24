package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"order/app/configs"
	"order/app/internal/test"
	"order/app/internal/user"
	"order/app/pkg/jwt"
	"testing"
	"time"
)

var caseUserLoginSuccess = []struct {
	Name string
	User user.RequestUserLogin
}{
	{Name: "Email", User: user.RequestUserLogin{Name: "DDAy", Email: "Dday1@gmail.com", Password: "o1uh-d9saaijxmzozxu=12-saa"}},
	{Name: "Phone", User: user.RequestUserLogin{Name: "DDAy", Phone: "7121321313", Password: "o1uh-d9saaijxmzozxu=12-saa"}},
}

func TestLoginSuccess(t *testing.T) {
	db := test.NewDbTest()
	db.InitUser(&test.UserTest)
	defer db.ClearDb()
	ts := httptest.NewServer(App())
	defer ts.Close()
	for _, tests := range caseUserLoginSuccess {
		data, _ := json.Marshal(&tests.User)
		resp, err := http.Post(ts.URL+"/user/login", "application/json", bytes.NewBuffer(data))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected %d got %d", http.StatusOK, resp.StatusCode)
		}
		dataResp, errBody := io.ReadAll(resp.Body)
		if errBody != nil {
			t.Fatal(errBody)
		}
		var tempUser user.TempUser
		errJs := json.Unmarshal(dataResp, &tempUser)
		if errJs != nil {
			t.Fatal(errJs)
		}
		if tempUser.Jwt == "" {
			t.Fatal("token not transferred")
		}
	}
}

var caseUserLoginNegative = []struct {
	Name string
	User user.RequestUserLogin
}{
	{Name: "No_Password", User: user.RequestUserLogin{Name: "DDAy", Email: "Dday1@gmail.com"}},
	{Name: "No_Name", User: user.RequestUserLogin{Phone: "7121321313", Password: "o1uh-d9saaijxmzozxu=12-saa"}},
	{Name: "No_Email_And_Phone", User: user.RequestUserLogin{Name: "DDAy", Password: "o1uh-d9saaijxmzozxu=12-saa"}},
}

func TestLoginNegative(t *testing.T) {
	ts := httptest.NewServer(App())
	defer ts.Close()
	for _, tests := range caseUserLoginNegative {
		data, _ := json.Marshal(&tests)
		resp, err := http.Post(ts.URL+"/user/login", "application/json", bytes.NewBuffer(data))
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected %d got %d", http.StatusBadRequest, resp.StatusCode)
		}
	}
}

func TestRegisterSuccess(t *testing.T) {
	db := test.NewDbTest()
	defer db.ClearDb()
	ts := httptest.NewServer(App())
	defer ts.Close()
	data, errJs := json.Marshal(&user.RequestUserRegister{
		Name:     test.UserTest.Name,
		Email:    test.UserTest.Email,
		Password: "1",
		Phone:    test.UserTest.Phone,
	})
	if errJs != nil {
		t.Fatal(errJs)
	}
	resp, errResp := http.Post(ts.URL+"/user/regist", "application/json", bytes.NewBuffer(data))
	if errResp != nil {
		t.Fatal(errResp)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, resp.StatusCode)
	}
	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		t.Fatal(errRead)
	}
	var payload user.TempUser
	errUn := json.Unmarshal(body, &payload)
	if errUn != nil {
		t.Fatal(errUn)
	}
	if payload.Jwt == "" {
		t.Fatal("jwt empty")
	}
}

func TestAuthPhoneSuccess(t *testing.T) {
	db := test.NewDbTest()
	db.InitTempUser(&test.UserTempTest)
	defer db.ClearDb()
	j := jwt.NewJWT(configs.NewConfig().Secret)
	token, errJwt := j.CreateTemporaryJWT(&jwt.DataJWt{Email: test.UserTempTest.Email})
	if errJwt != nil {
		t.Fatal(errJwt)
	}
	request := httptest.NewRequest(http.MethodPost, "/user/auth/"+"phone", nil)
	request.Header.Add("X-User-Token", "Bearer "+token)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, writer.Code)
	}
	body, errRead := io.ReadAll(writer.Body)
	if errRead != nil {
		t.Fatal(errRead)
	}
	var payload user.ResponseAuth
	errUn := json.Unmarshal(body, &payload)
	if errUn != nil {
		t.Fatal(errUn)
	}
	if payload.SessionId == "" {
		t.Fatal("session empty")
	}
}
func TestConfirmRegisterSuccess(t *testing.T) {
	db := test.NewDbTest()
	db.InitTempUser(&test.UserTempTest)
	db.InitSession(&user.Session{
		SessionId:    "T7wnAhtzGX",
		TempPassword: 448688,
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	})
	defer db.ClearDb()
	j := jwt.NewJWT(configs.NewConfig().Secret)
	token, errJwt := j.CreateTemporaryJWT(&jwt.DataJWt{Email: test.UserTempTest.Email})
	if errJwt != nil {
		t.Fatal(errJwt)
	}
	data, errJs := json.Marshal(&user.RequestConfirm{
		TempPassword: 448688,
	})
	if errJs != nil {
		t.Fatal(errJs)
	}
	request := httptest.NewRequest(http.MethodPost, "/user/auth", bytes.NewBuffer(data))
	query := request.URL.Query()
	query.Add("action", "register")
	request.URL.RawQuery = query.Encode()
	request.Header.Add("Authorization", "Bearer T7wnAhtzGX")
	request.Header.Add("X-User-Token", "Bearer "+token)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, writer.Code)
	}
	body, errRead := io.ReadAll(writer.Body)
	if errRead != nil {
		t.Fatal(errRead)
	}
	var payload user.ResponseConfirm
	errUn := json.Unmarshal(body, &payload)
	if errUn != nil {
		t.Fatal(errUn)
	}
	if payload.Jwt == "" {
		t.Fatal("jwt empty")
	}
}
func TestConfirmLoginSuccess(t *testing.T) {
	db := test.NewDbTest()
	db.InitTempUser(&test.UserTempTest)
	db.InitSession(&user.Session{
		SessionId:    "T7wnAhtzGX",
		TempPassword: 448688,
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	})
	db.InitUser(&test.UserTest)
	defer db.ClearDb()
	j := jwt.NewJWT(configs.NewConfig().Secret)
	token, errJwt := j.CreateTemporaryJWT(&jwt.DataJWt{Email: test.UserTempTest.Email})
	if errJwt != nil {
		t.Fatal(errJwt)
	}
	data, errJs := json.Marshal(&user.RequestConfirm{
		TempPassword: 448688,
	})
	if errJs != nil {
		t.Fatal(errJs)
	}
	request := httptest.NewRequest(http.MethodPost, "/user/auth", bytes.NewBuffer(data))
	query := request.URL.Query()
	query.Add("action", "login")
	request.URL.RawQuery = query.Encode()
	request.Header.Add("Authorization", "Bearer T7wnAhtzGX")
	request.Header.Add("X-User-Token", "Bearer "+token)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, writer.Code)
	}
	body, errRead := io.ReadAll(writer.Body)
	if errRead != nil {
		t.Fatal(errRead)
	}
	var payload user.ResponseConfirm
	errUn := json.Unmarshal(body, &payload)
	if errUn != nil {
		t.Fatal(errUn)
	}
	if payload.Jwt == "" {
		t.Fatal("jwt empty")
	}
}
