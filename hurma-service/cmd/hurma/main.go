package main

import (
	"context"
	"hurma-service/hurma-service/config"
	"hurma-service/hurma-service/database"
	"hurma-service/hurma-service/handlers"
	"strings"

	"github.com/labstack/echo/v4"
)

func main() {
	client := database.ConnectDb()
	defer client.Disconnect(context.Background())

	cfg := config.GetServer()
	addr := strings.Join([]string{cfg.Host, cfg.Port}, "")

	e := echo.New()

	e.POST("/sign-up", func(c echo.Context) error {
		return handlers.SignUpHandler(c, client)
	})

	e.POST("/login", func(c echo.Context) error {
		return handlers.LoginHandler(c, client)
	})

	e.POST("/create", func(c echo.Context) error {
		return handlers.CreateLinkHandler(c, client)
	})

	e.Logger.Fatal(e.Start(addr))
}
