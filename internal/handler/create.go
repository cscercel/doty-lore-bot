package handler

import (
	"context"
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/cscercel/doty-lore-bot/internal/db"
)

func HandleCreate(
	s *discordgo.Session, 
	i *discordgo.InteractionCreate, 
	opts []*discordgo.ApplicationCommandInteractionDataOption, 
	q *db.Queries,
) {
	optMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, o := range opts {
		optMap[o.Name] = o
	}

	var body *string
	if o, ok := optMap["body"]; ok {
		v := o.StringValue()
		body = &v
	}

	var imageURL *string
	if imgOpt, ok := optMap["image"]; ok {
		attachmentID := imgOpt.Value.(string)
		if att, ok := i.ApplicationCommandData().Resolved.Attachments[attachmentID]; ok {
			url := att.URL
			imageURL = &url
		}
	}

	params := db.CreateCardParams{
		Name:     optMap["name"].StringValue(),
		Type:     optMap["type"].StringValue(),
		Summary:  optMap["summary"].StringValue(),
		Body:     body,
		ImageUrl: imageURL,
	}

	card, err := q.CreateCard(context.Background(), params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			respond(s, i, "A card named \""+params.Name+"\" already exists.")
			return
		}
		log.Printf("create card error: %v", err)
		respond(s, i, "Something went wrong creating that card.")
		return
	}

	log.Printf("card created: id=%d name=%q type=%q by=%s", card.ID, card.Name, card.Type, i.Member.User.Username)

	embed := &discordgo.MessageEmbed{
		Title:       card.Name,
		Description: card.Summary,
		Color:       0x8B5CF6,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Type", Value: card.Type, Inline: true},
		},
	}
	if imageURL != nil {
		embed.Image = &discordgo.MessageEmbedImage{URL: *imageURL}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
