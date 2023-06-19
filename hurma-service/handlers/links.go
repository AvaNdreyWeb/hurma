package handlers

import (
	"hurma-service/hurma-service/config"
	"hurma-service/hurma-service/crud"
	"hurma-service/hurma-service/models"
	"log"
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

func EditLinkHandler(c echo.Context, cl *mongo.Client) error {
	linkId, err := primitive.ObjectIDFromHex(c.Param("linkId"))
	if err != nil {
		log.Fatal(err)
	}
	l := new(models.EditLinkDTO)
	if err := c.Bind(l); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	lm := new(crud.LinkManager)
	if l.Title != "" {
		err = lm.EditTitle(l.Title, linkId, cl)
		if err != nil {
			if err == crud.ErrLinkConflict {
				return c.String(http.StatusConflict, err.Error())
			}
			log.Fatal(err)
		}
	}
	if l.ExpiresAt != "" {
		err = lm.EditExpires(l.ExpiresAt, linkId, cl)
		if err != nil {
			if err == crud.ErrLinkConflict {
				return c.String(http.StatusConflict, err.Error())
			}
			log.Fatal(err)
		}
	}

	return c.JSON(http.StatusOK, l)
}

func DeleteLinkHandler(c echo.Context, cl *mongo.Client) error {
	authUserEmail := c.Get("user").(string)
	linkId, err := primitive.ObjectIDFromHex(c.Param("linkId"))
	if err != nil {
		log.Fatal(err)
	}

	lm := new(crud.LinkManager)
	err = lm.Delete(authUserEmail, linkId, cl)
	if err != nil {
		if err == crud.ErrLinkConflict {
			return c.String(http.StatusConflict, err.Error())
		}
		log.Fatal(err)
	}

	return c.String(http.StatusOK, "Link successfully deleted")
}

func RedirectHandler(c echo.Context, cl *mongo.Client) error {

	genPart := c.Param("genPart")
	cfg := config.GetService()
	addrPart := cfg.Host
	shortUrl := strings.Join([]string{addrPart, genPart}, "/")

	lm := new(crud.LinkManager)
	link, err := lm.GetFullUrl(shortUrl, cl)
	if err != nil {
		log.Fatal(err)
	}

	err = lm.IncTotal(link.Id, cl)
	if err != nil {
		if err == crud.ErrLinkConflict {
			return c.String(http.StatusConflict, err.Error())
		}
		log.Fatal(err)
	}

	return c.Redirect(http.StatusOK, link.FullUrl)
}
