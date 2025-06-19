package main

import (
	bot "discord-bot/Bot"
	db "discord-bot/Db"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN not found in environment")
	}

	RegisterCommands()

	database, err := db.NewDatabase("todos.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer func(database *db.Database) {
		err := database.Close()
		if err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			fmt.Println("Database connection closed successfully.")
		}
	}(database)

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	dg.AddHandler(ready)
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		interactionCreate(s, i, database)
	})

	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection:", err)
	}
	defer func(dg *discordgo.Session) {
		err := dg.Close()
		if err != nil {
			log.Printf("Error closing Discord session: %v", err)
		} else {
			fmt.Println("Discord session closed successfully.")
		}
	}(dg)

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("Gracefully shutting down...")
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Printf("Bot %s jest online!\n", s.State.User.Username)
	fmt.Printf("Bot jest na %d serwerach\n", len(event.Guilds))
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate, database *db.Database) {
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "todo":
			bot.HandleTodoCommand(s, i, database)
		}
	}

	if i.Type == discordgo.InteractionMessageComponent {
		if strings.HasPrefix(i.MessageComponentData().CustomID, "complete_todo_") {
			bot.HandleCompleteButton(s, i, database)
		}
	}
}
