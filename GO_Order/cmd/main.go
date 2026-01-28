package main

import (
	"net/http"
	"order/app/configs"
	"order/app/internal/product"
	"order/app/pkg/db"
)

func main(){
	conf := configs.NewConfig()
	DbConnect := db.NewDb(conf)
	Db := product.NewProductRepository(DbConnect)
	router := http.NewServeMux()
	product.NewProductHandler(router, product.ProductHandlerDep{ProductRepository: Db})
	server := http.Server{
		Addr: ":8081",
		Handler: router,
	}
	errApi := server.ListenAndServe()
	if errApi != nil {
		panic(errApi)
	}
}
