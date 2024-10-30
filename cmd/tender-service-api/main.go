package main

import (
	"TenderServiceApi/internal/config"
	"TenderServiceApi/internal/storage/postgres"
	"TenderServiceApi/internal/tender"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

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

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func StartServer(cfg *config.Config, log *slog.Logger, router http.Handler) {
	log.Info("server starting", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
		log.Error("Failed to start server")
	}
}
