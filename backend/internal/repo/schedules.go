package repo

import (
	"context"
	"encoding/json"

	"meetico/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SchedulesRepo struct {
	pool *pgxpool.Pool
}

func NewSchedulesRepo(pool *pgxpool.Pool) *SchedulesRepo {
	return &SchedulesRepo{pool: pool}
}

func (r *SchedulesRepo) GetByUserID(ctx context.Context, userID string) (*model.Schedule, error) {
	var s model.Schedule
	var hoursJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, timezone, weekly_hours FROM schedules WHERE user_id = $1`,
		userID,
	).Scan(&s.ID, &s.UserID, &s.Timezone, &hoursJSON)
	if err != nil {
		return nil, err
	}
	s.WeeklyHours = map[string][][]string{}
	if len(hoursJSON) > 0 {
		if err := json.Unmarshal(hoursJSON, &s.WeeklyHours); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

func (r *SchedulesRepo) Create(ctx context.Context, schedule *model.Schedule) error {
	if schedule.WeeklyHours == nil {
		schedule.WeeklyHours = map[string][][]string{}
	}
	hoursJSON, err := json.Marshal(schedule.WeeklyHours)
	if err != nil {
		return err
	}
	err = r.pool.QueryRow(ctx,
		`INSERT INTO schedules (user_id, timezone, weekly_hours)
		 VALUES ($1, $2, $3::jsonb)
		 RETURNING id`,
		schedule.UserID, schedule.Timezone, hoursJSON,
	).Scan(&schedule.ID)
	return err
}

func (r *SchedulesRepo) Update(ctx context.Context, userID string, timezone *string, weeklyHours *map[string][][]string) error {
	tz := "Europe/Moscow"
	hours := map[string][][]string{}

	if existing, err := r.GetByUserID(ctx, userID); err == nil {
		tz = existing.Timezone
		if existing.WeeklyHours != nil {
			hours = existing.WeeklyHours
		}
	}

	if timezone != nil {
		tz = *timezone
	}
	if weeklyHours != nil {
		hours = *weeklyHours
	}

	hoursJSON, err := json.Marshal(hours)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx,
		`INSERT INTO schedules (user_id, timezone, weekly_hours)
		 VALUES ($1, $2, $3::jsonb)
		 ON CONFLICT (user_id) DO UPDATE SET
		   timezone = EXCLUDED.timezone,
		   weekly_hours = EXCLUDED.weekly_hours`,
		userID, tz, hoursJSON,
	)
	return err
}

// DefaultWeeklyHours is Mon–Fri 09:00–18:00.
func DefaultWeeklyHours() map[string][][]string {
	mk := func() [][]string { return [][]string{{"09:00", "18:00"}} }
	return map[string][][]string{
		"mon": mk(),
		"tue": mk(),
		"wed": mk(),
		"thu": mk(),
		"fri": mk(),
		"sat": {},
		"sun": {},
	}
}

// EnsureDefault creates a default schedule if the user has none.
func (r *SchedulesRepo) EnsureDefault(ctx context.Context, userID string) error {
	_, err := r.GetByUserID(ctx, userID)
	if err == nil {
		return nil
	}
	return r.Create(ctx, &model.Schedule{
		UserID:      userID,
		Timezone:    "Europe/Moscow",
		WeeklyHours: DefaultWeeklyHours(),
	})
}
