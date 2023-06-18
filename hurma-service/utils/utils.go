package utils

import (
	"hurma-service/hurma-service/config"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const salt int = 10
const shortLen int = 6

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), salt)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ShortenURL() string {
	cfg := config.GetService()
	rand.Seed(time.Now().UnixNano())

	gen := make([]string, shortLen)
	for i := 0; i < shortLen; i++ {
		gen[i] = string(randomChar())
	}

	addrPart := cfg.Host
	genPart := strings.Join(gen, "")

	shortUrl := strings.Join([]string{addrPart, genPart}, "/")
	return shortUrl
}

func randomChar() byte {
	switch rand.Intn(3) {
	case 0:
		return randomInRange('A', 'Z')
	case 1:
		return randomInRange('a', 'z')
	default:
		return randomInRange('0', '9')
	}
}

func randomInRange(start, end byte) byte {
	return start + byte(rand.Intn(int(end-start+1)))
}
