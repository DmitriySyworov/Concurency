package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"order/app/internal/common"
	"order/app/internal/product"
	"order/app/internal/test"
	"testing"
)

func TestProductCreate(t *testing.T) {
	db := test.NewDbTest()
	db.InitUser(&test.UserTest)
	defer db.ClearDb()
	data, errJs := json.Marshal(&product.RequestProductCreate{
		Name:        test.ProductFirstTest.Name,
		Description: test.ProductFirstTest.Description,
		Images:      test.ProductFirstTest.Images,
		Category:    test.ProductFirstTest.Category,
	})
	if errJs != nil {
		t.Fatal(errJs)
	}

	request, errReq := http.NewRequest(http.MethodPost, "/users/products", bytes.NewBuffer(data))
	if errReq != nil {
		t.Fatal(errReq)
	}
	request.Header.Add("Authorization", test.UserJwt)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, writer.Code)
	}
	dataResp, errReader := io.ReadAll(writer.Body)
	if errReader != nil {
		t.Fatal(errReader)
	}
	var payload common.Product
	errUn := json.Unmarshal(dataResp, &payload)
	if errUn != nil {
		t.Fatal(errUn)
	}
	if payload.Hash == "" {
		t.Fatal("response empty")
	}
}
func TestProductUpdate(t *testing.T) {
	db := test.NewDbTest()
	db.InitProducts(&common.Product{
		Name:        test.ProductFirstTest.Name,
		Description: test.ProductFirstTest.Description,
		Images:      test.ProductFirstTest.Images,
		Category:    test.ProductFirstTest.Category,
		Hash:        test.ProductFirstTest.Hash,
		UserId:      test.UserTest.UserId,
	})
	defer db.ClearDb()
	data, errJs := json.Marshal(&product.RequestProductUpdate{
		Description: "green cucumber",
		Images:      []string{"https://www.google.com/imgres?q=cucumber", "https://www.google.com/imgres?q=cucumb"},
	})
	if errJs != nil {
		t.Fatal(errJs)
	}

	request, errReq := http.NewRequest(http.MethodPatch, "/users/products/"+test.ProductFirstTest.Hash, bytes.NewBuffer(data))
	if errReq != nil {
		t.Fatal(errReq)
	}
	request.Header.Add("Authorization", test.UserJwt)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, writer.Code)
	}
	dataResp, errReader := io.ReadAll(writer.Body)
	if errReader != nil {
		t.Fatal(errReader)
	}
	var payload common.Product
	errUn := json.Unmarshal(dataResp, &payload)
	if errUn != nil {
		t.Fatal(errUn)
	}
	if payload.Hash == "" {
		t.Fatal("response empty")
	}
	log.Println("initial description: ", test.ProductFirstTest.Description)
	log.Println("final description: ", payload.Description)
	log.Println("initial images: ", test.ProductFirstTest.Images)
	log.Println("final images: ", payload.Images)
}
func TestProductGet(t *testing.T) {
	db := test.NewDbTest()
	db.InitProducts(&common.Product{
		Name:        test.ProductFirstTest.Name,
		Description: test.ProductFirstTest.Description,
		Images:      test.ProductFirstTest.Images,
		Category:    test.ProductFirstTest.Category,
		Hash:        test.ProductFirstTest.Hash,
		UserId:      test.UserTest.UserId,
	})
	defer db.ClearDb()
	request, errReq := http.NewRequest(http.MethodGet, "/users/products/"+test.ProductFirstTest.Hash, nil)
	if errReq != nil {
		t.Fatal(errReq)
	}
	request.Header.Add("Authorization", test.UserJwt)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, writer.Code)
	}
	dataResp, errReader := io.ReadAll(writer.Body)
	if errReader != nil {
		t.Fatal(errReader)
	}
	var payload common.Product
	errUn := json.Unmarshal(dataResp, &payload)
	if errUn != nil {
		t.Fatal(errUn)
	}
	if payload.Hash == "" {
		t.Fatal("response empty")
	}
}
func TestProductDelete(t *testing.T) {
	db := test.NewDbTest()
	db.InitProducts(&common.Product{
		Name:        test.ProductFirstTest.Name,
		Description: test.ProductFirstTest.Description,
		Images:      test.ProductFirstTest.Images,
		Category:    test.ProductFirstTest.Category,
		Hash:        test.ProductFirstTest.Hash,
		UserId:      test.UserTest.UserId,
	})
	defer db.ClearDb()
	request, errReq := http.NewRequest(http.MethodDelete, "/users/products/"+test.ProductFirstTest.Hash, nil)
	if errReq != nil {
		t.Fatal(errReq)
	}
	request.Header.Add("Authorization", test.UserJwt)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusNoContent {
		t.Fatalf("expected %d got %d", http.StatusNoContent, writer.Code)
	}
}
func TestProductAllGet(t *testing.T) {
	db := test.NewDbTest()
	db.InitUser(&test.UserTest)
	db.InitProducts(&common.Product{
		Name:        test.ProductFirstTest.Name,
		Description: test.ProductFirstTest.Description,
		Images:      test.ProductFirstTest.Images,
		Category:    test.ProductFirstTest.Category,
		Hash:        test.ProductFirstTest.Hash,
		UserId:      test.UserTest.UserId,
	})
	defer db.ClearDb()
	request, errReq := http.NewRequest(http.MethodGet, "/users/products", nil)
	if errReq != nil {
		t.Fatal(errReq)
	}
	request.Header.Add("Authorization", test.UserJwt)
	//!
	query := request.URL.Query()
	query.Add("category", "Foods")
	request.URL.RawQuery = query.Encode()
	//!
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, writer.Code)
	}

	dataResp, errReader := io.ReadAll(writer.Body)
	if errReader != nil {
		t.Fatal(errReader)
	}
	var payload product.ResponseSliceProduct
	errUn := json.Unmarshal(dataResp, &payload)
	if errUn != nil {
		t.Fatal(errUn)
	}
	if len(payload.CategoryProduct) == 0 {
		t.Fatal("response empty")
	}
}
