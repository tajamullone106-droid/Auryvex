package bot

import (
	"context"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/tajamullone106-droid/Auryvex/internal/config"
	"github.com/tajamullone106-droid/Auryvex/internal/room"
)

func Start(cfg *config.Config, rooms *room.Manager) {
	b, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Auryvex started as @%s", b.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "start":
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Welcome to Auryvex 🎵\n\nUse /room to create a synchronized room.",
			)

			if _, err := b.Send(msg); err != nil {
				log.Println(err)
			}

		case "room":
			name := strings.TrimSpace(update.Message.CommandArguments())
			if name == "" {
				name = "Auryvex Room"
			}

			ctx := context.Background()

			roomID, err := rooms.Create(
				ctx,
				int64(update.Message.From.ID),
				update.Message.Chat.ID,
				name,
			)
			if err != nil {
				log.Println("room creation failed:", err)

				msg := tgbotapi.NewMessage(
					update.Message.Chat.ID,
					"❌ Failed to create the room. Please try again.",
				)
				_, _ = b.Send(msg)
				continue
			}

			text := fmt.Sprintf(
				"🎧 *Room Created!*\n\n"+
					"Name: %s\n"+
					"Room ID: `%s`\n\n"+
					"Mini App support is coming next.",
				name,
				roomID.String(),
			)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ParseMode = "Markdown"

			if _, err := b.Send(msg); err != nil {
				log.Println(err)
			}
		}
	}
}
