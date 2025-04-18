package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"TenderServiceApi/internal/config"
	bidsCreateGrpc "TenderServiceApi/internal/handlers/grpc/bids/create"
	bidsFetchGrpc "TenderServiceApi/internal/handlers/grpc/bids/fetch"
	bidsUpdateGrpc "TenderServiceApi/internal/handlers/grpc/bids/update"
	tenderCreateGrpc "TenderServiceApi/internal/handlers/grpc/tender/create"
	tenderFetchGrpc "TenderServiceApi/internal/handlers/grpc/tender/fetch"
	tenderUpdateGrpc "TenderServiceApi/internal/handlers/grpc/tender/update"
	bidsCreate "TenderServiceApi/internal/handlers/rest/bids/create"
	bidsFetch "TenderServiceApi/internal/handlers/rest/bids/fetch"
	bidsUpdate "TenderServiceApi/internal/handlers/rest/bids/update"
	pingFetch "TenderServiceApi/internal/handlers/rest/ping/fetch"
	swaggerFetch "TenderServiceApi/internal/handlers/rest/swagger/fetch"
	tenderCreate "TenderServiceApi/internal/handlers/rest/tender/create"
	tenderFetch "TenderServiceApi/internal/handlers/rest/tender/fetch"
	tenderUpdate "TenderServiceApi/internal/handlers/rest/tender/update"
	"TenderServiceApi/internal/kafka"
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
	"google.golang.org/grpc"
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

	log.Info("Starting tender api server", slog.String("Env", cfg.Env))

	storage, err := postgres.New(ctx, cfg.PostgresConn)
	if err != nil {
		log.Error("Failed to init storage", slog.Attr{Value: slog.StringValue(err.Error())})
		os.Exit(1) // nolint:gocritic
	}

	defer storage.Close()
	// <! kafka init
	producer, err := kafka.NewEventProducer(cfg.Kafka.GetAddresses(), cfg.Kafka.DefaultTopic)
	if err != nil {
		log.Error("Failed to init kafka", slog.Attr{Value: slog.StringValue(err.Error())})
		os.Exit(1) // nolint:gocritic
	}
	defer func() { _ = producer.Close() }()
	// kafka init !>
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
	handlerBidsCreate := bidsCreate.NewHandler(log, producer, useCaseBidsCreate, useCaseBidFeedbackCreate)
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

	gRPCServer := grpc.NewServer()

	// <! gRPC Handler Tender
	tenderCreateGrpc.NewHandler(gRPCServer, log, UseCaseTenderCreate)
	tenderFetchGrpc.NewHandler(gRPCServer, log, useCaseTenderFetch)
	tenderUpdateGrpc.NewHandler(gRPCServer, log, useCaseTenderEdit)
	// gRPC Handler Tender !>

	// <! gRPC Handler Bids
	bidsCreateGrpc.NewHandler(gRPCServer, log, useCaseBidsCreate, useCaseBidFeedbackCreate)
	bidsFetchGrpc.NewHandler(gRPCServer, log, useCaseBidsFetch, useCaseBidFeedbackFetch)
	bidsUpdateGrpc.NewHandler(gRPCServer, log, useCaseBidsEdit, useCaseBidsDecision)
	// gRPC Handler Tender !>

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		StartServerHttp(ctx, cfg, log, router)
	}()

	go func() {
		defer wg.Done()
		StartServerGrpc(ctx, cfg, log, gRPCServer)
	}()

	wg.Wait()
	log.Info("All servers stopped")
}

func StartServerHttp(ctx context.Context, cfg *config.Config, log *slog.Logger, router http.Handler) {
	log.Info("http server starting", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http listen and serve returned err:", slog.Attr{Value: slog.StringValue(err.Error())})
		}
	}()

	<-ctx.Done()

	if err := srv.Shutdown(ctx); err != nil {
		log.Info("server shutdown returned an err: %v\n", slog.Attr{Value: slog.StringValue(err.Error())})
	}
	log.Info("http server stopping")
}

func StartServerGrpc(ctx context.Context, cfg *config.Config, log *slog.Logger, gRPCServer *grpc.Server) {
	log.Info("grpc server starting", slog.String("address", cfg.GRPCServer.Address))

	_ = cfg

	lis, err := net.Listen("tcp", cfg.GRPCServer.Address)
	if err != nil {
		log.Error("gRPC listen error:", slog.String("error", err.Error()))
		return
	}

	go func() {
		if err := gRPCServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Error("gRPC serve error:", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()

	gRPCServer.GracefulStop()
	log.Info("grpc server stopping")
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
