package main

import (
	"context"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/cscercel/doty-lore-bot/internal/repository"
)

func main() {
	// Config
	godotenv.Load()
	db_url := os.Getenv("DATABASE_URL")
	token := os.Getenv("TOKEN")
	ctx := context.Background()

	// Connect to Database
	pool, err := repository.NewPool(ctx, db_url)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Connect bot to server
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
}
