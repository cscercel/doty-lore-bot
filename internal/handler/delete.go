package handler

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"

	"github.com/cscercel/doty-lore-bot/internal/database"
)

func HandleDelete(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	opts []*discordgo.ApplicationCommandInteractionDataOption,
	q *database.Queries,
) {
	name := opts[0].StringValue()

	existing, err := q.GetCardByName(context.Background(), name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond(s, i, "No card found named \""+name+"\".")
			return
		}
		log.Printf("delete lookup error: %v", err)
		respond(s, i, "Something went wrong looking that up.")
		return
	}

	if err := q.DeleteCard(context.Background(), existing.ID); err != nil {
		log.Printf("delete card error: %v", err)
		respond(s, i, "Something went wrong deleting that card.")
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Deleted **" + existing.Name + "**.",
		},
	})

	log.Printf("card deleted: name=%q by=%s", existing.Name, i.Member.User.Username)
}
