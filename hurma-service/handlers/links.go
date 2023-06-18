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
	authUserEmail := c.Get("user").(string)

	l := new(models.CreateLinkDTO)
	if err := c.Bind(l); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	lm := new(crud.LinkManager)
	linkId, err := lm.Create(l, cl)
	if err != nil {
		if err == crud.ErrLinkConflict {
			return c.String(http.StatusConflict, err.Error())
		}
		log.Fatal(err)
	}

	um := new(crud.UserManager)
	if err = um.AddLink(authUserEmail, linkId, cl); err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, l)
}
