package interpreter

import (
	"context"

	"github.com/DLzer/gojira/models"
	"go.opentelemetry.io/otel"
)

// Interpret reads the incoming message and returns a structure with the event context
func Interpret(ctx context.Context, message *models.JiraWebhookMessage) *models.EventMap {
	_, span := otel.Tracer("Receiver").Start(ctx, "interpreter.Interpret")
	defer span.End()

	var eventMap models.EventMap

	switch message.WebhookEvent {
	case "jira:issue_updated":
		eventMap.EventType = models.IssueUpdated
		eventMap.Created = false
		eventMap.Updated = true
	case "jira:issue_created":
		eventMap.EventType = models.IssueCreated
		eventMap.Created = false
		eventMap.Updated = true
	default:
		return nil
	}

	return &eventMap
}
