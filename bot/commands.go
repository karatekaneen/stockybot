package bot

import "github.com/bwmarrin/discordgo"

// var integerOptionMinValue = 1.0

func (b *DiscordBot) listCommands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "pending",
			Description: "List pending signals",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "all-lists",
					Description: "Show signals from all lists. Only listing Swedish Large, Mid and Small cap if false",
				},
			},
		},
	}
}

// var allCommands = []*discordgo.ApplicationCommand{
// {
// 	Name: "basic-command",
// 	// All commands and options must have a description
// 	// Commands/options without description will fail the registration
// 	// of the command.
// 	Description: "Basic command",
// },
// {
// 	Name:        "basic-command-with-files",
// 	Description: "Basic command with files",
// },
// {
// 	Name:        "options",
// 	Description: "Command for demonstrating options",
// 	Options: []*discordgo.ApplicationCommandOption{

// 		{
// 			Type:        discordgo.ApplicationCommandOptionString,
// 			Name:        "string-option",
// 			Description: "String option",
// 			Required:    true,
// 		},
// 		{
// 			Type:        discordgo.ApplicationCommandOptionInteger,
// 			Name:        "integer-option",
// 			Description: "Integer option",
// 			MinValue:    &integerOptionMinValue,
// 			MaxValue:    10,
// 			Required:    true,
// 		},
// 		{
// 			Type:        discordgo.ApplicationCommandOptionNumber,
// 			Name:        "number-option",
// 			Description: "Float option",
// 			MaxValue:    10.1,
// 			Required:    true,
// 		},
// 		{
// 			Type:        discordgo.ApplicationCommandOptionBoolean,
// 			Name:        "bool-option",
// 			Description: "Boolean option",
// 			Required:    true,

// 		},

// 		// Required options must be listed first since optional parameters
// 		// always come after when they're used.
// 		// The same concept applies to Discord's Slash-commands API

// 		{
// 			Type:        discordgo.ApplicationCommandOptionChannel,
// 			Name:        "channel-option",
// 			Description: "Channel option",
// 			// Channel type mask
// 			ChannelTypes: []discordgo.ChannelType{
// 				discordgo.ChannelTypeGuildText,
// 				discordgo.ChannelTypeGuildVoice,
// 			},
// 			Required: false,
// 		},
// 		{
// 			Type:        discordgo.ApplicationCommandOptionUser,
// 			Name:        "user-option",
// 			Description: "User option",
// 			Required:    false,
// 		},
// 		{
// 			Type:        discordgo.ApplicationCommandOptionRole,
// 			Name:        "role-option",
// 			Description: "Role option",
// 			Required:    false,
// 		},
// 	},
// },

// {
// 	Name:        "responses",
// 	Description: "Interaction responses testing initiative",
// 	Options: []*discordgo.ApplicationCommandOption{
// 		{
// 			Name:        "resp-type",
// 			Description: "Response type",
// 			Type:        discordgo.ApplicationCommandOptionInteger,
// 			Choices: []*discordgo.ApplicationCommandOptionChoice{
// 				{
// 					Name:  "Channel message with source",
// 					Value: 4,
// 				},
// 				{
// 					Name:  "Deferred response With Source",
// 					Value: 5,
// 				},
// 			},
// 			Required: true,
// 		},
// 	},
// },
// {
// 	Name:        "followups",
// 	Description: "Followup messages",
// },
// }
