package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"meetico/internal/model"
)

type EventTypesRepo struct {
	pool *pgxpool.Pool
}

func NewEventTypesRepo(pool *pgxpool.Pool) *EventTypesRepo {
	return &EventTypesRepo{pool: pool}
}

func (r *EventTypesRepo) Create(ctx context.Context, eventType *model.EventType) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO event_types (user_id, slug, title, duration_min, location_type, location_value, category) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7) 
		 RETURNING id, created_at`,
		eventType.UserID, eventType.Slug, eventType.Title, eventType.DurationMin,
		eventType.LocationType, eventType.LocationValue, eventType.Category,
	).Scan(&eventType.ID, &eventType.CreatedAt)
	return err
}

func (r *EventTypesRepo) GetByUserID(ctx context.Context, userID string) ([]*model.EventType, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, slug, title, duration_min, location_type, location_value, category, is_active, created_at 
		 FROM event_types WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var eventTypes []*model.EventType
	for rows.Next() {
		var et model.EventType
		err := rows.Scan(
			&et.ID, &et.UserID, &et.Slug, &et.Title, &et.DurationMin,
			&et.LocationType, &et.LocationValue, &et.Category, &et.IsActive, &et.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		eventTypes = append(eventTypes, &et)
	}
	return eventTypes, nil
}

func (r *EventTypesRepo) GetByID(ctx context.Context, id string) (*model.EventType, error) {
	var et model.EventType
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, slug, title, duration_min, location_type, location_value, category, is_active, created_at 
		 FROM event_types WHERE id = $1`,
		id,
	).Scan(
		&et.ID, &et.UserID, &et.Slug, &et.Title, &et.DurationMin,
		&et.LocationType, &et.LocationValue, &et.Category, &et.IsActive, &et.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &et, nil
}

func (r *EventTypesRepo) GetBySlug(ctx context.Context, slug string) (*model.EventType, error) {
	var et model.EventType
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, slug, title, duration_min, location_type, location_value, category, is_active, created_at 
		 FROM event_types WHERE slug = $1`,
		slug,
	).Scan(
		&et.ID, &et.UserID, &et.Slug, &et.Title, &et.DurationMin,
		&et.LocationType, &et.LocationValue, &et.Category, &et.IsActive, &et.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &et, nil
}

func (r *EventTypesRepo) Update(ctx context.Context, id string, title *string, durationMin *int, locationType, locationValue *string, isActive *bool) error {
	query := `UPDATE event_types SET `
	args := []interface{}{}
	argIdx := 1

	if title != nil {
		query += fmt.Sprintf("title = $%d, ", argIdx)
		args = append(args, *title)
		argIdx++
	}
	if durationMin != nil {
		query += fmt.Sprintf("duration_min = $%d, ", argIdx)
		args = append(args, *durationMin)
		argIdx++
	}
	if locationType != nil {
		query += fmt.Sprintf("location_type = $%d, ", argIdx)
		args = append(args, *locationType)
		argIdx++
	}
	if locationValue != nil {
		query += fmt.Sprintf("location_value = $%d, ", argIdx)
		args = append(args, *locationValue)
		argIdx++
	}
	if isActive != nil {
		query += fmt.Sprintf("is_active = $%d, ", argIdx)
		args = append(args, *isActive)
		argIdx++
	}

	if len(args) == 0 {
		return nil
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}

func (r *EventTypesRepo) Delete(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE event_types SET is_active = false WHERE id = $1`,
		id,
	)
	return err
}
