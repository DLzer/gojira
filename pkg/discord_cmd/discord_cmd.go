package discord_cmd

import (
	"fmt"
	"log"

	"github.com/DLzer/gojira/pkg/chatgpt"
	"github.com/bwmarrin/discordgo"
)

// Holds information about Dicsord Commands
type DiscordCommandsHandler struct {
	Session            *discordgo.Session
	GuildID            string
	ChatGPTToken       string
	RegisteredCommands []*discordgo.ApplicationCommand
	Commands           []*discordgo.ApplicationCommand
	CommandHandlers    map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

// A constructor to build a new Discord Command Handler
func NewDiscordCommandsHandler(s *discordgo.Session, gid string, token string) *DiscordCommandsHandler {
	return &DiscordCommandsHandler{Session: s, GuildID: gid, ChatGPTToken: token}
}

// Initializes commands and their handlers
func (d *DiscordCommandsHandler) Init() {
	d.Commands = d.GetCommands()
	d.CommandHandlers = d.GetCommandHandlers()

	d.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := d.CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

// Enables the range of commands available in the global struct
func (d *DiscordCommandsHandler) EnableCommands() {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(d.Commands))
	for i, v := range d.Commands {
		cmd, err := d.Session.ApplicationCommandCreate(d.Session.State.User.ID, d.GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	d.RegisteredCommands = registeredCommands
}

// Removes the range of commands available in the global struct
func (d *DiscordCommandsHandler) RemoveCommands() {
	log.Println("Removing commands...")
	// // We need to fetch the commands, since deleting requires the command ID.
	// // We are doing this from the returned commands on line 375, because using
	// // this will delete all the commands, which might not be desirable, so we
	// // are deleting only the commands that we added.

	if len(d.RegisteredCommands) == 0 {
		log.Println("No commands to remove...")
		return
	}

	for _, v := range d.RegisteredCommands {
		err := d.Session.ApplicationCommandDelete(d.Session.State.User.ID, d.GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

// Wrapper for making a GPT request
func GPTRequest(token string, message string) string {
	res, err := chatgpt.MakeGPTRequest(token, message)
	if err != nil {
		return "An error occurred in ChatGPT, please excuse my ignorance."
	}

	return res
}

// Returns a slice of the top level command options
func (d *DiscordCommandsHandler) GetCommands() []*discordgo.ApplicationCommand {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "go-gpt",
			Description: "Ask a question to GoJIRA-gpt",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "prompt",
					Description: "Write a text prompt",
					Required:    true,
				},
			},
		},
	}
	return commands
}

// Returns a slice of the command callback methods
func (d *DiscordCommandsHandler) GetCommandHandlers() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"go-gpt": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			var prompt string
			if option, ok := optionMap["prompt"]; ok {
				prompt = option.StringValue()
			}

			msgformat := fmt.Sprintf("**Your Prompt:** %s\n *Your answer should appear shortly*!\n", prompt)

			if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msgformat,
				},
			}); err != nil {
				log.Println("Error:", err)
			}

			response := GPTRequest(d.ChatGPTToken, prompt+" in less than 1500 chracters.")

			content := msgformat + "**Reply:** " + response
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})

			if err != nil {
				log.Println("Error:", err)
				if _, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				}); err != nil {
					log.Println("Error:", err)
				}
				return
			}
		},
	}
	return commandHandlers
}
