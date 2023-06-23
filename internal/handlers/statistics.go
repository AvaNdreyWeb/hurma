package handlers

import (
	"hurma/internal/config"
	"hurma/internal/crud"
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
		log.Fatal(err)
	}

	cfg := config.GetService()
	addrPart := cfg.Host
	shortUrl := strings.Join([]string{addrPart, genPart}, "/")
	lm := new(crud.LinkManager)
	link, err := lm.GetByShortUrl(shortUrl, cl)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DONE short url", link.Id)

	um := new(crud.UserManager)
	id, err := primitive.ObjectIDFromHex(link.Id)
	if err != nil {
		log.Fatal(err)
	}
	if !um.StatisticsAccess(authUserEmail, id, cl) {
		return c.String(http.StatusMethodNotAllowed, "Permission denied")
	}

	data, err := lm.GetLinkStatistics(link, days, cl)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, data)
}

func AllLinksStatisticsHandler(c echo.Context, cl *mongo.Client) error {
	authUserEmail := c.Get("user").(string)
	period := c.QueryParam("period")
	days, err := strconv.Atoi(period)
	if err != nil {
		log.Fatal(err)
	}

	lm := new(crud.LinkManager)
	um := new(crud.UserManager)

	user, err := um.Get(authUserEmail, cl)
	if err != nil {
		log.Fatal(err)
	}

	dataList := make([][]models.DailyDTO, 0)
	for _, id := range user.Links {
		link := lm.GetByID(id, cl)
		data, err := lm.GetLinkStatistics(link, days, cl)
		if err != nil {
			log.Fatal(err)
		}
		dataList = append(dataList, data)
	}

	linksData, err := utils.MergeStatistics(dataList)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, linksData)
}
