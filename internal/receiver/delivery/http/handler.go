package http

import (
	"net/http"

	"github.com/DLzer/gojira/config"
	"github.com/DLzer/gojira/internal/receiver"
	"github.com/DLzer/gojira/models"
	"github.com/DLzer/gojira/pkg/httpErrors"
	"github.com/DLzer/gojira/pkg/logger"
	"github.com/DLzer/gojira/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
)

type receiverHandlers struct {
	cfg             *config.Config
	receiverService receiver.Service
	logger          logger.Logger
}

func (h receiverHandlers) Accept() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := otel.Tracer("Receiver").Start(utils.GetRequestCtx(c), "receiverHandlers.Accept")

		p := &models.JiraWebhookMessage{}
		if err := c.Bind(p); err != nil {
			utils.LogResponseError(c, h.logger, err)
			span.End()
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		if err := h.receiverService.Accept(ctx, p); err != nil {
			utils.LogResponseError(c, h.logger, err)
			span.End()
			return c.JSON(httpErrors.ErrorResponse(err))
		}

		span.End()
		return c.JSON(http.StatusOK, map[string]interface{}{"message": "ok"})
	}
}
