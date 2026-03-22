package main

import (
	"fmt"
	"log"
	"net/http"

	"repo-stat/internal/gateway/controller/grpc"
	httpHandler "repo-stat/internal/gateway/controller/http"
)

func main() {
	collectorAddr := "localhost:50051"
	client, err := grpc.NewClient(collectorAddr)
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}
	defer client.Close()

	handler := httpHandler.NewHandler(client)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/repos/{owner}/{repo}", handler.GetRepoInfo)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("API Gateway started on http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
