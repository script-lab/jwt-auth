package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func handler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello Go!")
}

func main() {
	e := echo.New()

	// Routing
	e.GET("/", handler)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
