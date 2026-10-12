package handler

import (
	"context"
	"encoding/base64"
	"net/url"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/nekto-sns/server/model"
)

type scratchAuthService interface{
	LoginCallback(context.Context, string) (string, error)
}

type sessionManager interface{
	Cookie(string) *http.Cookie
}

type scratchAuthHandler struct{
	svc scratchAuthService
	sessionMN sessionManager
	loginRedirectURL string
}

func NewScratchAuth(svc scratchAuthService, sessionMN sessionManager, scratchAuthURL, loginCallbackURL string) (*scratchAuthHandler) {
	base64URL := base64.URLEncoding.EncodeToString([]byte(loginCallbackURL))
	redirect, _ := url.JoinPath(scratchAuthURL, "/auth")

	return &scratchAuthHandler{
		svc: svc,
		loginRedirectURL: redirect + "?redirect=" + base64URL,
		sessionMN: sessionMN,
	}
}

func (h *scratchAuthHandler) LoginRedirect(c *echo.Context) error {
	return c.Redirect(303, h.loginRedirectURL)
}

func (h *scratchAuthHandler) LoginCallback(c *echo.Context) error {
	privateCode := c.QueryParam("privateCode")
	if privateCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "privateCode not found")
	}

	req := c.Request().Context()

	session, err := h.svc.LoginCallback(req, privateCode)
	if err != nil {
		if errors.Is(err, model.ErrInternal) {
			return echo.NewHTTPError(http.StatusInternalServerError, "Login failed").Wrap(err)
		}
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized").Wrap(err)
	}

	c.SetCookie(h.sessionMN.Cookie(session))

	return c.JSON(http.StatusOK, map[string]string{"success": "true"})
}
