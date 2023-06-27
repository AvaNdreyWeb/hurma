package handlers

import (
	"hurma/internal/config"
	"hurma/internal/crud"
	"hurma/internal/models"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateLinkHandler(c echo.Context, cl *mongo.Client) error {
	authUserEmail := c.Get("user").(string)

	l := new(models.CreateLinkDTO)
	if err := c.Bind(l); err != nil {
		r = ResponseJSON{
			Code:    http.StatusBadRequest,
			Message: "Bad request",
		}
		return c.JSON(http.StatusBadRequest, r)
	}

	linkId, err := lm.Create(l, cl)
	if err != nil {
		if err == crud.ErrLinkConflict {
			r = ResponseJSON{
				Code:    http.StatusConflict,
				Message: "Link already exists",
			}
			return c.JSON(http.StatusConflict, r)
		}
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	if err = um.AddLink(authUserEmail, linkId, cl); err != nil {
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

func EditLinkHandler(c echo.Context, cl *mongo.Client) error {
	linkId, err := primitive.ObjectIDFromHex(c.Param("linkId"))
	if err != nil {
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	l := new(models.EditLinkDTO)
	if err := c.Bind(l); err != nil {
		r = ResponseJSON{
			Code:    http.StatusBadRequest,
			Message: "Bad request",
		}
		return c.JSON(http.StatusBadRequest, r)
	}

	if l.Title != "" {
		if err = lm.EditTitle(l.Title, linkId, cl); err != nil {
			r = ResponseJSON{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}
			return c.JSON(http.StatusInternalServerError, r)
		}
	}
	if l.ExpiresAt != "" {
		if err = lm.EditExpires(l.ExpiresAt, linkId, cl); err != nil {
			r = ResponseJSON{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}
			return c.JSON(http.StatusInternalServerError, r)
		}
	}

	r = ResponseJSON{
		Code:    http.StatusOK,
		Message: "OK",
	}
	return c.JSON(http.StatusOK, r)
}

func DeleteLinkHandler(c echo.Context, cl *mongo.Client) error {
	authUserEmail := c.Get("user").(string)
	linkId, err := primitive.ObjectIDFromHex(c.Param("linkId"))
	if err != nil {
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	if err = lm.Delete(authUserEmail, linkId, cl); err != nil {
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

func RedirectHandler(c echo.Context, cl *mongo.Client) error {
	genPart := c.Param("genPart")
	cfg := config.GetService()
	addrPart := cfg.Host
	shortUrl := strings.Join([]string{addrPart, genPart}, "/")

	link, err := lm.GetFullUrl(shortUrl, cl)
	if err != nil {
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	id, err := primitive.ObjectIDFromHex(link.Id)
	if err != nil {
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	err = lm.IncTotal(id, cl)
	if err != nil {
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	return c.Redirect(http.StatusPermanentRedirect, link.FullUrl)
}
