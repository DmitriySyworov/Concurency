package main

import (
	"net/http"
	"verif/app/configs"
	"verif/app/internal/handlers/verify"
)

func main() {
	conf := configs.LoadEnv()
	router := http.NewServeMux()
	verify.NewVeryHandler(router, verify.VerifyDep{Configs: *conf})
	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	errApi := server.ListenAndServe()
	if errApi != nil {
		server.ErrorLog.Fatal(errApi)
	}
}
