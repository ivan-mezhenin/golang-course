package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"repo-stat/collector/config"
	"repo-stat/collector/internal/adapter/github"
	"repo-stat/collector/internal/controller"
	"repo-stat/collector/internal/usecase"
	"repo-stat/platform/logger"
	proto "repo-stat/proto/collector"
)

func run() error {

	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "server configuration file")
	flag.Parse()

	cfg := config.MustLoad(configPath)

	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting server...")
	log.Debug("debug messages are enabled")

	ghAdapter := github.NewAdapter()

	usecase := usecase.NewGetRepoInfo(ghAdapter)

	handler := controller.NewHandler(usecase)

	grpcServer := grpc.NewServer()

	proto.RegisterCollectorServer(grpcServer, handler)

	lis, err := net.Listen("tcp", cfg.GRPC.Address)
	if err != nil {
		log.Error("failed to listen: ", "error", err)
		return err
	}

	log.Info("Collector gRPC server started", "address", cfg.GRPC.Address)
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
