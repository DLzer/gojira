package distributor

import (
	"context"
	"log"

	"github.com/DLzer/gojira/config"
	"github.com/DLzer/gojira/models"
	"go.opentelemetry.io/otel"
)

// MapDistribution readsd the incomning message and event to return an object with the relative outgoing project IDs
func MapDistribution(ctx context.Context, message *models.JiraWebhookMessage, event *models.EventMap) *models.ProjectMap {
	_, span := otel.Tracer("Receiver").Start(ctx, "distributor.MapDistribution")
	defer span.End()

	// Based on the incoming event type we will perform a few actions.
	// - Determine the github project ID from the JiraProjectKey
	// - Determine the discorcd channel ID from the JiraProjectKey
	// - Put together our struct for response

	return nil
}

func Distribute(ctx context.Context, cfg *config.Config) error {
	_, span := otel.Tracer("Receiver").Start(ctx, "distributor.Distribute")
	defer span.End()

	var repoTitle string = "repo_title"
	var issueTitle string = "testIssueTitle"
	var issueBody string = "testIssueBody"

	// Send data to GitHub project/GitHub issues
	gitHubRequest := &models.GitHubRequest{
		GitHubRepoOwner: &cfg.Github.Owner,
		GitHubRepoTitle: &repoTitle,
		GitHubToken:     &cfg.Github.Token,
	}
	gitHubRequest.Issue = *gitHubRequest.GenerateIssue(issueTitle, &issueBody)

	if err := gitHubRequest.Send(); err != nil {
		log.Fatal("GitHub Send Error:", err)
	}

	// Send data to Discord as a channel message with context

	return nil
}
