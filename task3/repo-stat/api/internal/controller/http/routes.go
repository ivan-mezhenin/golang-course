package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, pingUC *usecase.Ping, repoUC *usecase.GetRepoInfo) {
	mux.Handle("GET /api/ping", NewPingHandler(log, pingUC))

	mux.Handle("GET /api/repositories/info", NewRepoHandler(log, repoUC))
}
