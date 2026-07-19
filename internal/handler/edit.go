package handler

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"

	"github.com/cscercel/doty-lore-bot/internal/db"
)

func HandleEdit(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	opts []*discordgo.ApplicationCommandInteractionDataOption,
	q *db.Queries,
) {
	optMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, o := range opts {
		optMap[o.Name] = o
	}

	name := optMap["name"].StringValue()

	existing, err := q.GetCardByName(context.Background(), name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			respond(s, i, "No card found named \""+name+"\".")
			return
		}
		log.Printf("edit lookup error: %v", err)
		respond(s, i, "Something went wrong looking that up.")
		return
	}

	body := existing.Body
	if o, ok := optMap["body"]; ok {
		v := o.StringValue()
		body = &v
	}

	tags := existing.Tags
	if o, ok := optMap["tags"]; ok {
		raw := strings.Split(o.StringValue(), ",")
		tags = make([]string, 0, len(raw))
		for _, t := range raw {
			t = strings.TrimSpace(t)
			if t != "" {
				tags = append(tags, t)
			}
		}
	}

	summary := existing.Summary
	if s, ok := optMap["summary"]; ok {
		summary = s.StringValue()
	}

	imageURL := existing.ImageUrl
	if imgOpt, ok := optMap["image"]; ok {
		attachmentID := imgOpt.Value.(string)
		if att, ok := i.ApplicationCommandData().Resolved.Attachments[attachmentID]; ok {
			url := att.URL
			imageURL = &url
		}
	}

	updated, err := q.UpdateCard(context.Background(), db.UpdateCardParams{
		ID:       existing.ID,
		Summary:  summary,
		Body:     body,
		Tags:	  tags,
		ImageUrl: imageURL,
	})
	if err != nil {
		log.Printf("update card error: %v", err)
		respond(s, i, "Something went wrong updating that card.")
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       updated.Name,
		Description: updated.Summary,
		Color:       0x8B5CF6,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Type", Value: updated.Type, Inline: true},
		},
	}
	if len(updated.Tags) > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Tags", Value: strings.Join(updated.Tags, ", "), Inline: true,
		})
	}
	if updated.ImageUrl != nil {
		embed.Image = &discordgo.MessageEmbedImage{URL: *updated.ImageUrl}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: []*discordgo.MessageEmbed{embed}},
	})

	log.Printf("card edited: name=%q by=%s", updated.Name, i.Member.User.Username)
}
