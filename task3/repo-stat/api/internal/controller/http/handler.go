package http

import (
	"context"
	"log/slog"
	"net/http"

	"repo-stat/api/config"
	"repo-stat/api/internal/adapter/processor"
	"repo-stat/api/internal/adapter/subscriber"
	"repo-stat/api/internal/usecase"

	processorProto "repo-stat/proto/processor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewHandler(ctx context.Context, log *slog.Logger, cfg config.Config) (http.Handler, error) {
	processorConn, err := grpc.NewClient(
		cfg.Services.Processor,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("cannot connect to processor", "error", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			if err := processorConn.Close(); err != nil {
				log.Error("failed to close processor connection", "error", err)
			}
		}
	}()

	subscriberClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("cannot init subscriber adapter", "error", err)
		if err := processorConn.Close(); err != nil {
			log.Error("failed to close processor connection", "error", err)
		}

		return nil, err
	}

	processorPB := processorProto.NewProcessorClient(processorConn)
	processorAdapter := processor.NewClient(processorPB)

	repoUseCase := usecase.NewGetRepoInfo(processorAdapter)
	pingUseCase := usecase.NewPing(subscriberClient, processorAdapter)

	mux := http.NewServeMux()
	AddRoutes(mux, log, pingUseCase, repoUseCase)

	log.Info("HTTP handlers initialized successfully")

	var handler http.Handler = mux
	return handler, nil
}
