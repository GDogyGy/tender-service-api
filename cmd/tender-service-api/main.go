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
	tenderCreate "TenderServiceApi/internal/handlers/tender/create"
	tenderFetch "TenderServiceApi/internal/handlers/tender/fetch"
	tenderUpdate "TenderServiceApi/internal/handlers/tender/update"
	"TenderServiceApi/internal/repository/organization"
	tenderRepository "TenderServiceApi/internal/repository/tender"
	"TenderServiceApi/internal/storage/postgres"
	organizationUseCaseFetch "TenderServiceApi/internal/usecases/organization/fetch"
	organizationUseCaseVerify "TenderServiceApi/internal/usecases/organization/verification"
	tenderUseCaseCreate "TenderServiceApi/internal/usecases/tender/create"
	tenderUseCaseEdite "TenderServiceApi/internal/usecases/tender/edite"
	tenderUseCaseFetch "TenderServiceApi/internal/usecases/tender/fetch"
	tenderUseCaseVerify "TenderServiceApi/internal/usecases/tender/verification"
)

func main() {
	cfg := config.MustLoad()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// TODO: Тут не удалось победить линтер пришлось через nolint:gocritic решать ошибку
	log, err := setupLogger(cfg.DebugLevel)
	if err != nil {
		log.Error("Failed to init logger", slog.Attr{Value: slog.StringValue(err.Error())})
		os.Exit(1) // nolint:gocritic
	}

	log.Info("Starting organizationResponsible api server", slog.String("Env", cfg.Env))

	storage, err := postgres.New(ctx, cfg.PostgresConn)
	if err != nil {
		log.Error("Failed to init storage", slog.Attr{Value: slog.StringValue(err.Error())})
		os.Exit(1) // nolint:gocritic
	}

	defer storage.Close()

	router := http.NewServeMux()

	tenderRepository := tenderRepository.NewRepository(storage.Db)
	organizationRepository := organization.NewRepository(storage.Db)

	useCaseTenderVerify := tenderUseCaseVerify.NewService(tenderRepository)
	useCaseOrganizationVerify := organizationUseCaseVerify.NewService(organizationRepository)
	useCaseOrganizationFetch := organizationUseCaseFetch.NewService(organizationRepository)
	useCaseTenderFetch := tenderUseCaseFetch.NewService(tenderRepository, useCaseTenderVerify)
	useCaseTenderEdite := tenderUseCaseEdite.NewService(tenderRepository, useCaseOrganizationVerify, useCaseOrganizationFetch)
	tenderUseCaseCreate := tenderUseCaseCreate.NewService(tenderRepository, useCaseOrganizationVerify)

	handlerFetch := tenderFetch.NewHandler(log, useCaseTenderFetch)
	handlerCreate := tenderCreate.NewHandler(log, tenderUseCaseCreate)
	handlerUpdate := tenderUpdate.NewHandler(log, useCaseTenderEdite)
	handlerFetch.Register(router)
	handlerCreate.Register(router)
	handlerUpdate.Register(router)

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
