package handlers

import (
	"hurma-service/hurma-service/crud"
	"hurma-service/hurma-service/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoginHandler(c echo.Context, cl *mongo.Client) error {
	u := new(models.AuthUserDTO)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	// um := new(crud.UserManager)

	// if err := um.Validate(u, cl); err != nil {
	// 	return c.String(http.StatusUnauthorized, "invalid username or password")
	// }

	return c.JSON(http.StatusOK, u)
}

func SignUpHandler(c echo.Context, cl *mongo.Client) error {
	u := new(models.AuthUserDTO)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	um := new(crud.UserManager)
	if err := um.Create(u, cl); err != nil {
		if err == crud.ErrUsernameConflict {
			return c.String(http.StatusConflict, crud.ErrUsernameConflict.Error())
		}
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, u)
}
