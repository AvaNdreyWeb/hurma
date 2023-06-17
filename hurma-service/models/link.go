package models

type ExpireDate struct {
	Created_at string
	Expires_at string
}

type ClickStat struct {
	Total uint64
	Daily []uint64
}

type Link struct {
	Id        string
	Title     string
	Short_url string
	Full_url  string
	Expires   ExpireDate
	Clicks    ClickStat
}
