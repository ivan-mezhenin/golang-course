package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"repo-stat/pkg/pb"
	_ "repo-stat/services/gateway/docs"
	"repo-stat/services/gateway/internal/adapter/collector"
	httpHandler "repo-stat/services/gateway/internal/controller/httpController"
	"repo-stat/services/gateway/internal/usecase"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const gatewayPort = ":8080"
const collectorAddress = "localhost:50051"

func main() {
	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = collectorAddress
	}

	conn, err := grpc.Dial(
		collectorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to collector: %v", err)
	}

	defer conn.Close()

	pbClient := pb.NewRepoServiceClient(conn)
	collClient, err := collector.NewClient(pbClient)
	if err != nil {
		log.Fatalf("failed to create collector adapter: %v", err)
	}
	repoUsecase := usecase.NewGetRepoInfo(collClient)
	handler := httpHandler.NewHandler(repoUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /repos/{owner}/{repo}", handler.Get)
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
	))

	srv := &http.Server{
		Addr:    gatewayPort,
		Handler: mux,
	}

	fmt.Println("API Gateway started on http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
