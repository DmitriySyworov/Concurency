package main

import (
	"order/app/configs"
	"order/app/internal/common"
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
	errMigrate := db.AutoMigrate(&user.User{}, &common.Order{}, &common.Product{}, &user.Session{}, &user.TempUser{})
	if errMigrate != nil {
		panic(errMigrate)
	}
}
