package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"meetico/internal/notify"
	"meetico/internal/repo"
)

func Start(ctx context.Context, db *pgxpool.Pool, mailer notify.Mailer, tg notify.TelegramNotifier) {
	notificationsRepo := repo.NewNotificationsRepo(db)
	usersRepo := repo.NewUsersRepo(db)
	bookingsRepo := repo.NewBookingsRepo(db)
	eventTypesRepo := repo.NewEventTypesRepo(db)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	slog.Info("reminders worker started")

	for {
		select {
		case <-ctx.Done():
			slog.Info("reminders worker stopped")
			return
		case <-ticker.C:
			processNotifications(ctx, notificationsRepo, usersRepo, bookingsRepo, eventTypesRepo, mailer, tg)
		}
	}
}

func processNotifications(
	ctx context.Context,
	notificationsRepo *repo.NotificationsRepo,
	usersRepo *repo.UsersRepo,
	bookingsRepo *repo.BookingsRepo,
	eventTypesRepo *repo.EventTypesRepo,
	mailer notify.Mailer,
	tg notify.TelegramNotifier,
) {
	notifications, err := notificationsRepo.GetPending(ctx)
	if err != nil {
		slog.Error("failed to get pending notifications", "error", err)
		return
	}

	for _, notification := range notifications {
		booking, err := bookingsRepo.GetByID(ctx, notification.BookingID)
		if err != nil {
			slog.Error("failed to get booking", "notification_id", notification.ID, "error", err)
			continue
		}

		eventType, err := eventTypesRepo.GetByID(ctx, booking.EventTypeID)
		if err != nil {
			slog.Error("failed to get event type", "notification_id", notification.ID, "error", err)
			continue
		}

		user, err := usersRepo.GetByID(ctx, eventType.UserID)
		if err != nil {
			slog.Error("failed to get user", "notification_id", notification.ID, "error", err)
			continue
		}

		if notification.Channel == "email" {
			err = mailer.SendReminder(booking.CandidateEmail, booking.CandidateName, user.DisplayName, booking.StartsAt, booking.CancelToken)
			if err != nil {
				slog.Error("failed to send email reminder", "notification_id", notification.ID, "error", err)
				continue
			}
		}

		if notification.Channel == "telegram" && user.TelegramChatID != "" {
			err = tg.SendMessage(user.TelegramChatID, "Напоминание о встрече")
			if err != nil {
				slog.Error("failed to send telegram reminder", "notification_id", notification.ID, "error", err)
				continue
			}
		}

		err = notificationsRepo.MarkSent(ctx, notification.ID)
		if err != nil {
			slog.Error("failed to mark notification as sent", "notification_id", notification.ID, "error", err)
			continue
		}

		slog.Info("reminder sent", "type", notification.Type, "to", booking.CandidateEmail)
	}
}
