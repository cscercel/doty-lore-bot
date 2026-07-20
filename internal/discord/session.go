package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/cscercel/doty-lore-bot/internal/db"
	"github.com/cscercel/doty-lore-bot/internal/handler"
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
				handler.HandleCreateStart(s, i, sub.Options, q)
			case "view":
				handler.HandleView(s, i, sub.Options, q)
			case "list":
				handler.HandleList(s, i, sub.Options, q)
			case "edit":
				handler.HandleEditStart(s, i, sub.Options, q)
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

		case discordgo.InteractionModalSubmit:
			customID := i.ModalSubmitData().CustomID
			switch {
			case strings.HasPrefix(customID, "lore_create_modal:"):
				handler.HandleCreateSubmit(s, i, q)
			case strings.HasPrefix(customID, "lore_edit_modal:"):
				handler.HandleEditSubmit(s, i, q)
			}
		}
	})
}
