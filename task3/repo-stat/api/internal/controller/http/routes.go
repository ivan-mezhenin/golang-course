package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"

	httpSwagger "github.com/swaggo/http-swagger"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, pingUC *usecase.Ping, repoUC *usecase.GetRepoInfo) {
	mux.Handle("GET /api/ping", NewPingHandler(log, pingUC))

	mux.Handle("GET /api/repositories/info", NewRepoHandler(log, repoUC))
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
}
