package handlers

import (
	"hurma/internal/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// @Summary Subscribe to statistics
// @Description Subscribe to the user's email statistics
// @Tags Profile
// @Produce json
// @Success 200 {object} ResponseJSON
// @Failure 400 {object} ResponseJSON
// @Router /subscribe [post]
func SubscribeHandler(c echo.Context) error {
	authUserEmail := c.Get("user").(string)
	err := um.Subscribe(authUserEmail, cl)
	if err != nil {
		log.Println(err.Error())
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

// @Summary Unubscribe from statistics
// @Description Unubscribe from the user's email statistics
// @Tags Profile
// @Produce json
// @Success 200 {object} ResponseJSON
// @Failure 400 {object} ResponseJSON
// @Router /unsubscribe [post]
func UnsubscribeHandler(c echo.Context) error {
	authUserEmail := c.Get("user").(string)
	err := um.Unsubscribe(authUserEmail, cl)
	if err != nil {
		log.Println(err.Error())
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

// @Summary Get user info
// @Description Getting email, tg chat id and subscribtion status
// @Tags Profile
// @Produce json
// @Success 200 {object} models.ProfileUserDTO
// @Failure 400 {object} ResponseJSON
// @Router /profile [get]
func ProfileHandler(c echo.Context) error {
	authUserEmail := c.Get("user").(string)
	user, err := um.Get(authUserEmail, cl)
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	u := models.ProfileUserDTO{
		Email:        user.Email,
		ChatId:       user.ChatId,
		Subscription: user.Subscription,
	}

	return c.JSON(http.StatusOK, u)
}
