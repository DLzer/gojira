package models

type DiscordMessageComponentType = int

// DiscordMessageComponent represents a Discord Message model.
type DiscordMessageComponent struct {
	Type        DiscordMessageComponentType `json:"type"`
	Style       int                         `json:"style"`
	Label       string                      `json:"label"`
	CustomID    string                      `json:"custom_id"`
	Url         string                      `json:"url,omitempty"`
	Disabled    bool                        `json:"disabled"`
	Placeholder string                      `json:"placeholder"`
	MinValues   int                         `json:"min_values"`
	MaxValues   int                         `json:"max_values"`
	Required    bool                        `json:"required"`
	Value       string                      `json:"value,omitempty"`
}
