package models

import "time"

type JiraWebhookMessage struct {
	ID           int           `json:"id"`
	Timestamp    *time.Time    `json:"timestamp"`
	Issue        JiraIssue     `json:"issue"`
	User         JiraUser      `json:"user"`
	ChangeLog    JiraChangeLog `json:"changelog"`
	Comments     JiraComment   `json:"comments"`
	WebhookEvent string        `json:"webhookEvent"`
}

type JiraIssue struct {
	ID      int        `json:"id"`
	SelfURL string     `json:"self"`
	Key     string     `json:"key"`
	Fields  JiraFields `json:"fields"`
}

type JiraFields struct {
	Summary     string   `json:"summary"`
	Created     string   `json:"created"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
	Priority    string   `json:"priority"`
}

type JiraUser struct {
	SelfURL      string `json:"self"`
	Name         string `json:"name"`
	Key          string `json:"key"`
	EmailAddress string `json:"emailAddress"`
	DisplayName  string `json:"displayName"`
	Active       string `json:"active"`
}

type JiraChangeLog struct {
	ID    int                  `json:"id"`
	Items []JiraChangeLogItems `json:"items"`
}

type JiraChangeLogItems struct {
	ToString   string `json:"toString"`
	To         string `json:"to"`
	FromString string `json:"fromString"`
	From       string `json:"from"`
	FieldType  string `json:"fieldType"`
	Field      string `json:"field"`
}

type JiraComment struct {
	SelfURL string `json:"self"`
	ID      string `json:"id"`
	Body    string `json:"body"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}
