package repo

import (
	"context"
	"fmt"

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
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, timezone, weekly_hours FROM schedules WHERE user_id = $1`,
		userID,
	).Scan(&s.ID, &s.UserID, &s.Timezone, &s.WeeklyHours)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SchedulesRepo) Create(ctx context.Context, schedule *model.Schedule) error {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO schedules (user_id, timezone, weekly_hours) 
		 VALUES ($1, $2, $3) 
		 RETURNING id`,
		schedule.UserID, schedule.Timezone, schedule.WeeklyHours,
	).Scan(&schedule.ID)
	return err
}

func (r *SchedulesRepo) Update(ctx context.Context, userID string, timezone *string, weeklyHours *map[string][][]string) error {
	query := `UPDATE schedules SET `
	args := []interface{}{}
	argIdx := 1

	if timezone != nil {
		query += fmt.Sprintf("timezone = $%d, ", argIdx)
		args = append(args, *timezone)
		argIdx++
	}
	if weeklyHours != nil {
		query += fmt.Sprintf("weekly_hours = $%d, ", argIdx)
		args = append(args, *weeklyHours)
		argIdx++
	}

	if len(args) == 0 {
		return nil
	}

	query = query[:len(query)-2] + fmt.Sprintf(" WHERE user_id = $%d", argIdx)
	args = append(args, userID)

	_, err := r.pool.Exec(ctx, query, args...)
	return err
}
