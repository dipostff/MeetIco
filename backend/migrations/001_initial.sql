-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email            TEXT NOT NULL UNIQUE,
  password_hash    TEXT NOT NULL,
  display_name     TEXT,
  username         TEXT UNIQUE,
  bio              TEXT,
  photo_url        TEXT,
  telegram_chat_id TEXT,
  plan             TEXT NOT NULL DEFAULT 'free',
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX users_username_idx ON users(username)
  WHERE username IS NOT NULL;

CREATE TABLE schedules (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  timezone     TEXT NOT NULL DEFAULT 'Europe/Moscow',
  weekly_hours JSONB NOT NULL DEFAULT '{}',
  UNIQUE(user_id)
);

CREATE TABLE event_types (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  slug           TEXT NOT NULL,
  title          TEXT NOT NULL,
  duration_min   INTEGER NOT NULL DEFAULT 30,
  location_type  TEXT NOT NULL DEFAULT 'zoom',
  location_value TEXT,
  category       TEXT NOT NULL DEFAULT 'other',
  is_active      BOOLEAN NOT NULL DEFAULT TRUE,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(user_id, slug)
);

CREATE TABLE bookings (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_type_id   UUID NOT NULL REFERENCES event_types(id) ON DELETE CASCADE,
  candidate_name  TEXT NOT NULL,
  candidate_email TEXT NOT NULL,
  starts_at       TIMESTAMPTZ NOT NULL,
  ends_at         TIMESTAMPTZ NOT NULL,
  status          TEXT NOT NULL DEFAULT 'confirmed',
  cancel_token    UUID NOT NULL DEFAULT gen_random_uuid(),
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX bookings_event_type_starts_idx ON bookings(event_type_id, starts_at);
CREATE INDEX bookings_cancel_token_idx      ON bookings(cancel_token);
CREATE INDEX event_types_slug_idx           ON event_types(slug);

-- +goose Down
DROP INDEX IF EXISTS event_types_slug_idx;
DROP INDEX IF EXISTS bookings_cancel_token_idx;
DROP INDEX IF EXISTS bookings_event_type_starts_idx;
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS event_types;
DROP TABLE IF EXISTS schedules;
DROP INDEX IF EXISTS users_username_idx;
DROP TABLE IF EXISTS users;
DROP EXTENSION IF EXISTS "pgcrypto";
