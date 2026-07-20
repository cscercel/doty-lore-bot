package handler

import (
	"context"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/cscercel/doty-lore-bot/internal/db"
)

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func awaitImageReply(s *discordgo.Session, channelID, userID string, cardID int32, q *db.Queries) {
	done := make(chan *discordgo.MessageCreate, 1)

	removeHandler := s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.ChannelID != channelID || m.Author.ID != userID || len(m.Attachments) == 0 {
			return
		}
		select {
		case done <- m:
		default:
		}
	})

	go func() {
		defer removeHandler()
		select {
		case m := <-done:
			url := m.Attachments[0].URL
			_, err := q.UpdateCardImage(context.Background(), db.UpdateCardImageParams{
				ID:       cardID,
				ImageUrl: &url,
			})
			if err != nil {
				log.Printf("update card image error: %v", err)
				s.ChannelMessageSend(channelID, "Couldn't save that image — try `/lore edit` instead.")
				return
			}
			s.ChannelMessageSend(channelID, "🖼️ Image saved.")
			log.Printf("card image updated: id=%d by=%s", cardID, userID)

		case <-time.After(60 * time.Second):
			s.ChannelMessageSend(channelID, "No image received — skipped.")
		}
	}()
}
