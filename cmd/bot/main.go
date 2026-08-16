package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/tursodatabase/libsql-client-go/libsql"

	"github.com/cscercel/doty-lore-bot/internal/config"
	"github.com/cscercel/doty-lore-bot/internal/database"
	"github.com/cscercel/doty-lore-bot/internal/discord"
)

func main() {
	// Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v:", err)
	}

	// Sync to Database
	db, err := sql.Open("libsql", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Load queries
	queries := database.New(db)

	// Connect bot to server
	session, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatalf("failed to create discord session: %v", err)
	}

	// Load Handlers
	discord.RegisterHandlers(session, queries)

	if err := session.Open(); err != nil {
		log.Fatalf("failed to open discord connection: %v", err)
	}
	defer session.Close()

	registeredCmds, err := discord.RegisterCommands(session, cfg.GuildID)
	if err != nil {
		log.Fatalf("failed to register commands: %v", err)
	}
	defer discord.CleanupCommands(session, cfg.GuildID, registeredCmds)

	log.Println("bot is running, press CTRL-C to exit")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}
