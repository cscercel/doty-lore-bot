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
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "name",
						Description: "Card name",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "type",
						Description: "Card type",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "Location", Value: "location"},
							{Name: "NPC", Value: "npc"},
							{Name: "Rule", Value: "rule"},
							{Name: "Religion", Value: "religion"},
							{Name: "Event", Value: "event"},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "summary",
						Description: "Short description",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "body",
						Description: "Full lore text (optional, can add later)",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionAttachment,
						Name:        "image",
						Description: "Portrait or map image",
						Required:    false,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "view",
				Description: "View a lore card",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "name",
						Description:  "Card name",
						Required:     true,
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
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "type",
						Description: "Filter by card type",
						Required:    false,
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
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "name",
						Description:  "Card to edit",
						Required:     true,
						Autocomplete: true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "summary",
						Description: "New summary (leave blank to keep current)",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "body",
						Description: "New full lore text (leave blank to keep current)",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "tags",
						Description: "Comma-separated tags (leave blank to keep current)",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionAttachment,
						Name:        "image",
						Description: "New image (leave blank to keep current)",
						Required:    false,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "delete",
				Description: "Delete a lore card",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionString,
						Name:         "name",
						Description:  "Card to delete",
						Required:     true,
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
	// TODO
}
