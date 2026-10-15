package handler

import (
	"time"
	"errors"
	"context"
	"net/http"
	"github.com/labstack/echo/v5"

	"github.com/nekto-sns/server/domain"
	"github.com/nekto-sns/server/model"
)

type userService interface{
	ByName(context.Context, string) (*model.User, error)
	ByID(context.Context, string) (*model.User, error)
}

type userHandler struct{
	svc userService
}

type userJSON struct{
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Bio         string    `json:"bio"`
	CreatedAt   time.Time `json:"created_at"`
}

func userResponse(u *model.User) *userJSON {
	return &userJSON{
		ID:          u.ID,
		Name:        u.Name,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
		CreatedAt:   u.CreatedAt,
	}
}

func NewUser(svc userService) (*userHandler) {
	return &userHandler{
		svc: svc,
	}
}

func (h *userHandler) ByName(c *echo.Context) error {
	name := c.Param("name")

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	user, err := h.svc.ByName(ctx, name)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error").Wrap(err)
	}
	return c.JSON(http.StatusOK, userResponse(user))
}

func (h *userHandler) Me(c *echo.Context) error {
	userID := c.Get("user").(string)

	ctx := c.Request().Context()

	user, err := h.svc.ByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error").Wrap(err)
	}
	return c.JSON(http.StatusOK, userResponse(user))
}
