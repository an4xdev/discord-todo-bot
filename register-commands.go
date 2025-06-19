package main

import (
    "log"
    "os"

    "github.com/bwmarrin/discordgo"
    "github.com/joho/godotenv"
)

func RegisterCommands() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    token := os.Getenv("DISCORD_TOKEN")
    clientID := os.Getenv("CLIENT_ID")

    if token == "" || clientID == "" {
        log.Fatal("DISCORD_TOKEN and CLIENT_ID must be set in environment")
    }

    dg, err := discordgo.New("Bot " + token)
    if err != nil {
        log.Fatal("Error creating Discord session:", err)
    }

    commands := []*discordgo.ApplicationCommand{
        {
            Name:        "todo",
            Description: "Zarządzanie zadaniami TODO",
            Options: []*discordgo.ApplicationCommandOption{
                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "add",
                    Description: "Dodaj nowe zadanie",
                    Options: []*discordgo.ApplicationCommandOption{
                        {
                            Type:        discordgo.ApplicationCommandOptionString,
                            Name:        "task",
                            Description: "Treść zadania",
                            Required:    true,
                        },
                    },
                },
                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "list",
                    Description: "Pokaż wszystkie zadania na tym kanale",
                },
                {
                    Type:        discordgo.ApplicationCommandOptionSubCommand,
                    Name:        "reset",
                    Description: "Usuń wszystkie zadania z tego kanału",
                },
            },
        },
    }

    log.Println("Rozpoczęcie rejestracji komend...")

    for _, command := range commands {
        _, err := dg.ApplicationCommandCreate(clientID, "", command)
        if err != nil {
            log.Fatalf("Error creating command %s: %v", command.Name, err)
        }
    }

    log.Println("Komendy zostały zarejestrowane!")
}
