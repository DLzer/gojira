package discord_cmd

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type DiscordCommandsHandler struct {
	Session            *discordgo.Session
	GuildID            string
	RegisteredCommands []*discordgo.ApplicationCommand
}

func NewDiscordCommandsHandler(s *discordgo.Session, gid string) *DiscordCommandsHandler {
	return &DiscordCommandsHandler{Session: s, GuildID: gid}
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = true
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	commands = []*discordgo.ApplicationCommand{
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

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
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

			msgformat := fmt.Sprintf("Your Prompt: %s\n", prompt)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, they will be discussed in "responses"
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msgformat,
				},
			})
		},
	}
)

func (d *DiscordCommandsHandler) Init() {
	d.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func (d *DiscordCommandsHandler) EnableCommands() {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := d.Session.ApplicationCommandCreate(d.Session.State.User.ID, d.GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	d.RegisteredCommands = registeredCommands
}

func (d *DiscordCommandsHandler) RemoveCommands() {
	log.Println("Removing commands...")
	// // We need to fetch the commands, since deleting requires the command ID.
	// // We are doing this from the returned commands on line 375, because using
	// // this will delete all the commands, which might not be desirable, so we
	// // are deleting only the commands that we added.
	// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
	// if err != nil {
	// 	log.Fatalf("Could not fetch registered commands: %v", err)
	// }

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
