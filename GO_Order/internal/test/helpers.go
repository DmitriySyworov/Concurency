package test

import (
	"log"
	"order/app/internal/auth"
	"order/app/internal/common"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbTest struct {
	*gorm.DB
}

var UserTest = common.User{
	Name:     "DDAy",
	Email:    "Dday1@gmail.com",
	Password: "$2a$10$LNAtNw8kSDPFeotw44F9UOkalXNjIlurXpG0UoOtyKiNzqCvmvqZq",
	Phone:    "7121321313",
	UserId:   777463957883,
}
var UserTempTest = auth.TempUser{
	Name:     "DDAy",
	Email:    "Dday1@gmail.com",
	Password: "$2a$10$abjhCeIrA0emF5d1VDCs1OZ98hoxHy1RnCo3qzSPLyrptHn9MmPKW",
	Phone:    "7121321313",
	UserId:   777463957883,
}
var ProductFirstTest = common.Product{
	Name:        "cucumber",
	Description: "big",
	Images:      []string{"https://www.google.com/imgres?q=pro"},
	Category:    "Foods",
	Hash:        "ObgZ8BPp",
	UserId:      478829231251,
}
var ProductSecondTest = common.Product{
	Name:        "videocard",
	Description: "very power",
	Images:      []string{"https://www.google.com/imgres?q=prsaaao"},
	Category:    "Electronics",
	Hash:        "v0MikWOS",
	UserId:      134595623139,
}

const (
	UserJwt  = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZF91c2VyIjo3Nzc0NjM5NTc4ODN9.GiOV_GZDQ5iF9DNQXpdaQYL85ih5vqiNwPQkm4xdbtQ"
	OrderId  = "918cdc94-021b-4eb3-80c3-5c4d241dab03"
	Password = "qwjaixzmx1w"
)

func NewDbTest() *DbTest {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic(errEnv)
	}
	db, errDb := gorm.Open(postgres.Open(os.Getenv("DSN")))
	if errDb != nil {
		panic(errDb)
	}
	return &DbTest{
		DB: db, //.Begin(),
	}
}
func (db *DbTest) InitProducts(product *common.Product) {
	res := db.Create(product)
	if res.Error != nil {
		log.Println("error initialization product")
	}
}
func (db *DbTest) InitUser(user *common.User) {
	res := db.Create(&user)
	if res.Error != nil {
		log.Println("error initialization auth")
	}
}
func (db *DbTest) InitTempUser(user *auth.TempUser) {
	res := db.Create(&user)
	if res.Error != nil {
		log.Println("error initialization tempUser")
	}
}

func (db *DbTest) InitOrder(order *common.Order) {
	res := db.Create(&order)
	if res.Error != nil {
		log.Println(res.Error)
	}
}
func (db *DbTest) InitSession(session *auth.Session) {
	res := db.Create(&session)
	if res.Error != nil {
		log.Println("error init session")
	}
}
func (db *DbTest) InitSoftDeleteUser(user *common.User) {
	res1 := db.Create(&user)
	if res1.Error != nil {
		log.Println("error soft-delete user")
	}
	res2 := db.Where("user_id = ?", user.UserId).Delete(&common.User{})
	if res2.Error != nil {
		log.Println("error initialization soft-delete user")
	}
}
func (db *DbTest) DropTableAndMigrate() {
	db.Migrator().DropTable(&common.Order{})
	db.Migrator().DropTable(&common.Product{})
	db.Migrator().DropTable(&common.User{})
	db.AutoMigrate(&common.Order{})
	db.AutoMigrate(&common.Product{})
	db.AutoMigrate(&common.User{})
}
func (db *DbTest) ClearDb() {
	db.Exec("DELETE FROM sessions")
	db.Exec("DELETE FROM temp_users")
	db.Exec("DELETE FROM order_products")
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM users")
}
