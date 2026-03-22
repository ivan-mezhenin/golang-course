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
	// 1. Создаём адаптер (GitHub клиент)
	ghAdapter := github.NewAdapter()

	usecase := repo.NewGetRepoInfo(ghAdapter)

	// 2. Создаём use-case, передаём ему адаптер (порт → адаптер)
	ourServer := grpcserver.NewServer(usecase) // ← вызываем конструктор из ТВОЕГО пакета

	// 4. Создаём экземпляр gRPC-сервера из библиотеки
	grpcServer := grpc.NewServer()

	// 5. Регистрируем наш сервер как обработчик сервиса RepoService
	proto.RegisterRepoServiceServer(grpcServer, ourServer)

	// 5. Запускаем прослушивание порта
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Println("Collector gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
