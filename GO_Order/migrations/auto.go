package main

import (
	"order/app/configs"
	"order/app/internal/auth"
	"order/app/internal/common"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	conf := configs.NewConfig()
	db, errDb := gorm.Open(postgres.Open(conf.Dsn))
	if errDb != nil {
		panic(errDb)
	}
	errMigrate := db.AutoMigrate(&common.User{}, &common.Order{}, &common.Product{}, &auth.Session{}, &auth.TempUser{})
	if errMigrate != nil {
		panic(errMigrate)
	}
}
