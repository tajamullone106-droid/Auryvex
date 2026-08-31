package main

import (
	"fmt"
	"log"

	"github.com/tajamullone106-droid/Auryvex/internal/config"
	"github.com/tajamullone106-droid/Auryvex/internal/music"
)

func main() {
	cfg := config.Load()
	client := music.NewClient(cfg)

	results, err := client.Search("Alan Walker Faded", 3)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ArcMusic returned %d results:\n\n", len(results))

	for i, r := range results {
		fmt.Printf("%d. %s\n", i+1, r.Title)
		fmt.Printf("   Video ID: %s\n", r.VideoID)
		fmt.Printf("   Duration: %s\n", r.Duration)
		fmt.Printf("   Channel: %s\n", r.Channel)
		fmt.Printf("   URL: %s\n\n", r.URL)
	}
}
