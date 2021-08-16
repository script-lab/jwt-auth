package database

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"os"
)

var Mysql *gorm.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")

	dsn := user + ":" + password + "@tcp(" + "DOCKER_MYSQL" + ":" + port + ")/" + databaseName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}
	Mysql = db
}
