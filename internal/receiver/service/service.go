package service

import (
	"context"
	"fmt"

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

// Accept receives the JIRA Webhook Message and starts the process of interpreting and dispatching to all enabled services
func (s *receiverService) Accept(ctx context.Context, message *models.JiraWebhookMessage, dg *discordgo.Session) error {
	ctx, span := otel.Tracer("Receiver").Start(ctx, "receiverService.Accept")
	defer span.End()

	// Interpret
	eventMap := interpreter.Interpret(ctx, message)

	fmt.Println("Event Map", eventMap)

	// Distribution Map
	distributionMap, err := distributor.MapDistribution(ctx, message, eventMap)
	if err != nil {
		s.logger.Info("Error", err)
		return nil
	}

	fmt.Println("Distribution Map", distributionMap)

	// Distribute
	err = distributor.Distribute(ctx, s.cfg, message, distributionMap, eventMap, dg)
	if err != nil {
		s.logger.Info(err)
		return nil
	}

	return nil
}
