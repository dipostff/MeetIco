package repo

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"meetico/internal/model"
)

type NotificationsRepo struct {
	pool *pgxpool.Pool
}

func NewNotificationsRepo(pool *pgxpool.Pool) *NotificationsRepo {
	return &NotificationsRepo{pool: pool}
}

func (r *NotificationsRepo) Create(ctx context.Context, notification *model.Notification) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO notifications (booking_id, channel, type, send_at) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id`,
		notification.BookingID, notification.Channel, notification.Type, notification.SendAt,
	).Scan(&notification.ID)
	return err
}

func (r *NotificationsRepo) GetPending(ctx context.Context) ([]*model.Notification, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, booking_id, channel, type, send_at, sent 
		 FROM notifications 
		 WHERE send_at <= NOW() AND sent = FALSE 
		 FOR UPDATE SKIP LOCKED`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		var n model.Notification
		err := rows.Scan(
			&n.ID, &n.BookingID, &n.Channel, &n.Type, &n.SendAt, &n.Sent,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}
	return notifications, nil
}

func (r *NotificationsRepo) MarkSent(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE notifications SET sent = TRUE WHERE id = $1`,
		id,
	)
	return err
}

func (r *NotificationsRepo) CreateForBooking(ctx context.Context, bookingID string, startsAt time.Time) error {
	notifications := []struct {
		channel string
		typ     string
		sendAt  time.Time
	}{
		{"email", "confirmation", time.Now()},
		{"email", "reminder_24h", startsAt.Add(-24 * time.Hour)},
		{"email", "reminder_1h", startsAt.Add(-1 * time.Hour)},
	}

	for _, n := range notifications {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO notifications (booking_id, channel, type, send_at) 
			 VALUES ($1, $2, $3, $4)`,
			bookingID, n.channel, n.typ, n.sendAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
