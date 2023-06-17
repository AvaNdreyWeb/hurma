package main

import (
	"context"
	"hurma-service/hurma-service/handlers"
	"log"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set up MongoDB client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Ping the MongoDB server to ensure a connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.POST("/sign-up", func(c echo.Context) error {
		return handlers.SignUpHandler(c, client)
	})
	e.POST("/login", func(c echo.Context) error {
		return handlers.LoginHandler(c, client)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
