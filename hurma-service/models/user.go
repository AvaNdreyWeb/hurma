package models

type User struct {
	Id           string
	Username     string
	Password     string
	Email        string
	ChatId       string
	Links        []string
	Subscription bool
}

type AuthUserDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
