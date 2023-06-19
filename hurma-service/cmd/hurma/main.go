package main

import (
	"context"
	"hurma-service/hurma-service/config"
	"hurma-service/hurma-service/database"
	"hurma-service/hurma-service/handlers"
	"hurma-service/hurma-service/mw"
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
	}, mw.JwtMiddleware)

	e.PATCH("/edit/:linkId", func(c echo.Context) error {
		return handlers.EditLinkHandler(c, client)
	}, mw.JwtMiddleware)

	e.DELETE("/delete/:linkId", func(c echo.Context) error {
		return handlers.DeleteLinkHandler(c, client)
	}, mw.JwtMiddleware)

	e.POST("/subscribe", func(c echo.Context) error {
		return handlers.SubscribeHandler(c, client)
	}, mw.JwtMiddleware)

	e.POST("/unsubscribe", func(c echo.Context) error {
		return handlers.UnsubscribeHandler(c, client)
	}, mw.JwtMiddleware)

	e.GET("/:genPart", func(c echo.Context) error {
		return handlers.RedirectHandler(c, client)
	})

	// e.GET("/statistics", func(c echo.Context) error {
	// 	return handlers.CreateLinkHandler(c, client)
	// }, mw.JwtMiddleware)

	// e.GET("/statistics/:linkId", func(c echo.Context) error {
	// 	return handlers.CreateLinkHandler(c, client)
	// }, mw.JwtMiddleware)

	e.GET("/profile", func(c echo.Context) error {
		return handlers.ProfileHandler(c, client)
	}, mw.JwtMiddleware)

	e.GET("/links", func(c echo.Context) error {
		return handlers.UserLinksHandler(c, client)
	}, mw.JwtMiddleware)

	e.Logger.Fatal(e.Start(addr))
}
