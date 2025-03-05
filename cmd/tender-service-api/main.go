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
	bidsCreate "TenderServiceApi/internal/handlers/bids/create"
	bidsFetch "TenderServiceApi/internal/handlers/bids/fetch"
	bidsUpdate "TenderServiceApi/internal/handlers/bids/update"
	pingFetch "TenderServiceApi/internal/handlers/ping/fetch"
	swaggerFetch "TenderServiceApi/internal/handlers/swagger/fetch"
	tenderCreate "TenderServiceApi/internal/handlers/tender/create"
	tenderFetch "TenderServiceApi/internal/handlers/tender/fetch"
	tenderUpdate "TenderServiceApi/internal/handlers/tender/update"
	bidDecisionRepository "TenderServiceApi/internal/repository/bid_decision"
	bidFeedbackRepository "TenderServiceApi/internal/repository/bid_feedback"
	bidsRepository "TenderServiceApi/internal/repository/bids"
	organizationRepository "TenderServiceApi/internal/repository/organization"
	tenderRepository "TenderServiceApi/internal/repository/tender"
	"TenderServiceApi/internal/storage/postgres"
	bidFeedbackUseCaseCreate "TenderServiceApi/internal/usecases/bid_feedback/create"
	bidFeedbackUseCaseFetch "TenderServiceApi/internal/usecases/bid_feedback/fetch"
	bidFeedbackUseCaseVerification "TenderServiceApi/internal/usecases/bid_feedback/verification"
	bidsUseCaseCreate "TenderServiceApi/internal/usecases/bids/create"
	bidsUseCaseDecision "TenderServiceApi/internal/usecases/bids/decision"
	bidsUseCaseEdit "TenderServiceApi/internal/usecases/bids/edit"
	bidsUseCaseFetch "TenderServiceApi/internal/usecases/bids/fetch"
	bidsUseCaseVerify "TenderServiceApi/internal/usecases/bids/verification"
	organizationUseCaseVerify "TenderServiceApi/internal/usecases/organization/verification"
	tenderUseCaseCreate "TenderServiceApi/internal/usecases/tender/create"
	tenderUseCaseEdit "TenderServiceApi/internal/usecases/tender/edit"
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

	// <! Repository
	repositoryTender := tenderRepository.NewRepository(storage.Db)
	repositoryBids := bidsRepository.NewRepository(storage.Db)
	repositoryOrganization := organizationRepository.NewRepository(storage.Db)
	repositoryBidFeedback := bidFeedbackRepository.NewRepository(storage.Db)
	repositoryBidDecision := bidDecisionRepository.NewRepository(storage.Db)
	// Repository !>

	useCaseOrganizationVerify := organizationUseCaseVerify.NewService(repositoryOrganization)

	// <! useCase Tender
	useCaseTenderVerify := tenderUseCaseVerify.NewService(repositoryTender)
	useCaseTenderFetch := tenderUseCaseFetch.NewService(repositoryTender, useCaseTenderVerify)
	useCaseTenderEdit := tenderUseCaseEdit.NewService(repositoryTender, useCaseTenderVerify)
	UseCaseTenderCreate := tenderUseCaseCreate.NewService(repositoryTender, useCaseOrganizationVerify)
	// useCase Tender !>

	// <! useCase Bids
	useCaseBidsCreate := bidsUseCaseCreate.NewService(repositoryBids, useCaseOrganizationVerify)
	useCaseBidsVerify := bidsUseCaseVerify.NewService(repositoryBids)
	useCaseBidsFetch := bidsUseCaseFetch.NewService(repositoryBids, repositoryBidFeedback, useCaseTenderVerify, useCaseBidsVerify)
	useCaseBidsEdit := bidsUseCaseEdit.NewService(repositoryBids, useCaseBidsVerify)
	useCaseBidsDecision := bidsUseCaseDecision.NewService(repositoryBids, repositoryBidDecision, repositoryTender, useCaseOrganizationVerify, useCaseTenderVerify)
	// useCase Bids !>

	// <! useCase Bids_Feedback
	useCaseBidFeedbackVerification := bidFeedbackUseCaseVerification.NewService(repositoryBidFeedback)
	useCaseBidFeedbackCreate := bidFeedbackUseCaseCreate.NewService(repositoryBidFeedback, repositoryOrganization, useCaseBidFeedbackVerification)
	useCaseBidFeedbackFetch := bidFeedbackUseCaseFetch.NewService(repositoryBidFeedback, useCaseTenderVerify, useCaseOrganizationVerify)
	// useCase Bids_Feedback !>

	// <! Handler Tender
	handlerTenderFetch := tenderFetch.NewHandler(log, useCaseTenderFetch)
	handlerTenderCreate := tenderCreate.NewHandler(log, UseCaseTenderCreate)
	handlerTenderUpdate := tenderUpdate.NewHandler(log, useCaseTenderEdit)
	// Handler Tender !>

	// <! Handler Bids
	handlerBidsCreate := bidsCreate.NewHandler(log, useCaseBidsCreate, useCaseBidFeedbackCreate)
	handlerBidsFetch := bidsFetch.NewHandler(log, useCaseBidsFetch, useCaseBidFeedbackFetch)
	handlerBidsUpdate := bidsUpdate.NewHandler(log, useCaseBidsEdit, useCaseBidsDecision)
	// Handler Bids !>

	// <! Handler Ping
	handlerPingFetch := pingFetch.NewHandler()
	// Handler Ping !>

	// <! Handler Swagger
	handlerSwaggerFetch := swaggerFetch.NewHandler()
	// Handler Swagger !>

	handlerTenderFetch.Register(router)
	handlerTenderCreate.Register(router)
	handlerTenderUpdate.Register(router)

	handlerBidsCreate.Register(router)
	handlerBidsFetch.Register(router)
	handlerBidsUpdate.Register(router)
	handlerPingFetch.Register(router)
	handlerSwaggerFetch.Register(router)

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
