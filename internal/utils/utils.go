package utils

import (
	"hurma/internal/config"
	"hurma/internal/models"
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
	cfg := config.Get().Service
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

func MergeStatistics(rawData [][]models.DailyDTO) []models.DailyDTO {
	mx := 0
	for _, data := range rawData {
		if len(data) > mx {
			mx = len(data)
		}
	}
	merged := make([]models.DailyDTO, mx)
	for _, data := range rawData {
		shift := mx - len(data)
		for j := mx - 1; j >= 0; j-- {
			if len(data) == mx {
				merged[j].Date = data[j].Date
			}
			if j < len(data) {
				merged[j+shift].Clicks += data[j].Clicks
			}
		}
	}
	return merged
}
