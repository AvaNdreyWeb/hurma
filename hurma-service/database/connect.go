package database

import (
	"context"
	"hurma-service/hurma-service/config"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDb() *mongo.Client {
	cfg := config.GetMongoDb()
	uri := strings.Join([]string{cfg.Protocol, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Options}, "")
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
