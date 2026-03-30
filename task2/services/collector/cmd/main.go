package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"repo-stat/pkg/pb"
	"repo-stat/services/collector/internal/adapter/github"
	"repo-stat/services/collector/internal/controller"
	"repo-stat/services/collector/internal/usecase"
)

func main() {
	ghAdapter := github.NewAdapter()

	usecase := usecase.NewGetRepoInfo(ghAdapter)

	handler := controller.NewHandler(usecase)

	grpcServer := grpc.NewServer()

	pb.RegisterRepoServiceServer(grpcServer, handler)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Collector gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
