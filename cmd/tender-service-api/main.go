package main

import (
	"TenderServiceApi/internal/repository/facade/organizationResponsible"
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
	tenderRepository "TenderServiceApi/internal/repository/tender"
	"TenderServiceApi/internal/storage/postgres"
	tenderUseCase "TenderServiceApi/internal/usecase/tender"
)

func main() {
	cfg := config.MustLoad()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log, err := setupLogger(cfg.DebugLevel)
	if err != nil {
		log.Error("Failed to init logger", slog.Attr{Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}

	log.Info("Starting organizationResponsible api server", slog.String("Env", cfg.Env))

	storage, err := postgres.New(ctx, cfg.PostgresConn)
	defer storage.Close()

	if err != nil {
		log.Error("Failed to init storage", slog.Attr{Value: slog.StringValue(err.Error())})
		os.Exit(1)
	}

	router := http.NewServeMux()
	organizationResponsibleFacade := organizationResponsible.NewOrganizationResponsibleFacade(storage.Db)

	tenderRepository := tenderRepository.NewRepository(storage.Db)
	tenderService := tenderUseCase.NewService(tenderRepository)
	handler := tender.NewHandler(log, tenderService, organizationResponsibleFacade)

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
			log.Error("listen and serve returned err:", slog.Attr{Value: slog.StringValue(err.Error())})
		}
	}()

	<-ctx.Done()

	log.Info("got interruption signal")
	if err := srv.Shutdown(ctx); err != nil {
		log.Info("server shutdown returned an err: %v\n", slog.Attr{Value: slog.StringValue(err.Error())})
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
