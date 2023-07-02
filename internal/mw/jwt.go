package mw

import (
	"fmt"
	"hurma/internal/handlers"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var jwtSigningKey = []byte("secret")

func JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("hurmaToken")
		if err != nil {
			r := handlers.ResponseJSON{
				Code:    http.StatusUnauthorized,
				Message: "Missing or invalid JWT token",
			}
			return c.JSON(http.StatusUnauthorized, r)
		}
		tokenString := cookie.Value

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return jwtSigningKey, nil
		})

		if err != nil || !token.Valid {
			r := handlers.ResponseJSON{
				Code:    http.StatusUnauthorized,
				Message: "Invalid JWT token",
			}
			return c.JSON(http.StatusUnauthorized, r)
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Set("user", claims["user"])

		return next(c)
	}
}
