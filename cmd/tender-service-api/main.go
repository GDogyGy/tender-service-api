package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"TenderServiceApi/internal/config"
	"TenderServiceApi/internal/handlers/tender"
	"TenderServiceApi/internal/repository"
	"TenderServiceApi/internal/service"
	"TenderServiceApi/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log, err := setupLogger(cfg.DebugLevel)
	if err != nil {
		log.Error("Failed to init logger", slog.StringValue(err.Error()))
		os.Exit(1)
	}

	log.Info("Starting tender api server", slog.String("Env", cfg.Env))

	storage, err := postgres.New(ctx, cfg.PostgresConn)
	defer storage.Close()

	if err != nil {
		log.Error("Failed to init storage", slog.StringValue(err.Error()))
		os.Exit(1)
	}

	router := http.NewServeMux()
	teR := repository.NewTenderRepository(storage.Db)
	teS := service.NewTenderService(teR)
	handler := tender.NewHandler(log, teS)

	handler.Register(router)

	StartServer(ctx, cfg, log, router)
}

func StartServer(ctx context.Context, cfg *config.Config, log *slog.Logger, router http.Handler) {
	log.Info("server starting", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("listen and serve returned err:", err)
		}
	}()

	<-ctx.Done()

	log.Info("got interruption signal")
	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Info("server shutdown returned an err: %v\n", err)
	}

	log.Info("final")
}

func setupLogger(lvl string) (*slog.Logger, error) {
	const op = "main.setupLogger"
	var sl slog.Level
	err := sl.UnmarshalText([]byte(lvl))
	if err != nil {
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})), fmt.Errorf("%s:%v", op, err)
	}
	return slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: sl}),
	), nil
}
