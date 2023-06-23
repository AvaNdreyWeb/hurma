package handlers

import (
	"hurma/internal/crud"
	"hurma/internal/models"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoginHandler(c echo.Context, cl *mongo.Client) error {
	u := new(models.AuthUserDTO)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	um := new(crud.UserManager)

	if err := um.Validate(u, cl); err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = u.Email

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate JWT token")
	}

	return c.String(http.StatusOK, tokenString)
}

func SignUpHandler(c echo.Context, cl *mongo.Client) error {
	u := new(models.AuthUserDTO)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	um := new(crud.UserManager)
	if err := um.Create(u, cl); err != nil {
		if err == crud.ErrEmailConflict {
			return c.String(http.StatusConflict, err.Error())
		}
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, u)
}
