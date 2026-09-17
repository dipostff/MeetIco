package model

import (
	"time"
)

type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	DisplayName    string    `json:"display_name"`
	Username       string    `json:"username"`
	Bio            string    `json:"bio"`
	PhotoURL       string    `json:"photo_url"`
	TelegramChatID string    `json:"telegram_chat_id"`
	Plan           string    `json:"plan"`
	Timezone       string    `json:"timezone"`
	CreatedAt      time.Time `json:"created_at"`
}

type Schedule struct {
	ID          string                `json:"id"`
	UserID      string                `json:"user_id"`
	Timezone    string                `json:"timezone"`
	WeeklyHours map[string][][]string `json:"weekly_hours"`
}

type EventType struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Slug          string    `json:"slug"`
	Title         string    `json:"title"`
	DurationMin   int       `json:"duration_min"`
	LocationType  string    `json:"location_type"`
	LocationValue string    `json:"location_value"`
	Category      string    `json:"category"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

type Booking struct {
	ID             string    `json:"id"`
	EventTypeID    string    `json:"event_type_id"`
	CandidateName  string    `json:"candidate_name"`
	CandidateEmail string    `json:"candidate_email"`
	StartsAt       time.Time `json:"starts_at"`
	EndsAt         time.Time `json:"ends_at"`
	Status         string    `json:"status"`
	CancelToken    string    `json:"cancel_token"`
	CreatedAt      time.Time `json:"created_at"`
}

type Notification struct {
	ID        string    `json:"id"`
	BookingID string    `json:"booking_id"`
	Channel   string    `json:"channel"`
	Type      string    `json:"type"`
	SendAt    time.Time `json:"send_at"`
	Sent      bool      `json:"sent"`
}
