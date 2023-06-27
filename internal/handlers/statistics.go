package handlers

import (
	"hurma/internal/config"
	"hurma/internal/models"
	"hurma/internal/utils"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func OneLinkStatisticsHandler(c echo.Context, cl *mongo.Client) error {
	authUserEmail := c.Get("user").(string)
	genPart := c.Param("genPart")
	period := c.QueryParam("period")
	days, err := strconv.Atoi(period)
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	cfg := config.GetService()
	addrPart := cfg.Host
	shortUrl := strings.Join([]string{addrPart, genPart}, "/")
	link, err := lm.GetByShortUrl(shortUrl, cl)
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	id, err := primitive.ObjectIDFromHex(link.Id)
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}
	if !um.StatisticsAccess(authUserEmail, id, cl) {
		r = ResponseJSON{
			Code:    http.StatusForbidden,
			Message: "Permission denied",
		}
		return c.JSON(http.StatusForbidden, r)
	}

	data, err := lm.GetLinkStatistics(link, days, cl)
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	return c.JSON(http.StatusOK, data)
}

func AllLinksStatisticsHandler(c echo.Context, cl *mongo.Client) error {
	authUserEmail := c.Get("user").(string)
	period := c.QueryParam("period")
	days, err := strconv.Atoi(period)
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

	dataList := make([][]models.DailyDTO, 0)
	for _, id := range user.Links {
		link := lm.GetByID(id, cl)
		data, err := lm.GetLinkStatistics(link, days, cl)
		if err != nil {
			log.Println(err.Error())
			r = ResponseJSON{
				Code:    http.StatusInternalServerError,
				Message: "Internal Server Error",
			}
			return c.JSON(http.StatusInternalServerError, r)
		}
		dataList = append(dataList, data)
	}

	linksData := utils.MergeStatistics(dataList)
	return c.JSON(http.StatusOK, linksData)
}
