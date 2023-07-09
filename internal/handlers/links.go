package handlers

import (
	"context"
	"encoding/json"
	"hurma/internal/config"
	"hurma/internal/crud"
	"hurma/internal/models"
	"hurma/internal/utils"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Create link
// @Description Create new short link
// @Tags Links
// @Param createLinkDTO body models.CreateLinkDTO true "Create link data"
// @Accept json
// @Produce json
// @Success 200 {object} ResponseJSON
// @Failure 400 {object} ResponseJSON
// @Router /create [post]
func CreateLinkHandler(c echo.Context) error {
	cl := config.Clients.MongoDB
	authUserEmail := c.Get("user").(string)
	l := new(models.CreateLinkDTO)
	if err := c.Bind(l); err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusBadRequest,
			Message: "Bad request",
		}
		return c.JSON(http.StatusBadRequest, r)
	}

	linkId, err := lm.Create(l, cl)
	if err != nil {
		log.Println(err.Error())
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
	err = um.AddLink(authUserEmail, linkId, cl)
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}
	user, err := um.Get(authUserEmail, cl)
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}
	pages := int(math.Ceil(float64(len(user.Links)) / 10))
	utils.ClearCachedPages(authUserEmail, 1, pages)

	r = ResponseJSON{
		Code:    http.StatusOK,
		Message: "OK",
	}
	return c.JSON(http.StatusOK, r)
}

// @Summary Edit link
// @Description Edit link title or expire date
// @Tags Links
// @Param editLinkDTO body models.EditLinkDTO true "Edit link data"
// @Param linkId path string true "Link Id"
// @Accept json
// @Produce json
// @Success 200 {object} ResponseJSON
// @Failure 400 {object} ResponseJSON
// @Router /edit/{linkId} [patch]
func EditLinkHandler(c echo.Context) error {
	authUserEmail := c.Get("user").(string)
	cl := config.Clients.MongoDB
	linkId, err := primitive.ObjectIDFromHex(c.Param("linkId"))
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	l := new(models.EditLinkDTO)
	if err := c.Bind(l); err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusBadRequest,
			Message: "Bad request",
		}
		return c.JSON(http.StatusBadRequest, r)
	}

	if l.Title != "" {
		if err = lm.EditTitle(l.Title, linkId, cl); err != nil {
			log.Println(err.Error())
			r = ResponseJSON{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}
			return c.JSON(http.StatusInternalServerError, r)
		}
	}
	if l.ExpiresAt != "" {
		if err = lm.EditExpires(l.ExpiresAt, linkId, cl); err != nil {
			log.Println(err.Error())
			r = ResponseJSON{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}
			return c.JSON(http.StatusInternalServerError, r)
		}
	}

	user, err := um.Get(authUserEmail, cl)
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}
	number := 0
	for i := 0; i < len(user.Links); i++ {
		if user.Links[i] == linkId {
			number = i
			break
		}
	}
	page := 1 + int(math.Floor(float64(number)/10))
	utils.ClearCachedPages(authUserEmail, page, page)

	r = ResponseJSON{
		Code:    http.StatusOK,
		Message: "OK",
	}
	return c.JSON(http.StatusOK, r)
}

// @Summary Delete link
// @Description Delete link with linkId
// @Tags Links
// @Param linkId path string true "Link Id"
// @Produce json
// @Success 200 {object} ResponseJSON
// @Failure 400 {object} ResponseJSON
// @Router /delete/{linkId} [delete]
func DeleteLinkHandler(c echo.Context) error {
	cl := config.Clients.MongoDB
	authUserEmail := c.Get("user").(string)
	linkId, err := primitive.ObjectIDFromHex(c.Param("linkId"))
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	if err = lm.Delete(authUserEmail, linkId, cl); err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}
	user, err := um.Get(authUserEmail, cl)
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}
	number := 0
	for i := 0; i < len(user.Links); i++ {
		if user.Links[i] == linkId {
			number = i
			break
		}
	}
	page := 1 + int(math.Floor(float64(number)/10))
	pages := int(math.Ceil(float64(len(user.Links)) / 10))
	utils.ClearCachedPages(authUserEmail, page, pages)

	r = ResponseJSON{
		Code:    http.StatusOK,
		Message: "OK",
	}
	return c.JSON(http.StatusOK, r)
}

// @Summary Get statistics of all links
// @Description Getting array of clicks by day of all links
// @Tags Links
// @Param page query int true "Current page" minimum(1)
// @Produce json
// @Success 200 {object} models.UserLinksDTO
// @Failure 400 {object} ResponseJSON
// @Router /links [get]
func UserLinksHandler(c echo.Context) error {
	cl := config.Clients.MongoDB
	authUserEmail := c.Get("user").(string)
	queryPage := c.QueryParam("page")
	page, err := strconv.Atoi(queryPage)
	if err != nil {
		page = 1
	}

	q := utils.SearchCache{
		Email: authUserEmail,
		Page:  page,
	}
	hashKey, err := utils.GetHashKey(q)
	if err != nil {
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}
	val, err := config.Clients.Redis.Get(context.TODO(), hashKey).Result()
	if err == nil {
		cache := new(models.UserLinksDTO)
		json.Unmarshal([]byte(val), cache)
		log.Println("FROM CACHE")
		log.Println("Getting from", hashKey, "...")
		return c.JSON(http.StatusOK, cache)
	}

	links, err := um.GetLinks(authUserEmail, page, cl)
	if err != nil {
		log.Println(err.Error())
		if err == crud.ErrUserNotFound {
			r = ResponseJSON{
				Code:    http.StatusNotFound,
				Message: "User not found",
			}
			return c.JSON(http.StatusNotFound, r)
		}
		// if err == crud.ErrPageNotFound {
		// 	r = ResponseJSON{
		// 		Code:    http.StatusNotFound,
		// 		Message: "Page not found",
		// 	}
		// 	return c.JSON(http.StatusNotFound, r)
		// }
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	user, err := um.Get(authUserEmail, cl)
	if err != nil {
		log.Println(err.Error())
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
	s, err := utils.Stringify(data)
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}
	err = config.Clients.Redis.Set(context.TODO(), hashKey, s, 0).Err()
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	log.Println("FROM MONGODB")
	return c.JSON(http.StatusOK, data)
}

// Redirect from short to full url
func RedirectHandler(c echo.Context) error {
	cl := config.Clients.MongoDB
	genPart := c.Param("genPart")
	cfg := config.App.Service
	addrPart := cfg.Host
	shortUrl := strings.Join([]string{addrPart, genPart}, "/")

	link, err := lm.GetFullUrl(shortUrl, cl)
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	id := link.Id
	// if err != nil {
	// 	log.Println(err.Error())
	// 	r = ResponseJSON{
	// 		Code:    http.StatusInternalServerError,
	// 		Message: "Internal Server Error",
	// 	}
	// 	return c.JSON(http.StatusInternalServerError, r)
	// }

	err = lm.IncTotal(id, cl)
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	return c.Redirect(http.StatusPermanentRedirect, link.FullUrl)
}
