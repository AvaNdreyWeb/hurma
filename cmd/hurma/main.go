package main

import (
	"hurma/internal/config"
	"hurma/internal/crud"
	"hurma/internal/handlers"
	"hurma/internal/mw"

	_ "hurma/docs"

	"github.com/robfig/cron/v3"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

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
	config.Init()
	defer config.Close()

	cron := cron.New()
	cron.AddFunc("0 0 0 * * *", crud.UpdateAll)
	cron.Start()

	e := echo.New()
	e.Use(mw.CORS)
	// Auth handlers
	e.POST("/sign-up", handlers.SignUpHandler)
	e.POST("/login", handlers.LoginHandler)
	// Links handlers
	e.GET("/:genPart", handlers.RedirectHandler)
	e.GET("/links", handlers.UserLinksHandler, mw.JwtMiddleware)
	e.POST("/create", handlers.CreateLinkHandler, mw.JwtMiddleware)
	e.PATCH("/edit/:linkId", handlers.EditLinkHandler, mw.JwtMiddleware)
	e.DELETE("/delete/:linkId", handlers.DeleteLinkHandler, mw.JwtMiddleware)
	// Statistics handlers
	e.GET("/statistics", handlers.AllLinksStatisticsHandler, mw.JwtMiddleware)
	e.GET("/statistics/:genPart", handlers.OneLinkStatisticsHandler, mw.JwtMiddleware)
	// Profile handlers
	e.GET("/profile", handlers.ProfileHandler, mw.JwtMiddleware)
	e.POST("/subscribe", handlers.SubscribeHandler, mw.JwtMiddleware)
	e.POST("/unsubscribe", handlers.UnsubscribeHandler, mw.JwtMiddleware)
	// Documentation
	e.GET("/docs*", echoSwagger.WrapHandler)

	cfg := config.Get().Server
	addr := cfg.GetAddr()
	e.Logger.Fatal(e.Start(addr))
}
