package config

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

func GetMongoDb() *MongoConfig {
	return &MongoConfig{
		Protocol: "mongodb://",
		User:     "admin",
		Password: ":password",
		Host:     "@mongo",
		Port:     ":27017",
		Options:  "",
	}
}

func GetServer() *ServerConfig {
	return &ServerConfig{
		Host: "0.0.0.0",
		Port: ":8080",
	}
}

func GetService() *ServiceConfig {
	return &ServiceConfig{
		Host: "hur.ma",
	}
}
