package middleware

import (
	"net/http"
	"context"
	"errors"

	"github.com/labstack/echo/v5"

	"github.com/nekto-sns/server/model"
)

type sessionManager interface{
	Cookie(string) *http.Cookie
}

type sessionService interface{
	Update(context.Context, string) (string, bool, error)
}

func NewAuth(m sessionManager, s sessionService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			cookie, err := c.Cookie("session")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication is required").Wrap(err)
			}

			ctx := c.Request().Context()

			userID, isUpdated, err := s.Update(ctx, cookie.Value)

			if err != nil {
				if errors.Is(err, model.ErrNotFound) {
					return echo.NewHTTPError(http.StatusUnauthorized, "Invalid session")
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}

			if isUpdated {
				c.SetCookie(m.Cookie(cookie.Value))
			}

			c.Set("user", userID)

			return next(c)
		}
	}
}
