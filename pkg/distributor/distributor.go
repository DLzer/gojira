package distributor

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/DLzer/gojira/config"
	"github.com/DLzer/gojira/models"
	"github.com/DLzer/gojira/pkg/utils"
	"github.com/bwmarrin/discordgo"
	"go.opentelemetry.io/otel"
)

// MapDistribution readsd the incomning message and event to return an object with the relative outgoing project IDs
func MapDistribution(ctx context.Context, message *models.JiraWebhookMessage, event *models.EventMap) (*models.ProjectMap, error) {
	_, span := otel.Tracer("Receiver").Start(ctx, "distributor.MapDistribution")
	defer span.End()

	// Load ProjectMap File
	file, err := ioutil.ReadFile("project_map.json")
	if err != nil {
		log.Fatal("Error", err)
	}

	var projectMap []models.ProjectMap

	err = json.Unmarshal([]byte(file), &projectMap)
	if err != nil {
		log.Fatal("Error", err)
	}

	for x := range projectMap {
		if utils.StringSliceContains(projectMap[x].ProjectKey, event.EventKey) {
			return &projectMap[x], nil
		}
	}

	return nil, nil
}

// Distribute will broadcast multiple message distributions to multiple sources if those sources are enabled
func Distribute(ctx context.Context, cfg *config.Config, message *models.JiraWebhookMessage, projectMap *models.ProjectMap, eventMap *models.EventMap, dg *discordgo.Session) error {
	_, span := otel.Tracer("Receiver").Start(ctx, "distributor.Distribute")
	defer span.End()

	// Distribute to Github
	if cfg.Github.Enable {
		_ = DistributeToGithub(ctx, cfg)
	}

	// Distribute to Discord
	if cfg.Discord.Enable {
		err := DistributeToDiscord(ctx, cfg, message, projectMap, eventMap, dg)
		if err != nil {
			return err
		}
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
func DistributeToDiscord(ctx context.Context, cfg *config.Config, message *models.JiraWebhookMessage, projectMap *models.ProjectMap, event *models.EventMap, dg *discordgo.Session) error {
	_, span := otel.Tracer("Receiver").Start(ctx, "distributor.Distribute.DistributeToDiscord")
	defer span.End()

	fmt.Println("Attempting Discord Dispatch to Channel: ", projectMap.DiscordChannelID)

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s-%s %s", event.EventKey, event.EventID, message.Issue.Fields.Summary),
		URL:   fmt.Sprintf("%s/%s%s", cfg.Jira.BaseUrl, event.EventKey, event.EventID),
		Author: &discordgo.MessageEmbedAuthor{
			Name:    strings.Title(message.User.Name),
			IconURL: message.User.AvatarUrls.Small,
		},
		Color:       0x00ff00,
		Description: message.Issue.Fields.Description,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Powered by your friend, GoJIRA",
		},
	}

	_, err := dg.ChannelMessageSendEmbed(projectMap.DiscordChannelID, embed)
	if err != nil {
		log.Println("Embed Send Error", err)
		return err
	}

	return nil
}
