package main

import (
	"net/http"
	"order/app/configs"
	"order/app/internal/product"
	"order/app/internal/user"
	"order/app/pkg/db"
	"order/app/pkg/middleware"
)

func main(){
	conf := configs.NewConfig()
	DbConnect := db.NewDb(conf)
	router := http.NewServeMux()
	productService := product.NewProductService(product.NewProductRepository(DbConnect))
	userService := user.NewUserService(user.NewUserRepository(DbConnect),  conf)
	product.NewProductHandler(router, product.ProductHandlerDep{ProductService: productService})
	user.NewUserHandler(router, &user.UserHandlerDep{UserService: userService})
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	server := http.Server{
		Addr: ":8081",
		Handler: stack(router),
	}
	errApi := server.ListenAndServe()
	if errApi != nil {
		panic(errApi)
	}
}
