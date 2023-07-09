package mw

import (
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func allowOrigin(origin string) (bool, error) {
	return regexp.MatchString(`^*$`, origin)
}

var CORS echo.MiddlewareFunc = middleware.CORSWithConfig(middleware.CORSConfig{
	AllowOriginFunc:  allowOrigin,
	AllowCredentials: true,
	AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderCookie, echo.HeaderAccessControlAllowCredentials},
})
