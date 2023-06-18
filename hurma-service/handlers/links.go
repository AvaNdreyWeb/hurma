package handlers

import (
	"hurma-service/hurma-service/crud"
	"hurma-service/hurma-service/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateLinkHandler(c echo.Context, cl *mongo.Client) error {
	l := new(models.CreateLinkDTO)
	if err := c.Bind(l); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	lm := new(crud.LinkManager)
	if err := lm.Create(l, cl); err != nil {
		if err == crud.ErrLinkConflict {
			return c.String(http.StatusConflict, crud.ErrUsernameConflict.Error())
		}
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, l)
}
