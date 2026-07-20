package discord

import "github.com/bwmarrin/discordgo"

var commandDefs = []*discordgo.ApplicationCommand{
	{
		Name:        "lore",
		Description: "Manage D&D lore cards",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "create",
				Description: "Create a new lore card",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type: discordgo.ApplicationCommandOptionString,
						Name: "name",
						Description: "Card name",
						Required: true,
					},
					{
						Type: discordgo.ApplicationCommandOptionString,
						Name: "type",
						Description: "Card type",
						Required: true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "Location", Value: "location"},
							{Name: "NPC", Value: "npc"},
							{Name: "Rule", Value: "rule"},
							{Name: "Religion", Value: "religion"},
							{Name: "Event", Value: "event"},
							{Name: "Organization", Value: "organization"},
						},
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "view",
				Description: "View a lore card",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type: discordgo.ApplicationCommandOptionString,
						Name: "name",
						Description: "Card name",
						Required: true,
						Autocomplete: true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "list",
				Description: "List lore cards",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type: discordgo.ApplicationCommandOptionString,
						Name: "type",
						Description: "Filter by card type",
						Required: false,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "Location", Value: "location"},
							{Name: "NPC", Value: "npc"},
							{Name: "Rule", Value: "rule"},
							{Name: "Religion", Value: "religion"},
							{Name: "Event", Value: "event"},
							{Name: "Organization", Value: "organization"},
						},
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "edit",
				Description: "Edit an existing lore card",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type: discordgo.ApplicationCommandOptionString,
						Name: "name",
						Description: "Card to edit",
						Required: true,
						Autocomplete: true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "delete",
				Description: "Delete a lore card",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type: discordgo.ApplicationCommandOptionString,
						Name: "name",
						Description: "Card to delete",
						Required: true,
						Autocomplete: true,
					},
				},
			},
		},
	},
}

func RegisterCommands(s *discordgo.Session, guildID string) ([]*discordgo.ApplicationCommand, error) {
	return s.ApplicationCommandBulkOverwrite(s.State.User.ID, guildID, commandDefs)
}

func CleanupCommands(s *discordgo.Session, guildID string, cmds []*discordgo.ApplicationCommand) {
	//TODO
}
