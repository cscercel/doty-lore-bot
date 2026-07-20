package handler

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/cscercel/doty-lore-bot/internal/db"
)

func HandleCreateStart(
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
	cardType := optMap["type"].StringValue()
	customID := "lore_create_modal:" + cardType + ":" + name

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: customID,
			Title:    "New Lore Card",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "summary",
						Label:       "Summary",
						Style:       discordgo.TextInputShort,
						Required:    true,
						MaxLength:   200,
						Placeholder: "A one-line description",
					},
				}},
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "body",
						Label:       "Full Details",
						Style:       discordgo.TextInputParagraph,
						Required:    false,
						MaxLength:   4000,
						Placeholder: "Backstory, notes, whatever's useful",
					},
				}},
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "tags",
						Label:       "Tags (comma-separated)",
						Style:       discordgo.TextInputShort,
						Required:    false,
						MaxLength:   200,
						Placeholder: "tavern, ally, act1",
					},
				}},
			},
		},
	})
	if err != nil {
		log.Printf("modal open error: %v", err)
	}
}

func HandleCreateSubmit(s *discordgo.Session, i *discordgo.InteractionCreate, q *db.Queries) {
	data := i.ModalSubmitData()
	parts := strings.SplitN(data.CustomID, ":", 3)
	if len(parts) != 3 {
		respond(s, i, "Something went wrong — please try creating the card again.")
		return
	}
	cardType, name := parts[1], parts[2]

	var summary string
	var body *string
	var tags []string
	for _, row := range data.Components {
		actionRow, ok := row.(*discordgo.ActionsRow)
		if !ok || len(actionRow.Components) == 0 {
			continue
		}
		input, ok := actionRow.Components[0].(*discordgo.TextInput)
		if !ok {
			continue
		}
		switch input.CustomID {
		case "summary":
			summary = input.Value
		case "body":
			if input.Value != "" {
				v := input.Value
				body = &v
			}
		case "tags":
			if input.Value != "" {
				raw := strings.Split(input.Value, ",")
				tags = make([]string, 0, len(raw))
				for _, t := range raw {
					t = strings.TrimSpace(t)
					if t != "" {
						tags = append(tags, t)
					}
				}
			}
		}
	}

	card, err := q.CreateCard(context.Background(), db.CreateCardParams{
		Name:    name,
		Type:    cardType,
		Summary: summary,
		Body:    body,
		Tags:    tags,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			respond(s, i, "A card named \""+name+"\" already exists.")
			return
		}
		log.Printf("create card error: %v", err)
		respond(s, i, "Something went wrong creating that card.")
		return
	}
	log.Printf("card created: id=%d name=%q type=%q by=%s", card.ID, card.Name, card.Type, i.Member.User.Username)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: card.Name + "** created. Reply here with an image within 60 seconds to add a portrait, or it'll be skipped.",
		},
	})

	awaitImageReply(s, i.ChannelID, i.Member.User.ID, card.ID, q)
}
