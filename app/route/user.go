package route

import (
	"github.com/labstack/echo/v5"
)

type userHandler interface{
	ByName(*echo.Context) error
	Me(*echo.Context) error
}

func userRoute(e *echo.Echo, uh userHandler, am echo.MiddlewareFunc) {
	u := e.Group("/users")
	u.GET("/me", uh.Me, am)
	u.GET("/:name", uh.ByName)
}
