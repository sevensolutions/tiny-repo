package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func ValidateAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Get("user") == nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
		}

		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		name := claims["name"].(string)
		prefix := claims["prefix"].(string)

		if !strings.HasPrefix(prefix, "/") {
			prefix = "/" + prefix
		}

		log.Debug("Username", name)

		path := c.Request().URL.Path

		if strings.HasPrefix(path, prefix) {
			return next(c)
		}

		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthorized to access path " + path})
	}
}
