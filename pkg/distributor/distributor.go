package distributor

import (
	"context"
	"fmt"
	"log"

	"github.com/DLzer/gojira/config"
	"github.com/DLzer/gojira/models"
	"github.com/bwmarrin/discordgo"
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

// Distribute will broadcast multiple message distributions to multiple sources if those sources are enabled
func Distribute(ctx context.Context, cfg *config.Config, message *models.JiraWebhookMessage, projectMap *models.ProjectMap, dg *discordgo.Session) error {
	_, span := otel.Tracer("Receiver").Start(ctx, "distributor.Distribute")
	defer span.End()

	// Distribute to Github
	if cfg.Github.Enable {
		_ = DistributeToGithub(ctx, cfg)
	}

	// Distribute to Discord
	if cfg.Discord.Enable {
		_ = DistributeToDiscord(ctx, cfg, message, projectMap, dg)
	}

	return nil
}

// DistributeToGithub will craft and send an issue genration request
func DistributeToGithub(ctx context.Context, cfg *config.Config) error {
	_, span := otel.Tracer("Receiver").Start(ctx, "distributor.Distribute.DistributeToGithub")
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
		return err
	}

	return nil
}

// DistributeToDiscord will craft a message embed and send it to discord
func DistributeToDiscord(ctx context.Context, cfg *config.Config, message *models.JiraWebhookMessage, projectMap *models.ProjectMap, dg *discordgo.Session) error {
	_, span := otel.Tracer("Receiver").Start(ctx, "distributor.Distribute.DistributeToDiscord")
	defer span.End()

	embed := &discordgo.MessageEmbed{
		URL:         fmt.Sprintf("%s%s", cfg.Jira.BaseUrl, projectMap.JiraProjectKey),
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x00ff00,
		Description: message.Issue.Fields.Summary,
	}

	_, err := dg.ChannelMessageSendEmbed(projectMap.DiscordProjectID, embed)
	if err != nil {
		log.Print("Embed Send Error", err)
		return err
	}

	return nil
}
