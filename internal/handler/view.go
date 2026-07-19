package handler

import (
	"context"
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"

	"github.com/cscercel/doty-lore-bot/internal/db"
)

func HandleViewAutocomplete(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	opts []*discordgo.ApplicationCommandInteractionDataOption, q *db.Queries,
) {
	var input string
	for _, o := range opts {
		if o.Name == "name" {
			input = o.StringValue()
		}
	}

	searchTerm := input + "%"
	cards, err := q.SearchCardsByName(context.Background(), &searchTerm)
	if err != nil {
		log.Printf("autocomplete search error: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{Choices: []*discordgo.ApplicationCommandOptionChoice{}},
		})
		return
	}

	choices := make([]*discordgo.ApplicationCommandOptionChoice, 0, len(cards))
	for _, c := range cards {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  c.Name,
			Value: c.Name,
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{Choices: choices},
	})
}

func HandleView(
	s *discordgo.Session, 
	i *discordgo.InteractionCreate, 
	opts []*discordgo.ApplicationCommandInteractionDataOption, q *db.Queries,
) {
	var name string
	for _, o := range opts {
		if o.Name == "name" {
			name = o.StringValue()
		}
	}

	card, err := q.GetCardByName(context.Background(), name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respond(s, i, "No card found named \""+name+"\".")
			return
		}
		log.Printf("view card error: %v", err)
		respond(s, i, "Something went wrong looking that up.")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       card.Name,
		Description: card.Summary,
		Color:       0x8B5CF6,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Type", Value: card.Type, Inline: true},
		},
	}
	if card.ImageUrl != nil {
		embed.Image = &discordgo.MessageEmbedImage{URL: *card.ImageUrl}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})

	log.Printf("card viewed: name=%q by=%s", card.Name, i.Member.User.Username)
}
