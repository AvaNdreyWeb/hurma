package handlers

import (
	"hurma/internal/config"
	"hurma/internal/crud"
	"hurma/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// @Summary Login user
// @Description Login user and getting JWT in response
// @Tags Auth
// @Param authUserDTO body models.AuthUserDTO true "Login user data"
// @Accept json
// @Produce json
// @Success 200 {object} ResponseJSON
// @Failure 400 {object} ResponseJSON
// @Router /login [post]
func LoginHandler(c echo.Context) error {
	cl := config.Clients.MongoDB
	u := new(models.AuthUserDTO)
	if err := c.Bind(u); err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusBadRequest,
			Message: "Bad request",
		}
		return c.JSON(http.StatusBadRequest, r)
	}

	if err := um.Validate(u, cl); err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusUnauthorized,
			Message: "Invalid email or password",
		}
		return c.JSON(http.StatusUnauthorized, r)
	}
	const hours = 1
	expirationTime := time.Now().Add(hours * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"exp": expirationTime.Unix(), "user": u.Email})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}
		return c.JSON(http.StatusInternalServerError, r)
	}

	cookie := &http.Cookie{
		Name:     "hurmaToken",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
		Domain:   "localhost",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		Expires:  time.Now().Add(hours * time.Hour),
	}
	c.SetCookie(cookie)

	r = ResponseJSON{
		Code:    http.StatusOK,
		Message: "OK",
	}
	return c.JSON(http.StatusOK, r)
}

// @Summary Create new user
// @Description Create new user with email and password
// @Tags Auth
// @Param authUserDTO body models.AuthUserDTO true "Create new user data"
// @Accept json
// @Produce json
// @Success 200 {object} ResponseJSON
// @Failure 400 {object} ResponseJSON
// @Router /sign-up [post]
func SignUpHandler(c echo.Context) error {
	cl := config.Clients.MongoDB
	u := new(models.AuthUserDTO)
	if err := c.Bind(u); err != nil {
		log.Println(err.Error())
		r = ResponseJSON{
			Code:    http.StatusBadRequest,
			Message: "Bad request",
		}
		return c.JSON(http.StatusBadRequest, r)
	}

	if err := um.Create(u, cl); err != nil {
		log.Println(err.Error())
		if err == crud.ErrEmailConflict {
			r = ResponseJSON{
				Code:    http.StatusConflict,
				Message: "User with this email already exists",
			}
			return c.JSON(http.StatusConflict, r)
		}
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
