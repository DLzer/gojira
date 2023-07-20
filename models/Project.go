package models

// ProjectMap represents a structure that holds mappings from JIRA to it's distribution sources
type ProjectMap struct {
	ProjectKey           string `json:"ProjectKey,omitempty"`
	GitHubRepositoryName string `json:"GitHubRepositoryName,omitempty"`
	DiscordChannelID     string `json:"DiscordChannelID,omitempty"`
}
