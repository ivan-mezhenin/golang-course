package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"repo-stat/processor/internal/adapter/collector"
	"repo-stat/processor/internal/controller"
	"repo-stat/processor/internal/usecase"
	collectorClient "repo-stat/proto/collector"
	processorServer "repo-stat/proto/processor"
)

const (
	processorAddress = ":50052"
	collectorAddress = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(
		collectorAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to collector: %v", err)
	}
	defer conn.Close()

	collectorClient := collectorClient.NewCollectorClient(conn)

	collectorAdapter, err := collector.NewClient(collectorClient)
	if err != nil {
		log.Fatalf("failed to create collector client: %v", err)
	}

	repoUsecase := usecase.NewGetRepoInfo(collectorAdapter)

	handler := controller.NewHandler(repoUsecase)

	grpcServer := grpc.NewServer()
	processorServer.RegisterProcessorServer(grpcServer, handler)

	lis, err := net.Listen("tcp", processorAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("Processor gRPC server started on %s\n", processorAddress)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
