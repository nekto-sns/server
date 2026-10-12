package route

import (
	"github.com/labstack/echo/v5"
)

func Setup(e *echo.Echo, uh userHandler, sh scratchAuthHandler, am echo.MiddlewareFunc) {
	userRoute(e, uh, am)
	saRoute(e, sh)
}
