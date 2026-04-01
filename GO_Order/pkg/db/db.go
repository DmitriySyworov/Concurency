package db

import (
	"order/app/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct{
	*gorm.DB
}
func NewDb(conf *configs.Config)*Db{
db, errDb := gorm.Open(postgres.Open(conf.Dsn))
if errDb != nil {
	panic(errDb)
}
return &Db{db}
}