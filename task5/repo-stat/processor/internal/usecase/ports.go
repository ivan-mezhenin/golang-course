package usecase

import (
	"context"
	"repo-stat/processor/internal/domain"
)

type RepoGetter interface {
	Get(ctx context.Context, owner, repo string) (*domain.Repository, error)
	GetSubscriptionsInfo(ctx context.Context) (*domain.SubscriptionInfo, error)
}

type Repository interface {
	ListSubscriptions(ctx context.Context) ([]*domain.Subscription, error)
	ReplaceAllSubscriptions(ctx context.Context, subs []*domain.Subscription) error

	GetRepoFromCache(ctx context.Context, owner, repo string) (*domain.Repository, error)
	UpsertRepoCache(ctx context.Context, repo *domain.Repository) error
}
