package models

type ProjectMap struct {
	JiraProjectKey   string `json:"jiraProjectKey"`
	DiscordProjectID string `json:"discordProjectID"`
	GithubProjectKey string `json:"gituhbProjectKey"`
}
