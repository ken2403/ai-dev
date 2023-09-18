package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var e *echo.Echo

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		// message = fmt.Sprintf("%v", he.Message)
	}
	c.Logger().Error(err)

	if err := c.JSON(code, map[string]string{"custom": fmt.Sprintf("%v", message)}); err != nil {
		c.Logger().Error(err)
	}
}

func main() {
	e = echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler

	e.Use(middleware.Logger())
	e.GET("/map", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Hello",
		})
	})
	e.GET("/str", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello")
	})
	e.GET("/string", func(c echo.Context) error {
		return echo.NewHTTPError(http.StatusCreated)
	})

	log.Fatal(e.Start(":8080"))
}
