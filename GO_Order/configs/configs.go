package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	*Db
}
type Db struct {
	Dsn string
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
	}
}
