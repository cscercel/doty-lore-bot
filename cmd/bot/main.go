package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"

	"github.com/cscercel/doty-lore-bot/internal/config"
	"github.com/cscercel/doty-lore-bot/internal/db"
	"github.com/cscercel/doty-lore-bot/internal/repository"
	"github.com/cscercel/doty-lore-bot/internal/discord"
)

func main() {
	// Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v:", err)
	}

	ctx := context.Background()

	// Connect to Database
	pool, err := repository.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Load queries
	queries := db.New(pool)

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
