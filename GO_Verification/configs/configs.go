package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configs struct{
	EmailApi string
	Password string
	Address string
	AddressHost string
}

func LoadEnv()*Configs{
	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("could not read .env file, verification email could not be sent")
	}
	return &Configs{
		EmailApi: os.Getenv("email"),
		Password: os.Getenv("password"),
		Address: os.Getenv("address"),
		AddressHost: os.Getenv("addresshost"),
	}
}