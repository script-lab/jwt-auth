package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/script-lab/jwt-auth/database"
	"github.com/script-lab/jwt-auth/model"
)

func handler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello Go!")
}

func signUp(c echo.Context) error {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return err
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	password := string(hash)

	user := model.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: password,
	}

	database.Mysql.Create(&user)
	return c.JSON(http.StatusCreated, user)
}

func main() {
	e := echo.New()

	// Database
	database.Connect()
	database.Mysql.AutoMigrate(&model.User{})
	if database.Mysql.Migrator().HasTable(&model.User{}) == false {
		database.Mysql.Migrator().CreateTable(&model.User{})
	}
	sqlDB, _ := database.Mysql.DB()
	defer sqlDB.Close()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routing
	e.GET("/", handler)
	e.POST("/signUp", signUp)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
