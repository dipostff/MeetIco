package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"meetico/internal/config"
	"meetico/internal/db"
	"meetico/internal/handler"
	"meetico/internal/middleware"
	"meetico/internal/notify"
	"meetico/internal/repo"
	"meetico/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to create db pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	sqlDB := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer sqlDB.Close()

	if err := runMigrations(sqlDB); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	usersRepo := repo.NewUsersRepo(pool)
	schedulesRepo := repo.NewSchedulesRepo(pool)
	eventTypesRepo := repo.NewEventTypesRepo(pool)
	bookingsRepo := repo.NewBookingsRepo(pool)
	notificationsRepo := repo.NewNotificationsRepo(pool)

	slotsService := service.NewSlotsService(schedulesRepo, bookingsRepo)
	mailer := notify.NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom, cfg.BaseURL)
	telegramNotifier := notify.NewBot(cfg.TelegramBotToken)
	bookingService := service.NewBookingService(bookingsRepo, eventTypesRepo, usersRepo, notificationsRepo, schedulesRepo, slotsService, mailer, telegramNotifier)

	authHandler := handler.NewAuthHandler(usersRepo, schedulesRepo, cfg.JWTSecret)
	scheduleHandler := handler.NewScheduleHandler(schedulesRepo)
	eventTypesHandler := handler.NewEventTypesHandler(eventTypesRepo)
	bookingsHandler := handler.NewBookingsHandler(bookingsRepo, bookingService)
	publicHandler := handler.NewPublicHandler(usersRepo, eventTypesRepo, bookingsRepo, schedulesRepo, slotsService, bookingService)

	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware()

	r := chi.NewRouter()
	r.Use(middleware.CORS)
	r.Use(middleware.Compress)

	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Middleware)
		r.Get("/api/me", authHandler.GetMe)
		r.Patch("/api/me", authHandler.UpdateMe)
		r.Post("/api/auth/change-password", authHandler.ChangePassword)

		r.Get("/api/schedule", scheduleHandler.GetSchedule)
		r.Patch("/api/schedule", scheduleHandler.UpdateSchedule)

		r.Get("/api/event-types", eventTypesHandler.GetEventTypes)
		r.Post("/api/event-types", eventTypesHandler.CreateEventType)
		r.Patch("/api/event-types/{id}", eventTypesHandler.UpdateEventType)
		r.Delete("/api/event-types/{id}", eventTypesHandler.DeleteEventType)

		r.Get("/api/bookings", bookingsHandler.GetBookings)
		r.Delete("/api/bookings/{id}", bookingsHandler.CancelBooking)
	})

	r.Group(func(r chi.Router) {
		r.Use(rateLimitMiddleware.Middleware)
		r.Get("/api/public/{username}/{slug}", publicHandler.GetPublicInfo)
		r.Get("/api/public/{username}/{slug}/slots", publicHandler.GetPublicSlots)
		r.Post("/api/public/{username}/{slug}/book", publicHandler.CreatePublicBooking)
	})

	r.Get("/api/public/cancel/{token}", publicHandler.GetCancelInfo)
	r.Post("/api/public/cancel/{token}", publicHandler.CancelBooking)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		slog.Info("server started", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown error", "error", err)
	}

	slog.Info("server stopped")
}

func runMigrations(db *sql.DB) error {
	return goose.Up(db, "./migrations")
}
