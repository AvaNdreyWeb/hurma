package main

import (
	"context"
	"hurma/internal/config"
	"hurma/internal/crud"
	"hurma/internal/handlers"
	"hurma/internal/mw"
	"hurma/internal/storage"
	"net/http"
	"regexp"

	"strings"

	_ "hurma/docs"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron/v3"
)

func allowOrigin(origin string) (bool, error) {
	return regexp.MatchString(`^*$`, origin)
}

// @title Hurma API
// @version 1.0
// @description Hurma URL shortener and conversion analysis.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	client := storage.ConnectDb()
	defer client.Disconnect(context.Background())

	cfg := config.GetServer()
	addr := strings.Join([]string{cfg.Host, cfg.Port}, "")

	cron := cron.New()
	cron.AddFunc("0 0 0 * * *", func() {
		crud.UpdateAll(client)
	})
	cron.Start()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc:  allowOrigin,
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderCookie, echo.HeaderAccessControlAllowCredentials},
	}))

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

	e.GET("/statistics", func(c echo.Context) error {
		return handlers.AllLinksStatisticsHandler(c, client)
	}, mw.JwtMiddleware)

	e.GET("/statistics/:genPart", func(c echo.Context) error {
		return handlers.OneLinkStatisticsHandler(c, client)
	}, mw.JwtMiddleware)

	e.GET("/profile", func(c echo.Context) error {
		return handlers.ProfileHandler(c, client)
	}, mw.JwtMiddleware)

	e.GET("/links", func(c echo.Context) error {
		return handlers.UserLinksHandler(c, client)
	}, mw.JwtMiddleware)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(addr))
}
