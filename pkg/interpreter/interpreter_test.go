package interpreter

import (
	"context"
	"testing"

	"github.com/DLzer/gojira/models"
	"github.com/stretchr/testify/assert"
)

func TestInterpret(t *testing.T) {
	webhookUpdateEvent := &models.JiraWebhookMessage{
		WebhookEvent: "jira:issue_updated",
	}

	webhookCreatedEvent := &models.JiraWebhookMessage{
		WebhookEvent: "jira:issue_created",
	}

	interpretUpdate := Interpret(context.Background(), webhookUpdateEvent)
	assert.Equal(t, "IssueUpdated", interpretUpdate.EventType)

	interpretCreate := Interpret(context.Background(), webhookCreatedEvent)
	assert.Equal(t, "IssueCreated", interpretCreate.EventType)
}

func TestBadInterpret(t *testing.T) {
	webhookMiscEvent := &models.JiraWebhookMessage{
		WebhookEvent: "jira:issue",
	}

	interpretEvent := Interpret(context.Background(), webhookMiscEvent)
	assert.Empty(t, interpretEvent)
}
