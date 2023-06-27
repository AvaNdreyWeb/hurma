package handlers

import (
	"hurma/internal/crud"
	"hurma/internal/models"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoginHandler(c echo.Context, cl *mongo.Client) error {
	u := new(models.AuthUserDTO)
	if err := c.Bind(u); err != nil {
		r = ResponseJSON{
			Code:    http.StatusBadRequest,
			Message: "Bad request",
		}
		return c.JSON(http.StatusBadRequest, r)
	}

	if err := um.Validate(u, cl); err != nil {
		r = ResponseJSON{
			Code:    http.StatusUnauthorized,
			Message: "Invalid email or password",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user"] = u.Email
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	access := tokenDTO{AccessToken: tokenString}
	return c.JSON(http.StatusOK, access)
}

func SignUpHandler(c echo.Context, cl *mongo.Client) error {
	u := new(models.AuthUserDTO)
	if err := c.Bind(u); err != nil {
		r = ResponseJSON{
			Code:    http.StatusBadRequest,
			Message: "Bad request",
		}
		return c.JSON(http.StatusBadRequest, r)
	}

	if err := um.Create(u, cl); err != nil {
		if err == crud.ErrEmailConflict {
			r = ResponseJSON{
				Code:    http.StatusConflict,
				Message: "User with this email already exists",
			}
			return c.JSON(http.StatusConflict, r)
		}
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	r = ResponseJSON{
		Code:    http.StatusOK,
		Message: "OK",
	}
	return c.JSON(http.StatusOK, r)
}
