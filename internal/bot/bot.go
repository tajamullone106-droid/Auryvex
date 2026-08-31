package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/tajamullone106-droid/Auryvex/internal/config"
)

func Start(cfg *config.Config) {
	b, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Auryvex started as @%s", b.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() && update.Message.Command() == "start" {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				"Welcome to Auryvex 🎵\n\nUse /room to create a synchronized room.",
			)

			if _, err := b.Send(msg); err != nil {
				log.Println(err)
			}
		}
	}
}
