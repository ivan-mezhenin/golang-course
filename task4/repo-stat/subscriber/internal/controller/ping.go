package controller

import (
	"context"
	"log/slog"
	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/internal/usecase"
)

type PingHandler struct {
	subscriberpb.UnimplementedSubscriberServer
	log  *slog.Logger
	ping *usecase.Ping
}

func NewPingHandler(log *slog.Logger, ping *usecase.Ping) *PingHandler {
	return &PingHandler{
		log:  log,
		ping: ping,
	}
}

func (ph *PingHandler) Ping(ctx context.Context, _ *subscriberpb.PingRequest) (*subscriberpb.PingResponse, error) {
	ph.log.Debug("subscriberp ping request received")

	return &subscriberpb.PingResponse{
		Reply: ph.ping.Execute(ctx),
	}, nil
}
