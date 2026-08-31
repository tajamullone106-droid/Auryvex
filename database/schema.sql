CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username TEXT,
    first_name TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rooms (
    id UUID PRIMARY KEY,
    chat_id BIGINT,
    owner_telegram_id BIGINT NOT NULL,
    name TEXT NOT NULL DEFAULT 'Auryvex Room',

    current_track_id TEXT,
    current_track_title TEXT,
    current_track_url TEXT,

    is_playing BOOLEAN NOT NULL DEFAULT FALSE,
    position_seconds DOUBLE PRECISION NOT NULL DEFAULT 0,

    state_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS room_members (
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    telegram_id BIGINT NOT NULL REFERENCES users(telegram_id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (room_id, telegram_id)
);

CREATE TABLE IF NOT EXISTS queue (
    id BIGSERIAL PRIMARY KEY,
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    track_id TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    duration_seconds INTEGER NOT NULL DEFAULT 0,
    position INTEGER NOT NULL,
    added_by BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_room_members_room
    ON room_members(room_id);

CREATE INDEX IF NOT EXISTS idx_queue_room_position
    ON queue(room_id, position);
