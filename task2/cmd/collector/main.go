package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"repo-stat/internal/collector/adapter/github"
	"repo-stat/internal/collector/delivery/grpcserver"
	"repo-stat/internal/usecase/repo"
	"repo-stat/proto"
)

func main() {
	ghAdapter := github.NewAdapter()

	usecase := repo.NewGetRepoInfo(ghAdapter)

	ourServer := grpcserver.NewServer(usecase)

	grpcServer := grpc.NewServer()

	proto.RegisterRepoServiceServer(grpcServer, ourServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Collector gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
