package app

import "github.com/labstack/echo/v4"

type httpController interface {
	Register(e *echo.Group)
}
