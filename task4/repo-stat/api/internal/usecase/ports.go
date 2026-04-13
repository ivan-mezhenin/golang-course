package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type Pinger interface {
	Ping(ctx context.Context) domain.PingStatus
}

type SubscriptionRepository interface {
	Create(ctx context.Context, owner, repo string) error

	Delete(ctx context.Context, owner, repo string) error

	List(ctx context.Context) ([]*domain.Subscription, error)
}
