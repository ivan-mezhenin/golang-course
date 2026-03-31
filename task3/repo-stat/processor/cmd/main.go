package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"repo-stat/platform/logger"
	"repo-stat/processor/config"
	"repo-stat/processor/internal/adapter/collector"
	"repo-stat/processor/internal/controller"
	"repo-stat/processor/internal/usecase"
	collectorClient "repo-stat/proto/collector"
	processorServer "repo-stat/proto/processor"
)

func run() error {

	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting server...")
	log.Debug("debug messages are enabled")

	conn, err := grpc.NewClient(
		cfg.Services.Collector,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to collector: ", "error", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Error("failed to close collector connection", "error", err)
		}
	}()

	collectorClient := collectorClient.NewCollectorClient(conn)

	collectorAdapter, err := collector.NewClient(collectorClient)
	if err != nil {
		log.Error("failed to create collector client: ", "error", err)
		return err
	}

	repoUsecase := usecase.NewGetRepoInfo(collectorAdapter)

	handler := controller.NewHandler(repoUsecase)

	grpcServer := grpc.NewServer()
	processorServer.RegisterProcessorServer(grpcServer, handler)

	lis, err := net.Listen("tcp", cfg.GRPC.Address)
	if err != nil {
		log.Error("failed to listen: ", "error", err)
		return err
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Error("failed to serve: ", "error", err)
		return err
	}

	return nil
}

func main() {
	_, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := run(); err != nil {
		_, err = fmt.Fprintln(os.Stderr, err)
		if err != nil {
			fmt.Printf("launching server error: %s\n", err)
		}
		os.Exit(1)
	}
}
