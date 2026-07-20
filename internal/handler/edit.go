package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v5"

	"github.com/cscercel/doty-lore-bot/internal/db"
)

func HandleEditStart(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	opts []*discordgo.ApplicationCommandInteractionDataOption,
	q *db.Queries,
) {
	name := opts[0].StringValue()

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

	bodyVal := ""
	if existing.Body != nil {
		bodyVal = *existing.Body
	}
	tagsVal := ""
	if len(existing.Tags) > 0 {
		tagsVal = strings.Join(existing.Tags, ", ")
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "lore_edit_modal:" + fmt.Sprint(existing.ID),
			Title:    "Edit: " + existing.Name,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:  "name",
						Label:     "Name",
						Style:     discordgo.TextInputShort,
						Required:  true,
						MaxLength: 100,
						Value:     existing.Name,
					},
				}},
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:  "summary",
						Label:     "Summary",
						Style:     discordgo.TextInputShort,
						Required:  true,
						MaxLength: 200,
						Value:     existing.Summary,
					},
				}},
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:  "body",
						Label:     "Full Details",
						Style:     discordgo.TextInputParagraph,
						Required:  false,
						MaxLength: 4000,
						Value:     bodyVal,
					},
				}},
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:  "tags",
						Label:     "Tags (comma-separated)",
						Style:     discordgo.TextInputShort,
						Required:  false,
						MaxLength: 200,
						Value:     tagsVal,
					},
				}},
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "image",
						Label:       "Image: keep / edit / delete",
						Style:       discordgo.TextInputShort,
						Required:    true,
						MaxLength:   10,
						Value:       "keep",
						Placeholder: "keep, edit, or delete",
					},
				}},
			},
		},
	})
	if err != nil {
		log.Printf("edit modal open error: %v", err)
	}
}

func HandleEditSubmit(s *discordgo.Session, i *discordgo.InteractionCreate, q *db.Queries) {
	data := i.ModalSubmitData()
	parts := strings.SplitN(data.CustomID, ":", 2)
	if len(parts) != 2 {
		respond(s, i, "Something went wrong — please try editing again.")
		return
	}
	id, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		respond(s, i, "Something went wrong — please try editing again.")
		return
	}

	var newName, summary, imageFlag string
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
		case "name":
			newName = strings.TrimSpace(input.Value)
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
		case "image":
			imageFlag = strings.ToLower(strings.TrimSpace(input.Value))
		}
	}

	if imageFlag != "keep" && imageFlag != "edit" && imageFlag != "delete" {
		respond(s, i, "Image field must be \"keep\", \"edit\", or \"delete\" — please try again.")
		return
	}

	// If the name changed, make sure it's not already taken by another card.
	if conflict, err := q.GetCardByName(context.Background(), newName); err == nil && conflict.ID != int32(id) {
		respond(s, i, "A card named \""+newName+"\" already exists — pick a different name.")
		return
	} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Printf("name conflict check error: %v", err)
		respond(s, i, "Something went wrong checking that name.")
		return
	}

	params := db.UpdateCardParams{
		ID:      int32(id),
		Name:    newName,
		Summary: summary,
		Body:    body,
		Tags:    tags,
	}

	if imageFlag == "delete" {
		params.ImageUrl = nil // clears the image
	}
	// for "keep" and "edit" we leave ImageUrl untouched here — see note below

	updated, err := q.UpdateCard(context.Background(), params)
	if err != nil {
		log.Printf("update card error: %v", err)
		respond(s, i, "Something went wrong updating that card.")
		return
	}
	log.Printf("card edited: name=%q by=%s", updated.Name, i.Member.User.Username)

	switch imageFlag {
	case "edit":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "**" + updated.Name + "** updated. Reply here with a new image within 60 seconds.",
			},
		})
		awaitImageReply(s, i.ChannelID, i.Member.User.ID, updated.ID, q)

	case "delete":
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "**" + updated.Name + "** updated. Image removed.",
			},
		})

	default: // "keep"
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "**" + updated.Name + "** updated.",
			},
		})
	}
}
