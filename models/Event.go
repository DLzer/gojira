package models

const (
	IssueUpdated string = "IssueUpdated"
	IssueCreated string = "IssueCreated"
)

type EventMap struct {
	EventType string `json:"eventType"`
	EventKey  string `json:"eventKey"`
	EventID   string `json:"eventID"`
	Created   bool   `json:"created"`
	Updated   bool   `json:"updated"`
}
