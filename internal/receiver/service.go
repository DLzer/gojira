package receiver

import (
	"context"

	"github.com/DLzer/gojira/models"
)

type Service interface {
	Accept(ctx context.Context, message *models.JiraWebhookMessage) error
}
