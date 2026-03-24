package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"order/app/internal/common"
	"order/app/internal/order"
	"order/app/internal/test"
	"testing"
)

var hash = []string{"ObgZ8BPp", "v0MikWOS"}

func TestCreateOrder(t *testing.T) {
	db := test.NewDbTest()
	db.DropTableAndMigrate()
	db.InitProducts(&test.ProductFirstTest)
	db.InitProducts(&test.ProductSecondTest)
	db.InitUser(&test.UserTest)
	defer db.ClearDb()
	data, errJs := json.Marshal(&order.CreateOrderRequest{ProductsHash: hash})
	if errJs != nil {
		t.Fatal(errJs)
	}
	request, errReq := http.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(data))
	if errReq != nil {
		t.Fatal(errReq)
	}
	request.Header.Add("Authorization", test.UserJwt)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, writer.Code)
	}
	body, errRead := io.ReadAll(writer.Body)
	if errRead != nil {
		t.Fatal(errRead)
	}
	var payload common.Order
	errUnmarshal := json.Unmarshal(body, &payload)
	if errUnmarshal != nil {
		t.Fatal(errUnmarshal)
	}
	if payload.Products == nil {
		t.Fatal("response empty")
	}
}
func TestGetOrder(t *testing.T) {
	db := test.NewDbTest()
	db.DropTableAndMigrate()
	db.InitProducts(&test.ProductFirstTest)
	db.InitProducts(&test.ProductSecondTest)
	db.InitUser(&test.UserTest)
	db.InitOrder(&common.Order{
		Products: []common.Product{test.ProductFirstTest, test.ProductSecondTest},
		UserId:   test.UserTempTest.UserId,
		OrderId:  test.OrderId,
	})
	defer db.ClearDb()
	request := httptest.NewRequest(http.MethodGet, "/order/"+test.OrderId, nil)
	request.Header.Add("Authorization", test.UserJwt)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, writer.Code)
	}
	body, errRead := io.ReadAll(writer.Body)
	if errRead != nil {
		t.Fatal(errRead)
	}
	var payload common.Order
	errUn := json.Unmarshal(body, &payload)
	if errUn != nil {
		t.Fatal(errUn)
	}
	if len(payload.Products) == 0 {
		t.Fatal("empty order")
	}
}

func TestAllOrderGet(t *testing.T) {
	db := test.NewDbTest()
	db.DropTableAndMigrate()
	db.InitProducts(&test.ProductFirstTest)
	db.InitProducts(&test.ProductSecondTest)
	db.InitUser(&test.UserTest)
	db.InitOrder(&common.Order{
		Products: []common.Product{test.ProductFirstTest, test.ProductSecondTest},
		UserId:   test.UserTempTest.UserId,
		OrderId:  test.OrderId,
	})
	defer db.ClearDb()
	request := httptest.NewRequest(http.MethodGet, "/my-orders", nil)
	request.Header.Add("Authorization", test.UserJwt)
	writer := httptest.NewRecorder()
	App().ServeHTTP(writer, request)
	if writer.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, writer.Code)
	}
	body, errRead := io.ReadAll(writer.Body)
	if errRead != nil {
		t.Fatal(errRead)
	}
	var payload order.AllOrdersResponse
	errUn := json.Unmarshal(body, &payload)
	if errUn != nil {
		t.Fatal(errUn)
	}
	if len(payload.MyOrders) == 0 {
		t.Fatal("empty orders")
	}
}
