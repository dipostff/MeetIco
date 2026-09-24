package repo

import (
	"context"
	"fmt"

	"meetico/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepo struct {
	pool *pgxpool.Pool
}

func NewUsersRepo(pool *pgxpool.Pool) *UsersRepo {
	return &UsersRepo{pool: pool}
}

func (r *UsersRepo) Create(ctx context.Context, email, passwordHash, displayName string) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, display_name)
		 VALUES ($1, $2, $3)
		 RETURNING id, email, password_hash, display_name,
		          COALESCE(username, ''), COALESCE(bio, ''), COALESCE(photo_url, ''),
		          COALESCE(telegram_chat_id, ''), plan, created_at`,
		email, passwordHash, displayName,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName,
		&user.Username, &user.Bio, &user.PhotoURL, &user.TelegramChatID,
		&user.Plan, &user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx,
		`SELECT u.id, u.email, u.password_hash, u.display_name,
		        COALESCE(u.username, ''), COALESCE(u.bio, ''), COALESCE(u.photo_url, ''),
		        COALESCE(u.telegram_chat_id, ''), u.plan, u.created_at,
		        COALESCE(s.timezone, 'Europe/Moscow')
		 FROM users u
		 LEFT JOIN schedules s ON s.user_id = u.id
		 WHERE u.id = $1`,
		id,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName,
		&user.Username, &user.Bio, &user.PhotoURL, &user.TelegramChatID,
		&user.Plan, &user.CreatedAt, &user.Timezone,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx,
		`SELECT u.id, u.email, u.password_hash, u.display_name,
		        COALESCE(u.username, ''), COALESCE(u.bio, ''), COALESCE(u.photo_url, ''),
		        COALESCE(u.telegram_chat_id, ''), u.plan, u.created_at,
		        COALESCE(s.timezone, 'Europe/Moscow')
		 FROM users u
		 LEFT JOIN schedules s ON s.user_id = u.id
		 WHERE u.email = $1`,
		email,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName,
		&user.Username, &user.Bio, &user.PhotoURL, &user.TelegramChatID,
		&user.Plan, &user.CreatedAt, &user.Timezone,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.pool.QueryRow(ctx,
		`SELECT u.id, u.email, u.password_hash, u.display_name,
		        COALESCE(u.username, ''), COALESCE(u.bio, ''), COALESCE(u.photo_url, ''),
		        COALESCE(u.telegram_chat_id, ''), u.plan, u.created_at,
		        COALESCE(s.timezone, 'Europe/Moscow')
		 FROM users u
		 LEFT JOIN schedules s ON s.user_id = u.id
		 WHERE u.username = $1`,
		username,
	).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName,
		&user.Username, &user.Bio, &user.PhotoURL, &user.TelegramChatID,
		&user.Plan, &user.CreatedAt, &user.Timezone,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepo) Update(ctx context.Context, id string, displayName, username, bio, photoURL, telegramChatID *string) error {
	query := `UPDATE users SET `
	args := []interface{}{}
	argIdx := 1

	if displayName != nil {
		query += fmt.Sprintf("display_name = $%d, ", argIdx)
		args = append(args, *displayName)
		argIdx++
	}
	if username != nil {
		query += fmt.Sprintf("username = $%d, ", argIdx)
		args = append(args, *username)
		argIdx++
	}
	if bio != nil {
		query += fmt.Sprintf("bio = $%d, ", argIdx)
		args = append(args, *bio)
		argIdx++
	}
	if photoURL != nil {
		query += fmt.Sprintf("photo_url = $%d, ", argIdx)
		args = append(args, *photoURL)
		argIdx++
	}
	if telegramChatID != nil {
		query += fmt.Sprintf("telegram_chat_id = $%d, ", argIdx)
		args = append(args, *telegramChatID)
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

func (r *UsersRepo) UpdatePassword(ctx context.Context, id, newPasswordHash string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET password_hash = $1 WHERE id = $2`,
		newPasswordHash, id,
	)
	return err
}
