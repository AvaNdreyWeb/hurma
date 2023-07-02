package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
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
