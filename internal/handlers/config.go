package handlers

import "hurma/internal/crud"

type ResponseJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type tokenDTO struct {
	AccessToken string `json:"accessToken"`
}

var r ResponseJSON

var um crud.UserManager
var lm crud.LinkManager
