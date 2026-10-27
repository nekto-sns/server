package handler

import (
	"context"
	"encoding/base64"
	"net/url"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/nekto-sns/server/domain"
)

type scratchAuthService interface{
	LoginCallback(context.Context, string) (string, error)
	RegisterCallback(context.Context, string) (string, error)
}

type sessionManager interface{
	Cookie(string) *http.Cookie
}

type scratchAuthHandler struct{
	svc scratchAuthService
	sessionMN sessionManager
	loginRedirectURL string
	registerRedirectURL string
}

func NewScratchAuth(svc scratchAuthService, sessionMN sessionManager, scratchAuthURL, loginCallbackURL, registerCallbackURL string) (*scratchAuthHandler) {
	base64Login    := base64.URLEncoding.EncodeToString([]byte(loginCallbackURL))
	base64Register := base64.URLEncoding.EncodeToString([]byte(registerCallbackURL))
	redirect, _    := url.JoinPath(scratchAuthURL, "/auth")

	return &scratchAuthHandler{
		svc: svc,
		loginRedirectURL:    redirect + "?redirect=" + base64Login,
		registerRedirectURL: redirect + "?redirect=" + base64Register,
		sessionMN:           sessionMN,
	}
}

func (h *scratchAuthHandler) LoginRedirect(c *echo.Context) error {
	return c.Redirect(303, h.loginRedirectURL)
}

func (h *scratchAuthHandler) RegisterRedirect(c *echo.Context) error {
	return c.Redirect(303, h.registerRedirectURL)
}

func (h *scratchAuthHandler) LoginCallback(c *echo.Context) error {
	privateCode := c.QueryParam("privateCode")
	if privateCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "privateCode not found")
	}

	req := c.Request().Context()

	session, err := h.svc.LoginCallback(req, privateCode)
	if err != nil {
		if errors.Is(err, domain.ErrInternal) {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error").Wrap(err)
		}
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	c.SetCookie(h.sessionMN.Cookie(session))

	return c.JSON(http.StatusOK, map[string]string{"success": "true"})
}

func (h *scratchAuthHandler) RegisterCallback(c *echo.Context) error {
	privateCode := c.QueryParam("privateCode")
	if privateCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "privateCode not found")
	}

	req := c.Request().Context()

	session, err := h.svc.RegisterCallback(req, privateCode)
	if err != nil {
		if errors.Is(err, domain.ErrInternal) {
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error").Wrap(err)
		} else if errors.Is(err, domain.ErrConflict) {
			return echo.NewHTTPError(http.StatusConflict, "Account already exists")
		}
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	c.SetCookie(h.sessionMN.Cookie(session))

	return c.JSON(http.StatusOK, map[string]string{"success": "true"})
}
