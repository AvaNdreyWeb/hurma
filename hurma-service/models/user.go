package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id           string
	Password     string
	Email        string
	ChatId       string
	Links        []primitive.ObjectID
	Subscription bool
}

type AuthUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
