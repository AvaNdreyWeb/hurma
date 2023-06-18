package handlers

import (
	"hurma-service/hurma-service/crud"
	"hurma-service/hurma-service/models"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserLinksHandler(c echo.Context, cl *mongo.Client) error {
	authUserEmail := c.Get("user").(string)
	queryPage := c.QueryParam("page")
	page, err := strconv.Atoi(queryPage)
	if err != nil {
		page = 1
	}

	um := new(crud.UserManager)
	links, err := um.GetLinks(authUserEmail, page, cl)
	if err != nil {
		if err == crud.ErrUserNotFound {
			return c.String(http.StatusNotFound, err.Error())
		}
		if err == crud.ErrPageNotFound {
			return c.String(http.StatusNotFound, err.Error())
		}
		log.Fatal(err)
	}

	user, err := um.Get(authUserEmail, cl)
	if err != nil {
		log.Fatal(err)
	}

	data := models.UserLinksDTO{
		Total: len(user.Links),
		Data:  links,
	}

	return c.JSON(http.StatusOK, data)
}
