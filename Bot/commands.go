package bot

import (
	db "discord-bot/Db"
	"fmt"
	"log"
	_ "strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HandleTodoCommand(s *discordgo.Session, i *discordgo.InteractionCreate, database *db.Database) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		return
	}

	subcommand := options[0].Name

	switch subcommand {
	case "add":
		handleAddTodo(s, i, options[0].Options, database)
	case "list":
		handleListTodos(s, i, database)
	case "reset":
		handleResetTodos(s, i, database)
	}
}

func handleAddTodo(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, database *db.Database) {
	if len(options) == 0 {
		return
	}

	task := options[0].StringValue()
	channelID := i.ChannelID

	embed := &discordgo.MessageEmbed{
		Title:       "📝 Nowe TODO",
		Description: task,
		Color:       0x0099ff,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: "complete_todo_temp",
					Label:    "✅ Oznacz jako wykonane",
					Style:    discordgo.SuccessButton,
				},
			},
		},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})

	if err != nil {
		log.Printf("Error responding to interaction: %v", err)
		return
	}

	response, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		log.Printf("Error getting interaction response: %v", err)
		return
	}

	todoID, err := database.InsertTodo(channelID, response.ID, task)
	if err != nil {
		log.Printf("Error inserting todo: %v", err)
		return
	}

	updatedComponents := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: fmt.Sprintf("complete_todo_%d", todoID),
					Label:    "✅ Oznacz jako wykonane",
					Style:    discordgo.SuccessButton,
				},
			},
		},
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Components: &updatedComponents,
	})

	if err != nil {
		log.Printf("Error updating interaction: %v", err)
	}
}

func handleListTodos(s *discordgo.Session, i *discordgo.InteractionCreate, database *db.Database) {
	channelID := i.ChannelID

	todos, err := database.GetTodosByChannel(channelID)
	if err != nil {
		log.Printf("Error getting todos: %v", err)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Błąd podczas pobierania zadań!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Printf("Error responding to list command: %v", err)
			return
		}
		return
	}

	if len(todos) == 0 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Brak zadań na tym kanale!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Printf("Error responding to empty list command: %v", err)
			return
		}
		return
	}

	var todoList strings.Builder
	for _, todo := range todos {
		var status, text string
		if todo.Completed {
			status = "✅"
			text = fmt.Sprintf("~~%s~~", todo.Content)
		} else {
			status = "❌"
			text = todo.Content
		}
		todoList.WriteString(fmt.Sprintf("%s %s\n", status, text))
	}

	embed := &discordgo.MessageEmbed{
		Title:       "📋 Lista TODO",
		Description: todoList.String(),
		Color:       0x00ae86,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})

	if err != nil {
		log.Printf("Error responding to list command: %v", err)
		return
	}

	response, err := s.InteractionResponse(i.Interaction)
	if err != nil {
		log.Printf("Error getting interaction response: %v", err)
		return
	}

	err = database.InsertListMessage(channelID, response.ID)
	if err != nil {
		log.Printf("Error inserting list message: %v", err)
	}
}

func handleResetTodos(s *discordgo.Session, i *discordgo.InteractionCreate, database *db.Database) {
	channelID := i.ChannelID

	todoMessageIDs, err := database.GetTodoMessageIDs(channelID)
	if err != nil {
		log.Printf("Error getting todo message IDs: %v", err)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Błąd podczas pobierania zadań!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Printf("Error responding to reset command: %v", err)
			return
		}
		return
	}

	listMessageIDs, err := database.GetListMessageIDs(channelID)
	if err != nil {
		log.Printf("Error getting list message IDs: %v", err)
	}

	deletedMessages := 0

	for _, messageID := range todoMessageIDs {
		err := s.ChannelMessageDelete(channelID, messageID)
		if err != nil {
			log.Printf("Cannot delete message %s: %v", messageID, err)
		} else {
			deletedMessages++
		}
	}

	for _, messageID := range listMessageIDs {
		err := s.ChannelMessageDelete(channelID, messageID)
		if err != nil {
			log.Printf("Cannot delete list message %s: %v", messageID, err)
		} else {
			deletedMessages++
		}
	}

	rowsAffected, err := database.DeleteTodosByChannel(channelID)
	if err != nil {
		log.Printf("Error deleting todos from database: %v", err)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Błąd podczas resetowania zadań!",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Printf("Error responding to reset command: %v", err)
			return
		}
		return
	}

	err = database.DeleteListMessagesByChannel(channelID)
	if err != nil {
		log.Printf("Error deleting list messages from database: %v", err)
	}

	embed := &discordgo.MessageEmbed{
		Title:       "🗑️ Reset TODO",
		Description: fmt.Sprintf("Usunięto %d zadań z bazy danych i %d wiadomości z kanału.", rowsAffected, deletedMessages),
		Color:       0xff6b6b,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		log.Printf("Error responding to reset command: %v", err)
		return
	}
}

func HandleCompleteButton(s *discordgo.Session, i *discordgo.InteractionCreate, database *db.Database) {
	customID := i.MessageComponentData().CustomID
	parts := strings.Split(customID, "_")
	if len(parts) != 3 {
		return
	}

	todoID := parts[2]
	todo, err := database.CompleteTodo(todoID)
	if err != nil {
		log.Printf("Error completing todo: %v", err)
		return
	}

	completedEmbed := &discordgo.MessageEmbed{
		Title:       "✅ TODO - Wykonane",
		Description: fmt.Sprintf("~~%s~~", todo.Content),
		Color:       0x90ee90,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{completedEmbed},
			Components: []discordgo.MessageComponent{}, // Remove components
		},
	})

	if err != nil {
		log.Printf("Error updating completed todo: %v", err)
	}
}
