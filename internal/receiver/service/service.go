package service

import (
	"context"

	"github.com/DLzer/gojira/config"
	"github.com/DLzer/gojira/internal/receiver"
	"github.com/DLzer/gojira/models"
	"github.com/DLzer/gojira/pkg/interpreter"
	"github.com/DLzer/gojira/pkg/logger"
	"go.opentelemetry.io/otel"
)

type receiverService struct {
	cfg    *config.Config
	logger logger.Logger
}

func NewReceiverService(cfg *config.Config, logger logger.Logger) receiver.Service {
	return &receiverService{cfg: cfg, logger: logger}
}

func (s *receiverService) Accept(ctx context.Context, message *models.JiraWebhookMessage) error {
	_, span := otel.Tracer("Receiver").Start(ctx, "receiverService.Accept")
	defer span.End()

	_ = interpreter.Interpret(ctx, message)

	return nil
}
