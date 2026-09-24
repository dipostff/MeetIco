package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"meetico/internal/model"
	"meetico/internal/notify"
	"meetico/internal/repo"
)

var ErrForbidden = errors.New("forbidden")

type BookingService struct {
	bookingsRepo      *repo.BookingsRepo
	eventTypesRepo    *repo.EventTypesRepo
	usersRepo         *repo.UsersRepo
	notificationsRepo *repo.NotificationsRepo
	schedulesRepo     *repo.SchedulesRepo
	slotsService      *SlotsService
	mailer            notify.Mailer
	telegramNotifier  notify.TelegramNotifier
}

func NewBookingService(
	bookingsRepo *repo.BookingsRepo,
	eventTypesRepo *repo.EventTypesRepo,
	usersRepo *repo.UsersRepo,
	notificationsRepo *repo.NotificationsRepo,
	schedulesRepo *repo.SchedulesRepo,
	slotsService *SlotsService,
	mailer notify.Mailer,
	telegramNotifier notify.TelegramNotifier,
) *BookingService {
	return &BookingService{
		bookingsRepo:      bookingsRepo,
		eventTypesRepo:    eventTypesRepo,
		usersRepo:         usersRepo,
		notificationsRepo: notificationsRepo,
		schedulesRepo:     schedulesRepo,
		slotsService:      slotsService,
		mailer:            mailer,
		telegramNotifier:  telegramNotifier,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, eventTypeID, candidateName, candidateEmail, slotTime, dateStr, timezone string) (*model.Booking, error) {
	eventType, err := s.eventTypesRepo.GetByID(ctx, eventTypeID)
	if err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	slotTimeParsed, err := time.Parse("15:04", slotTime)
	if err != nil {
		return nil, err
	}

	user, err := s.usersRepo.GetByID(ctx, eventType.UserID)
	if err != nil {
		return nil, err
	}

	schedule, err := s.schedulesRepo.GetByUserID(ctx, eventType.UserID)
	if err != nil {
		return nil, err
	}

	recruiterLoc, err := time.LoadLocation(schedule.Timezone)
	if err != nil {
		return nil, err
	}

	startsAt := time.Date(date.Year(), date.Month(), date.Day(), slotTimeParsed.Hour(), slotTimeParsed.Minute(), 0, 0, recruiterLoc)
	endsAt := startsAt.Add(time.Duration(eventType.DurationMin) * time.Minute)

	booking := &model.Booking{
		EventTypeID:    eventTypeID,
		CandidateName:  candidateName,
		CandidateEmail: candidateEmail,
		StartsAt:       startsAt,
		EndsAt:         endsAt,
		Status:         "confirmed",
	}

	err = s.bookingsRepo.Create(ctx, booking)
	if err != nil {
		return nil, err
	}

	err = s.notificationsRepo.CreateForBooking(ctx, booking.ID, startsAt)
	if err != nil {
		slog.Error("failed to create notifications", "error", err)
	}

	s.slotsService.InvalidateCache(eventType.UserID, date)

	err = s.mailer.SendBookingConfirmation(candidateEmail, candidateName, user.DisplayName, startsAt, booking.CancelToken)
	if err != nil {
		slog.Error("failed to send confirmation email", "error", err)
	}

	if user.TelegramChatID != "" {
		msg := fmt.Sprintf(
			"📅 <b>Новая запись!</b>\n👤 %s\n📧 %s\n🕐 %s\n📋 %s",
			candidateName, candidateEmail, startsAt.Format("2 января в 15:04 МСК"), eventType.Title,
		)
		err = s.telegramNotifier.SendMessage(user.TelegramChatID, msg)
		if err != nil {
			slog.Error("failed to send telegram notification", "error", err)
		}
	}

	return booking, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, token string) error {
	booking, err := s.bookingsRepo.GetByCancelToken(ctx, token)
	if err != nil {
		return err
	}
	return s.cancelBooking(ctx, booking)
}

func (s *BookingService) CancelBookingByOwner(ctx context.Context, userID, bookingID string) error {
	booking, err := s.bookingsRepo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	eventType, err := s.eventTypesRepo.GetByID(ctx, booking.EventTypeID)
	if err != nil {
		return err
	}
	if eventType.UserID != userID {
		return ErrForbidden
	}

	return s.cancelBooking(ctx, booking)
}

func (s *BookingService) cancelBooking(ctx context.Context, booking *model.Booking) error {
	if booking.Status == "cancelled" {
		return nil
	}

	if err := s.bookingsRepo.UpdateStatus(ctx, booking.ID, "cancelled"); err != nil {
		return err
	}

	eventType, err := s.eventTypesRepo.GetByID(ctx, booking.EventTypeID)
	if err != nil {
		return err
	}

	user, err := s.usersRepo.GetByID(ctx, eventType.UserID)
	if err != nil {
		return err
	}

	if user.TelegramChatID != "" {
		msg := fmt.Sprintf(
			"❌ <b>Встреча отменена</b>\n👤 %s\n🕐 %s",
			booking.CandidateName, booking.StartsAt.Format("2 января в 15:04 МСК"),
		)
		if err := s.telegramNotifier.SendMessage(user.TelegramChatID, msg); err != nil {
			slog.Error("failed to send telegram notification", "error", err)
		}
	}

	return nil
}
