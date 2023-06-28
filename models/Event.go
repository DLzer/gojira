package models

const (
	IssueUpdated string = "IssueUpdated"
	IssueCreated string = "IssueCreated"
)

type EventMap struct {
	EventType string `json:"eventType"`
	Created   bool   `json:"created"`
	Updated   bool   `json:"updated"`
}
