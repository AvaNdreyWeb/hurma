package handlers

import (
	"hurma/internal/crud"
	"hurma/internal/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func SubscribeHandler(c echo.Context, cl *mongo.Client) error {
	authUserEmail := c.Get("user").(string)

	um := new(crud.UserManager)
	err := um.Subscribe(authUserEmail, cl)
	if err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "Successfully subscribed")
}

func UnsubscribeHandler(c echo.Context, cl *mongo.Client) error {
	authUserEmail := c.Get("user").(string)

	um := new(crud.UserManager)
	err := um.Unsubscribe(authUserEmail, cl)
	if err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "Successfully subscribed")
}

func ProfileHandler(c echo.Context, cl *mongo.Client) error {
	authUserEmail := c.Get("user").(string)
	um := new(crud.UserManager)
	user, err := um.Get(authUserEmail, cl)
	if err != nil {
		log.Fatal(err)
	}

	u := models.ProfileUserDTO{
		Email:        user.Email,
		ChatId:       user.ChatId,
		Subscription: user.Subscription,
	}

	return c.JSON(http.StatusOK, u)
}
