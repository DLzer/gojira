package models

type ChatGPTResponse struct {
	ID        string           `json:"id,omitempty"`
	Object    string           `json:"object,omitempty"`
	CreatedAt int              `json:"created_at,omitempty"`
	Choices   []ChatGPTChoices `json:"choices,omitempty"`
}

type ChatGPTChoices struct {
	Index        int            `json:"index,omitempty"`
	Message      ChatGPTMessage `json:"message,omitempty"`
	FinishReason string         `json:"finish_reason,omitempty"`
}

type ChatGPTMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}
