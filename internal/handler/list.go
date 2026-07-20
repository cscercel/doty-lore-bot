package handler

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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
		title = fmt.Sprintf("%s Cards", cases.Title(language.AmericanEnglish).String(typeFilter))
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

	embed := &discordgo.MessageEmbed{
		Title:       title,
		Description: sb.String(),
		Color:       0x8B5CF6,
	}

	dmChannel, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		log.Printf("dm channel create error: %v", err)
		respond(s, i, "Couldn't open a DM — check your privacy settings allow DMs from server members.")
		return
	}
	if _, err := s.ChannelMessageSendEmbed(dmChannel.ID, embed); err != nil {
		log.Printf("dm send error: %v", err)
		respond(s, i, "Couldn't send you a DM — check your privacy settings allow DMs from server members.")
		return
	}

	respond(s, i, "Sent the card list to your DMs.")
	log.Printf("cards listed via dm: type=%q count=%d by=%s", typeFilter, len(cards), i.Member.User.Username)
}
