package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ExpireDate struct {
	CreatedAt string `bson:"createdAt"`
	ExpiresAt string `bson:"expiresAt"`
}

type ClickStat struct {
	Total uint64
	Daily []uint64
}

type Link struct {
	Id       primitive.ObjectID `bson:"_id"`
	Title    string
	ShortUrl string `bson:"shortUrl"`
	FullUrl  string `bson:"fullUrl"`
	Expires  ExpireDate
	Clicks   ClickStat
}

type CreateLinkDTO struct {
	Title     string `json:"title"`
	FullUrl   string `json:"fullUrl"`
	CreatedAt string `json:"createdAt"`
	ExpiresAt string `json:"expiresAt"`
}

type EditLinkDTO struct {
	Title     string `json:"title"`
	ExpiresAt string `json:"expiresAt"`
}

type TableLinkDTO struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	ShortUrl    string `json:"shortUrl"`
	ExpiresAt   string `json:"expiresAt"`
	ClicksTotal uint64 `json:"clicksTotal"`
}

type DailyDTO struct {
	Date   string
	Clicks uint64
}
