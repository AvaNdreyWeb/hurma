package handlers

import (
	"hurma/internal/config"
	"hurma/internal/crud"

	"go.mongodb.org/mongo-driver/mongo"
)

type ResponseJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// type tokenDTO struct {
// 	AccessToken string `json:"accessToken"`
// }

var r ResponseJSON

var um crud.UserManager
var lm crud.LinkManager

var cl *mongo.Client = config.Clients.MongoDB
