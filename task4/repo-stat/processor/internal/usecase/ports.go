package usecase

import (
	"context"
	"repo-stat/processor/internal/domain"
)

type RepoGetter interface {
	Get(ctx context.Context, owner, repo string) (*domain.Repository, error)
	GetSubscriptionsInfo(ctx context.Context) (*domain.SubscriptionInfo, error)
}
