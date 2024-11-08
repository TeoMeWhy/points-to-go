package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func OpenDBConnection() (*gorm.DB, error) {

	godotenv.Load(".env")
	HOST_DB := os.Getenv("HOST_DB")
	PORT_DB := os.Getenv("PORT_DB")
	USER_DB := os.Getenv("USER_DB")
	PASSWORD_DB := os.Getenv("PASSWORD_DB")

	dsn := "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	dsn = fmt.Sprintf(dsn, USER_DB, PASSWORD_DB, HOST_DB, PORT_DB, "points")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}
