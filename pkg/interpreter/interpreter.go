package interpreter

import (
	"context"
	"strings"

	"github.com/DLzer/gojira/models"
	"go.opentelemetry.io/otel"
)

// Interpret reads the incoming message and returns a structure with the event context
func Interpret(ctx context.Context, message *models.JiraWebhookMessage) *models.EventMap {
	_, span := otel.Tracer("Receiver").Start(ctx, "interpreter.Interpret")
	defer span.End()

	var eventMap models.EventMap

	issueData := strings.Split(message.Issue.Key, "-")

	switch message.WebhookEvent {
	case "jira:issue_updated":
		eventMap.EventType = models.IssueUpdated
		eventMap.EventKey = issueData[0]
		eventMap.EventID = issueData[1]
		eventMap.Created = false
		eventMap.Updated = true
	case "jira:issue_created":
		eventMap.EventType = models.IssueCreated
		eventMap.EventKey = issueData[0]
		eventMap.EventID = issueData[1]
		eventMap.Created = false
		eventMap.Updated = true
	default:
		return nil
	}

	return &eventMap
}
