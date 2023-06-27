package handlers

import (
	"hurma/internal/crud"
	"hurma/internal/models"
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

	links, err := um.GetLinks(authUserEmail, page, cl)
	if err != nil {
		if err == crud.ErrUserNotFound {
			r = ResponseJSON{
				Code:    http.StatusNotFound,
				Message: "User not found",
			}
			return c.JSON(http.StatusNotFound, r)
		}
		if err == crud.ErrPageNotFound {
			r = ResponseJSON{
				Code:    http.StatusNotFound,
				Message: "Page not found",
			}
			return c.JSON(http.StatusNotFound, r)
		}
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	user, err := um.Get(authUserEmail, cl)
	if err != nil {
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	data := models.UserLinksDTO{
		Total: len(user.Links),
		Data:  links,
	}
	return c.JSON(http.StatusOK, data)
}
