package receiver

import (
	"context"

	"github.com/DLzer/gojira/models"
	"github.com/bwmarrin/discordgo"
)

type Service interface {
	Accept(ctx context.Context, message *models.JiraWebhookMessage, dg *discordgo.Session) error
}
