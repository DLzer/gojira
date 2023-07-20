package service

import (
	"context"

	"github.com/DLzer/gojira/config"
	"github.com/DLzer/gojira/internal/receiver"
	"github.com/DLzer/gojira/models"
	"github.com/DLzer/gojira/pkg/distributor"
	"github.com/DLzer/gojira/pkg/interpreter"
	"github.com/DLzer/gojira/pkg/logger"
	"github.com/bwmarrin/discordgo"
	"go.opentelemetry.io/otel"
)

type receiverService struct {
	cfg    *config.Config
	logger logger.Logger
}

func NewReceiverService(cfg *config.Config, logger logger.Logger) receiver.Service {
	return &receiverService{cfg: cfg, logger: logger}
}

func (s *receiverService) Accept(ctx context.Context, message *models.JiraWebhookMessage, dg *discordgo.Session) error {
	ctx, span := otel.Tracer("Receiver").Start(ctx, "receiverService.Accept")
	defer span.End()

	eventMap := interpreter.Interpret(ctx, message)
	distributionMap := distributor.MapDistribution(ctx, message, eventMap)

	err := distributor.Distribute(ctx, s.cfg, message, distributionMap, dg)
	if err != nil {
		s.logger.Info(err)
	}

	return nil
}
