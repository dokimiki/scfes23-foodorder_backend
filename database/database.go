package database

import (
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	godotenv.Load()
	var USER = os.Getenv("SCFES23FOODORDER_DB_USER")
	var PASS = os.Getenv("SCFES23FOODORDER_DB_PASS")
	var PROTOCOL = os.Getenv("SCFES23FOODORDER_DB_PROTOCOL")
	var DATABASE = "scfes23"

	dsn := USER + ":" + PASS + "@" + PROTOCOL + "/" + DATABASE + "?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{SkipDefaultTransaction: false})
	if err != nil {
		panic(err)
	}
}
