package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cscercel/doty-lore-bot/internal/handler"
	"github.com/cscercel/doty-lore-bot/internal/db"
)

func RegisterHandlers(s *discordgo.Session, q *db.Queries) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			data := i.ApplicationCommandData()
			if data.Name != "lore" || len(data.Options) == 0 {
				return
			}
			sub := data.Options[0]
			switch sub.Name {
			case "create":
				handler.HandleCreate(s, i, sub.Options, q)
			case "view":
				handler.HandleView(s, i, sub.Options, q)
			case "list":
				handler.HandleList(s, i, sub.Options, q)
			case "edit":
				handler.HandleEdit(s, i, sub.Options, q)
			case "delete":
				handler.HandleDelete(s, i, sub.Options, q)
			}

		case discordgo.InteractionApplicationCommandAutocomplete:
			data := i.ApplicationCommandData()
			if data.Name != "lore" || len(data.Options) == 0 {
				return
			}
			sub := data.Options[0]
			if sub.Name == "view" || sub.Name == "edit" || sub.Name == "delete" {
				handler.HandleViewAutocomplete(s, i, sub.Options, q)
			}
		}
	})
}
