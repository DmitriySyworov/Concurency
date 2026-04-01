package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"order/app/configs"
	"order/app/internal/auth"
	"order/app/internal/test"
	"order/app/pkg/jwt"
	"testing"
	"time"
)

func TestLoginPhoneSuccess(t *testing.T) {
	db := test.NewDbTest()
	db.InitUser(&test.UserTest)
	defer db.ClearDb()
	ts := httptest.NewServer(App())
	defer ts.Close()
	data, _ := json.Marshal(&auth.RequestUserLogin{
		Name:     test.UserTest.Name,
		Phone:    test.UserTest.Phone,
		Password: test.Password,
	})
	resp, err := http.Post(ts.URL+"/auth/login/phone", "application/json", bytes.NewBuffer(data))
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
	var payload auth.ResponseAuth
	errJs := json.Unmarshal(dataResp, &payload)
	if errJs != nil {
		t.Fatal(errJs)
	}
	if payload.Jwt == "" {
		t.Fatal("token not transferred")
	}
}

var caseUserLoginNegative = []struct {
	Name string
	User auth.RequestUserLogin
}{
	{Name: "No_Password", User: auth.RequestUserLogin{Name: "DDAy", Email: "Dday1@gmail.com"}},
	{Name: "No_Name", User: auth.RequestUserLogin{Phone: "7121321313", Password: "o1uh-d9saaijxmzozxu=12-saa"}},
	{Name: "No_Email_And_Phone", User: auth.RequestUserLogin{Name: "DDAy", Password: "o1uh-d9saaijxmzozxu=12-saa"}},
}

func TestLoginNegative(t *testing.T) {
	ts := httptest.NewServer(App())
	defer ts.Close()
	for _, tests := range caseUserLoginNegative {
		data, _ := json.Marshal(&tests)
		resp, err := http.Post(ts.URL+"/auth/login/phone", "application/json", bytes.NewBuffer(data))
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
	data, errJs := json.Marshal(&auth.RequestUserRegister{
		Name:     test.UserTest.Name,
		Email:    test.UserTest.Email,
		Password: "1",
		Phone:    test.UserTest.Phone,
	})
	if errJs != nil {
		t.Fatal(errJs)
	}
	resp, errResp := http.Post(ts.URL+"/auth/register/phone", "application/json", bytes.NewBuffer(data))
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
	var payload auth.ResponseAuth
	errUn := json.Unmarshal(body, &payload)
	if errUn != nil {
		t.Fatal(errUn)
	}
	if payload.Jwt == "" {
		t.Fatal("jwt empty")
	}
}

func TestConfirmRegisterSuccess(t *testing.T) {
	db := test.NewDbTest()
	db.InitTempUser(&test.UserTempTest)
	db.InitSession(&auth.Session{
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
	data, errJs := json.Marshal(&auth.RequestConfirm{
		TempPassword: 448688,
	})
	if errJs != nil {
		t.Fatal(errJs)
	}
	request := httptest.NewRequest(http.MethodPost, "/auth/confirm", bytes.NewBuffer(data))
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
	var payload auth.ResponseConfirm
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
	db.InitSession(&auth.Session{
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
	data, errJs := json.Marshal(&auth.RequestConfirm{
		TempPassword: 448688,
	})
	if errJs != nil {
		t.Fatal(errJs)
	}
	request := httptest.NewRequest(http.MethodPost, "/auth/confirm", bytes.NewBuffer(data))
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
	var payload auth.ResponseConfirm
	errUn := json.Unmarshal(body, &payload)
	if errUn != nil {
		t.Fatal(errUn)
	}
	if payload.Jwt == "" {
		t.Fatal("jwt empty")
	}
}
func TestConfirmRestoreSuccess(t *testing.T) {
	db := test.NewDbTest()
	db.InitTempUser(&test.UserTempTest)
	db.InitSession(&auth.Session{
		SessionId:    "T7wnAhtzGX",
		TempPassword: 448688,
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	})
	db.InitSoftDeleteUser(&test.UserTest)
	defer db.ClearDb()
	j := jwt.NewJWT(configs.NewConfig().Secret)
	token, errJwt := j.CreateTemporaryJWT(&jwt.DataJWt{Email: test.UserTempTest.Email})
	if errJwt != nil {
		t.Fatal(errJwt)
	}
	data, errJs := json.Marshal(&auth.RequestConfirm{
		TempPassword: 448688,
	})
	if errJs != nil {
		t.Fatal(errJs)
	}
	request := httptest.NewRequest(http.MethodPost, "/auth/confirm", bytes.NewBuffer(data))
	query := request.URL.Query()
	query.Add("action", "restore")
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
	var payload auth.ResponseConfirm
	errUn := json.Unmarshal(body, &payload)
	if errUn != nil {
		t.Fatal(errUn)
	}
	if payload.Jwt == "" {
		t.Fatal("jwt empty")
	}
}
