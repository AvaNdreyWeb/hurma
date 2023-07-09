package utils

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"hurma/internal/config"
	"hurma/internal/models"
	"log"
)

type SearchCache struct {
	Email string `json:"email"`
	Page  int    `json:"page"`
}

func GetHashKey(obj SearchCache) (string, error) {
	s, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	hash := sha1.New()
	hash.Write(s)
	hashedBytes := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)
	return hashedString, nil
}

func Stringify(obj models.UserLinksDTO) (string, error) {
	s, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func ClearCachedPages(email string, from, to int) error {
	log.Println("$$$ Clear cache starts:", "from page -", from, "to page (include) -", to)
	for i := from; i <= to; i++ {
		cachedPage := SearchCache{
			Email: email,
			Page:  i,
		}
		log.Println("cleaning:", email, "page -", i)
		key, err := GetHashKey(cachedPage)
		if err != nil {
			return err
		}
		config.Clients.Redis.Del(context.TODO(), key)
	}
	log.Println("$$$ Clear cache ends")
	return nil
}
