package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	*Db
	*Auth
}
type Db struct {
	Dsn string
}
type Auth struct{
	Secret string
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
			Secret: os.Getenv("SECRET"),
		},
	}
}
