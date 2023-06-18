package models

type ExpireDate struct {
	CreatedAt string
	ExpiresAt string
}

type ClickStat struct {
	Total uint64
	Daily []uint64
}

type Link struct {
	Id       string
	Title    string
	ShortUrl string
	FullUrl  string
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
