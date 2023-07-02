package config

import (
	"context"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Protocol string
	User     string
	Password string
	Host     string
	Port     string
	Options  string
}

type ServerConfig struct {
	Host string
	Port string
}

type ServiceConfig struct {
	Host string
}

type RedisConfig struct {
	Host string
	Port string
}

type AppConfig struct {
	MongoDB *MongoConfig
	Server  *ServerConfig
	Service *ServiceConfig
	Redis   *RedisConfig
}

type AppClients struct {
	MongoDB *mongo.Client
	Redis   *redis.Client
}

var app *AppConfig
var Clients *AppClients

func Init() {
	app = &AppConfig{
		MongoDB: setMongoConfig(),
		Server:  setServerConfig(),
		Service: setServiceConfig(),
		Redis:   setRedisConfig(),
	}
	mongoClient := app.MongoDB.Client()
	redisClient := app.Redis.Client()

	Clients = &AppClients{
		MongoDB: mongoClient,
		Redis:   redisClient,
	}
}

func Close() {
	Clients.MongoDB.Disconnect(context.Background())
	if err := Clients.Redis.Close(); err != nil {
		log.Println(err.Error())
	}
}

func Get() *AppConfig {
	return app
}

func setMongoConfig() *MongoConfig {
	return &MongoConfig{
		Protocol: "mongodb://",
		User:     "admin",
		Password: ":password",
		Host:     "@mongo",
		Port:     ":27017",
		Options:  "",
	}
}

func setServerConfig() *ServerConfig {
	return &ServerConfig{
		Host: "0.0.0.0",
		Port: ":8080",
	}
}

func setServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		Host: "hur.ma",
	}
}

func setRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host: "redis",
		Port: ":6379",
	}
}

func (cfg *MongoConfig) Client() *mongo.Client {
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

func (cfg *RedisConfig) Client() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: strings.Join([]string{cfg.Host, cfg.Port}, ""),
	})
}

func (cfg *ServerConfig) GetAddr() string {
	return strings.Join([]string{cfg.Host, cfg.Port}, "")
}
