package room

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type State struct {
	RoomID         string  `json:"room_id"`
	TrackID        string  `json:"track_id,omitempty"`
	TrackTitle     string  `json:"track_title,omitempty"`
	TrackURL       string  `json:"track_url,omitempty"`
	IsPlaying      bool    `json:"is_playing"`
	Position       float64 `json:"position"`
	StateUpdatedAt int64   `json:"state_updated_at"`
}

type Manager struct {
	DB    *pgxpool.Pool
	Redis *redis.Client
}

func NewManager(db *pgxpool.Pool, redisClient *redis.Client) *Manager {
	return &Manager{
		DB:    db,
		Redis: redisClient,
	}
}

func (m *Manager) Create(ctx context.Context, ownerID int64, chatID int64, name string) (uuid.UUID, error) {
	id := uuid.New()

	_, err := m.DB.Exec(
		ctx,
		`INSERT INTO rooms
			(id, chat_id, owner_telegram_id, name)
		 VALUES ($1, $2, $3, $4)`,
		id,
		chatID,
		ownerID,
		name,
	)
	if err != nil {
		return uuid.Nil, err
	}

	state := State{
		RoomID:         id.String(),
		IsPlaying:      false,
		Position:       0,
		StateUpdatedAt: time.Now().UnixMilli(),
	}

	if err := m.SetState(ctx, id, state); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (m *Manager) SetState(ctx context.Context, roomID uuid.UUID, state State) error {
	state.StateUpdatedAt = time.Now().UnixMilli()

	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("room:%s:state", roomID)

	return m.Redis.Set(ctx, key, data, 24*time.Hour).Err()
}

func (m *Manager) GetState(ctx context.Context, roomID uuid.UUID) (*State, error) {
	key := fmt.Sprintf("room:%s:state", roomID)

	data, err := m.Redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}

	return &state, nil
}
