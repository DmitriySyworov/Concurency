package main

import (
	"order/app/configs"
	"order/app/internal/product"
	"order/app/internal/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	conf := configs.NewConfig()
	db, errDb := gorm.Open(postgres.Open(conf.Dsn))
	if errDb != nil {
		panic(errDb)
	}
	errMigrate := db.AutoMigrate(&product.Product{}, &user.User{}, &user.Session{})
	if errMigrate != nil {
		panic(errMigrate)
	}
}
