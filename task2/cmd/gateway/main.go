package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "repo-stat/docs"
	"repo-stat/internal/gateway/controller/grpc"
	httpHandler "repo-stat/internal/gateway/controller/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {
	collectorAddr := os.Getenv("COLLECTOR_ADDR")
	if collectorAddr == "" {
		collectorAddr = "localhost:50051"
	}
	client, err := grpc.NewClient(collectorAddr)
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}
	defer client.Close()

	handler := httpHandler.NewHandler(client)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /repos/{owner}/{repo}", handler.GetRepoInfo)
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
	))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("API Gateway started on http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
