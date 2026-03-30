package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"repo-stat/collector/internal/adapter/github"
	"repo-stat/collector/internal/controller"
	"repo-stat/collector/internal/usecase"
	proto "repo-stat/proto/collector"
)

const collectorAddress = ":50051"

func main() {
	ghAdapter := github.NewAdapter()

	usecase := usecase.NewGetRepoInfo(ghAdapter)

	handler := controller.NewHandler(usecase)

	grpcServer := grpc.NewServer()

	proto.RegisterCollectorServer(grpcServer, handler)

	lis, err := net.Listen("tcp", collectorAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Collector gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
