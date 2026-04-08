package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type Pinger interface {
	Ping(ctx context.Context) domain.PingStatus
}
