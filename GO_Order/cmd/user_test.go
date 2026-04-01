package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"order/app/internal/common"
	"order/app/internal/test"
	"order/app/internal/user"
	"testing"
)

func TestUserGet(t *testing.T) {
	db := test.NewDbTest()
	db.InitUser(&test.UserTest)
	defer db.ClearDb()
	request, errReq := http.NewRequest(http.MethodGet, "/users/my", nil)
	if errReq != nil {
		t.Fatal(errReq)
	}
	request.Header.Set("Authorization", test.UserJwt)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, writer.Code)
	}
	data, errRead := io.ReadAll(writer.Body)
	if errRead != nil {
		t.Fatal(errRead)
	}
	var payload common.User
	errJs := json.Unmarshal(data, &payload)
	if errJs != nil {
		t.Fatal(errJs)
	}
	if payload.Name == "" {
		t.Fatal("response empty")
	}
}

func TestUserUpdate(t *testing.T) {
	db := test.NewDbTest()
	db.InitUser(&test.UserTest)
	defer db.ClearDb()
	dataJs, errMarshal := json.Marshal(&user.RequestUpdateUser{
		Name: "newUser",
	})
	if errMarshal != nil {
		t.Fatal(errMarshal)
	}
	request, errReq := http.NewRequest(http.MethodPatch, "/users/my", bytes.NewBuffer(dataJs))
	if errReq != nil {
		t.Fatal(errReq)
	}
	request.Header.Set("Authorization", test.UserJwt)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, writer.Code)
	}
	data, errRead := io.ReadAll(writer.Body)
	if errRead != nil {
		t.Fatal(errRead)
	}
	var payload common.User
	errJs := json.Unmarshal(data, &payload)
	if errJs != nil {
		t.Fatal(errJs)
	}
	if payload.Name != "newUser" || payload.UserId != test.UserTest.UserId {
		t.Fatal("bad response")
	}
}
func TestUserDelete(t *testing.T) {
	db := test.NewDbTest()
	db.InitUser(&test.UserTest)
	defer db.ClearDb()
	request, errReq := http.NewRequest(http.MethodDelete, "/users/my", nil)
	if errReq != nil {
		t.Fatal(errReq)
	}
	request.Header.Set("Authorization", test.UserJwt)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusNoContent {
		t.Fatalf("expected %d got %d", http.StatusNoContent, writer.Code)
	}
}
