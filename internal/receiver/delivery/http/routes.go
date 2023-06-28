package http

import (
	"github.com/DLzer/gojira/internal/receiver"
	"github.com/labstack/echo/v4"
)

// Map receiver routes
func MapReceiverRoutes(receiverGroup *echo.Group, h receiver.Handlers) {
	receiverGroup.POST("/accept", h.Accept())
}
