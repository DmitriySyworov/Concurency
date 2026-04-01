package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	*Db
	*Auth
	*VerifyEmail
}
type Db struct {
	Dsn string
}
type Auth struct{
	Secret []byte
}
type VerifyEmail struct{
	EmailApi string
	PasswordApi string
	AddressHost string
	Address string
}

func NewConfig() *Config{
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic(errEnv)
	}
	return &Config{
		Db: &Db{
			Dsn: os.Getenv("DSN"),
		},
		Auth: &Auth{
			Secret: []byte(os.Getenv("SECRET")),
		},
		VerifyEmail: &VerifyEmail{
			PasswordApi: os.Getenv("PASSWORD"),
			EmailApi: os.Getenv("EMAIL"),
			AddressHost: os.Getenv("ADDRESS_HOST"),
			Address: os.Getenv("ADDRESS"),
		},
	}
}
