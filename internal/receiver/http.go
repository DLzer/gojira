package receiver

import "github.com/labstack/echo/v4"

type Handlers interface {
	Accept() echo.HandlerFunc
}
