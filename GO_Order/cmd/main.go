package main

import (
	"context"
	"net/http"
	"order/app/configs"
	"order/app/internal/order"
	"order/app/internal/product"
	"order/app/internal/user"
	"order/app/pkg/db"
	"order/app/pkg/middleware"
)

func main() {
	conf := configs.NewConfig()
	DbConnect := db.NewDb(conf)
	router := http.NewServeMux()
	ctxCancel, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := user.NewUserRepository(DbConnect)
	go r.DeleteTempDB(ctxCancel)
	//repositories
	productRepo := product.NewProductRepository(DbConnect)
	orderRepo := order.NewRepositoryOrder(DbConnect)
	//services
	productService := product.NewProductService(productRepo)
	userService := user.NewUserService(r, conf)
	orderService := order.NewServiceOrder(orderRepo, &order.ServiceOrderDep{IProductRepo: productRepo})
	//handlers
	product.NewHandlerProduct(router, product.HandlerProductDep{Service: productService, Config: conf})
	user.NewUserHandler(router, &user.UserHandlerDep{Service: userService, Config: conf})
	order.NewHandlerOrder(router, &order.HandlerOrderDep{Config: conf, Service: orderService})
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}
	errApi := server.ListenAndServe()
	if errApi != nil {
		panic(errApi)
	}
}
