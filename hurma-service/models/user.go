package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           string
	Password     string
	Email        string
	ChatId       string `bson:"chatId"`
	Links        []primitive.ObjectID
	Subscription bool
}

type AuthUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLinksDTO struct {
	Total int            `json:"total"`
	Data  []TableLinkDTO `json:"data"`
}
