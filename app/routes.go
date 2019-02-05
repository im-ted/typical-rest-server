package app

import (
	"net/http"

	"github.com/labstack/echo"
)

func initRoutes(server *echo.Echo) {
	server.GET("/", hello)

}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}