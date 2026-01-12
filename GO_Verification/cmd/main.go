package main

import (
	"flag"
	"net/http"
	"verif/app/configs"
	"verif/app/internal/handlers/verify"
)

func main() {
	conf := configs.LoadEnv()
	email := flag.String("email", "", "enter your email")
	flag.Parse()
	router := http.NewServeMux()
	verify.NewVeryHandler(router, *email, verify.VerifyDep{Configs: *conf})
	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	errApi := server.ListenAndServe()
	if errApi != nil {
		server.ErrorLog.Fatal(errApi)
	}
}
