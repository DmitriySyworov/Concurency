package main

import (
	"order/app/configs"
	"order/app/internal/product"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)



func main(){
	conf := configs.NewConfig()
	db, errDb := gorm.Open(postgres.Open(conf.Dsn))
	if errDb != nil {
		panic(errDb)
	}
	errMigrate := db.AutoMigrate(&product.Product{})
	if errMigrate != nil {
		panic(errMigrate)
	}
}