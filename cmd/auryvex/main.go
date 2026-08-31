package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/tajamullone106-droid/Auryvex/internal/bot"
	"github.com/tajamullone106-droid/Auryvex/internal/config"
	"github.com/tajamullone106-droid/Auryvex/internal/db"
	"github.com/tajamullone106-droid/Auryvex/internal/room"
	"github.com/tajamullone106-droid/Auryvex/internal/ws"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("database connection failed:", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("database ping failed:", err)
	}

	redisOptions, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatal("invalid Redis URL:", err)
	}

	redisClient := redis.NewClient(redisOptions)
	defer redisClient.Close()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	roomManager := room.NewManager(pool, redisClient)
	hub := ws.NewHub(roomManager)

	go startServer(hub)

	log.Println("Database connected")
	log.Println("Redis connected")
	log.Println("Starting Auryvex...")

	bot.Start(cfg, roomManager)
}
