package main

import (
	"net/http"
	"order/app/configs"
	"order/app/internal/order"
	"order/app/internal/product"
	"order/app/internal/user"
	"order/app/pkg/db"
	"order/app/pkg/middleware"
)

func main() {
	app := App()
	server := http.Server{
		Addr:    ":8081",
		Handler: app,
	}
	errApi := server.ListenAndServe()
	if errApi != nil {
		panic(errApi)
	}
}
func App() http.Handler {
	conf := configs.NewConfig()
	DbConnect := db.NewDb(conf)
	router := http.NewServeMux()
	userRepo := user.NewUserRepository(DbConnect)
	go userRepo.DeleteTempDB()
	//repositories
	productRepo := product.NewProductRepository(DbConnect)
	orderRepo := order.NewRepositoryOrder(DbConnect)
	//services
	productService := product.NewProductService(productRepo, userRepo)
	userService := user.NewUserService(userRepo, conf)
	orderService := order.NewServiceOrder(orderRepo, &order.ServiceOrderDep{IProductRepo: productRepo, IUserRepo: userRepo})
	//handlers
	product.NewHandlerProduct(router, product.HandlerProductDep{Service: productService, Config: conf})
	user.NewUserHandler(router, &user.UserHandlerDep{Service: userService, Config: conf})
	order.NewHandlerOrder(router, &order.HandlerOrderDep{Config: conf, Service: orderService})
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	return stack(router)
}
