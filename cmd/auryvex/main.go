package main

import (
	"log"

	"github.com/tajamullone106-driod/Auryvex/internal/bot"
	"github.com/tajamullone106-driod/Auryvex/internal/config"
)

func main() {
	cfg := config.Load()

	log.Println("Starting Auryvex...")
	bot.Start(cfg)
}
