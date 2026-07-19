package handler

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/cscercel/doty-lore-bot/internal/db"
)

func HandleList(
	s *discordgo.Session, 
	i *discordgo.InteractionCreate, 
	opts []*discordgo.ApplicationCommandInteractionDataOption, q *db.Queries,
) {
	var typeFilter string
	for _, o := range opts {
		if o.Name == "type" {
			typeFilter = o.StringValue()
		}
	}

	var cards []db.LoreCard
	var err error
	title := "All Lore Cards"

	if typeFilter != "" {
		cards, err = q.ListCardsByType(context.Background(), typeFilter)
		title = fmt.Sprintf("%s Cards", typeFilter)
	} else {
		cards, err = q.ListCards(context.Background())
	}

	if err != nil {
		log.Printf("list cards error: %v", err)
		respond(s, i, "Something went wrong listing cards.")
		return
	}

	if len(cards) == 0 {
		respond(s, i, "No cards found.")
		return
	}

	var sb strings.Builder
	for _, c := range cards {
		fmt.Fprintf(&sb, "**%s** — _%s_\n", c.Name, c.Type)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       title,
					Description: sb.String(),
					Color:       0x8B5CF6,
				},
			},
		},
	})

	log.Printf("cards listed: type=%q count=%d by=%s", typeFilter, len(cards), i.Member.User.Username)
}
