package repo

import (
	"context"
	"time"

	"meetico/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingsRepo struct {
	pool *pgxpool.Pool
}

func NewBookingsRepo(pool *pgxpool.Pool) *BookingsRepo {
	return &BookingsRepo{pool: pool}
}

func (r *BookingsRepo) Create(ctx context.Context, booking *model.Booking) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO bookings (event_type_id, candidate_name, candidate_email, starts_at, ends_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, cancel_token, created_at`,
		booking.EventTypeID, booking.CandidateName, booking.CandidateEmail,
		booking.StartsAt, booking.EndsAt,
	).Scan(&booking.ID, &booking.CancelToken, &booking.CreatedAt)
	return err
}

func (r *BookingsRepo) GetByUserID(ctx context.Context, userID string) ([]*model.Booking, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT b.id, b.event_type_id, b.candidate_name, b.candidate_email, b.starts_at, b.ends_at, b.status, b.cancel_token, b.created_at 
		 FROM bookings b
		 JOIN event_types et ON et.id = b.event_type_id
		 WHERE et.user_id = $1 AND b.status = 'confirmed' AND b.starts_at >= NOW()
		 ORDER BY b.starts_at ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]*model.Booking, 0)
	for rows.Next() {
		var b model.Booking
		err := rows.Scan(
			&b.ID, &b.EventTypeID, &b.CandidateName, &b.CandidateEmail,
			&b.StartsAt, &b.EndsAt, &b.Status, &b.CancelToken, &b.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &b)
	}
	return bookings, nil
}

func (r *BookingsRepo) GetByEventTypeAndDate(ctx context.Context, eventTypeID string, date time.Time) ([]*model.Booking, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	rows, err := r.pool.Query(ctx,
		`SELECT id, event_type_id, candidate_name, candidate_email, starts_at, ends_at, status, cancel_token, created_at 
		 FROM bookings 
		 WHERE event_type_id = $1 AND starts_at >= $2 AND starts_at < $3 AND status = 'confirmed'`,
		eventTypeID, startOfDay, endOfDay,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]*model.Booking, 0)
	for rows.Next() {
		var b model.Booking
		err := rows.Scan(
			&b.ID, &b.EventTypeID, &b.CandidateName, &b.CandidateEmail,
			&b.StartsAt, &b.EndsAt, &b.Status, &b.CancelToken, &b.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &b)
	}
	return bookings, nil
}

func (r *BookingsRepo) GetByUserIDAndDate(ctx context.Context, userID string, date time.Time) ([]*model.Booking, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	rows, err := r.pool.Query(ctx,
		`SELECT b.id, b.event_type_id, b.candidate_name, b.candidate_email, b.starts_at, b.ends_at, b.status, b.cancel_token, b.created_at 
		 FROM bookings b
		 JOIN event_types et ON et.id = b.event_type_id
		 WHERE et.user_id = $1 AND b.starts_at >= $2 AND b.starts_at < $3 AND b.status = 'confirmed'`,
		userID, startOfDay, endOfDay,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]*model.Booking, 0)
	for rows.Next() {
		var b model.Booking
		err := rows.Scan(
			&b.ID, &b.EventTypeID, &b.CandidateName, &b.CandidateEmail,
			&b.StartsAt, &b.EndsAt, &b.Status, &b.CancelToken, &b.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, &b)
	}
	return bookings, nil
}

func (r *BookingsRepo) GetByCancelToken(ctx context.Context, token string) (*model.Booking, error) {
	var b model.Booking
	err := r.pool.QueryRow(ctx,
		`SELECT id, event_type_id, candidate_name, candidate_email, starts_at, ends_at, status, cancel_token, created_at 
		 FROM bookings WHERE cancel_token = $1`,
		token,
	).Scan(
		&b.ID, &b.EventTypeID, &b.CandidateName, &b.CandidateEmail,
		&b.StartsAt, &b.EndsAt, &b.Status, &b.CancelToken, &b.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookingsRepo) GetByID(ctx context.Context, id string) (*model.Booking, error) {
	var b model.Booking
	err := r.pool.QueryRow(ctx,
		`SELECT id, event_type_id, candidate_name, candidate_email, starts_at, ends_at, status, cancel_token, created_at 
		 FROM bookings WHERE id = $1`,
		id,
	).Scan(
		&b.ID, &b.EventTypeID, &b.CandidateName, &b.CandidateEmail,
		&b.StartsAt, &b.EndsAt, &b.Status, &b.CancelToken, &b.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookingsRepo) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE bookings SET status = $1 WHERE id = $2`,
		status, id,
	)
	return err
}
