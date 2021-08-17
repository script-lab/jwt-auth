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
	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")

	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	user := model.User{
		Name:     name,
		Email:    email,
		Password: passwordHash,
	}

	if err := c.Bind(&user); err != nil {
		return err
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
