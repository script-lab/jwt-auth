package main

import (
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/form3tech-oss/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/script-lab/jwt-auth/database"
	"github.com/script-lab/jwt-auth/model"
)

const SecretKey = "secret"

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

func login(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}
	email := c.QueryParam("email")
	database.Mysql.Where("email = ?", email).First(&user)

	if user.ID == 0 {
		return c.JSON(http.StatusBadRequest, user)
	}

	comparePassword := c.QueryParam("password")
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(comparePassword)); err != nil {
		return err
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		return err
	}

	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour * 24)
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"token":   token,
		"message": "success",
	})
}

func authenticate(c echo.Context) error {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return err
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return err
	}

	claims := token.Claims.(*jwt.StandardClaims)

	user := model.User{}

	database.Mysql.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(http.StatusOK, user)
}

func main() {
	e := echo.New()

	// cors
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

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
	e.POST("/signUp", signUp)
	e.POST("/login", login)
	e.GET("/authenticate", authenticate)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
