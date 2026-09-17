package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"meetico/internal/config"
	"meetico/internal/db"
	"meetico/internal/notify"
	"meetico/internal/worker"
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

	mailer := notify.NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPFrom, cfg.BaseURL)
	telegramNotifier := notify.NewBot(cfg.TelegramBotToken)

	go worker.Start(ctx, pool, mailer, telegramNotifier)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down worker...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("worker stopped")
}
