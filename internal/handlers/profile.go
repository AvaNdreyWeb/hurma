package handlers

import (
	"hurma/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func SubscribeHandler(c echo.Context, cl *mongo.Client) error {
	authUserEmail := c.Get("user").(string)

	err := um.Subscribe(authUserEmail, cl)
	if err != nil {
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

func UnsubscribeHandler(c echo.Context, cl *mongo.Client) error {
	authUserEmail := c.Get("user").(string)

	err := um.Unsubscribe(authUserEmail, cl)
	if err != nil {
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

func ProfileHandler(c echo.Context, cl *mongo.Client) error {
	authUserEmail := c.Get("user").(string)
	user, err := um.Get(authUserEmail, cl)
	if err != nil {
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	u := models.ProfileUserDTO{
		Email:        user.Email,
		ChatId:       user.ChatId,
		Subscription: user.Subscription,
	}

	return c.JSON(http.StatusOK, u)
}
