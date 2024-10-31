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
	"TenderServiceApi/internal/storage/postgres"
	"TenderServiceApi/internal/tender"
)

func main() {
	cfg := config.MustLoad()

	log, err := setupLogger(cfg.DebugLevel)
	if err != nil {
		log.Error("Failed to init logger", slog.StringValue(err.Error()))
		os.Exit(1)
	}

	log.Info("Starting tender api server", slog.String("Env", cfg.Env))

	storage, err := postgres.New(cfg.PostgresConn)

	if err != nil {
		log.Error("Failed to init storage", slog.StringValue(err.Error()))
		os.Exit(1)
	}
	router := http.NewServeMux()
	//_ = storage
	handler := tender.NewHandler(log, storage)
	handler.Register(router)

	StartServer(cfg, log, router)
}

func StartServer(cfg *config.Config, log *slog.Logger, router http.Handler) {
	log.Info("server starting", slog.String("address", cfg.Address))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

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
