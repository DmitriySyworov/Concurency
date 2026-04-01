package main

import (
	"net/http"
	"order/app/configs"
	"order/app/internal/auth"
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
	//repositories
	authRepo := auth.NewRepositoryAuth(DbConnect)
	userRepo := user.NewUserRepository(DbConnect)
	productRepo := product.NewProductRepository(DbConnect)
	orderRepo := order.NewRepositoryOrder(DbConnect)
	//
	go authRepo.DeleteTempDB()
	//services
	authService := auth.NewAuthService(authRepo, &auth.ServiceAuthDep{IUserRepo: userRepo, Config: conf})
	userService := user.NewServiceUser(userRepo)
	productService := product.NewProductService(productRepo, userRepo)
	orderService := order.NewServiceOrder(orderRepo, &order.ServiceOrderDep{IProductRepo: productRepo, IUserRepo: userRepo})
	//handlers
	auth.NewUserHandler(router, &auth.HandlerAuthDep{Service: authService, Config: conf})
	user.NewHandlerUser(router, &user.HandlerUserDep{ServiceUser: userService, Config: conf})
	product.NewHandlerProduct(router, product.HandlerProductDep{Service: productService, Config: conf})
	order.NewHandlerOrder(router, &order.HandlerOrderDep{Service: orderService, Config: conf})
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	return stack(router)
}
