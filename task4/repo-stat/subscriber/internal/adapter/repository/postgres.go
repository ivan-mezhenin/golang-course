package repository

import (
	"context"
	"database/sql"
	"fmt"

	"repo-stat/subscriber/internal/adapter/repository/sqlc"
	"repo-stat/subscriber/internal/domain"
	"repo-stat/subscriber/internal/usecase"
)

type postgresRepository struct {
	queries *sqlc.Queries
}

func NewPostgresRepository(db *sql.DB) usecase.SubscriptionRepository {
	return &postgresRepository{
		queries: sqlc.New(db),
	}
}

func (r *postgresRepository) Create(ctx context.Context, sub *domain.Subscription) error {
	_, err := r.queries.CreateSubscription(ctx, sqlc.CreateSubscriptionParams{
		Owner: sub.Owner,
		Repo:  sub.Repo,
	})
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	return nil
}

func (r *postgresRepository) Delete(ctx context.Context, owner, repo string) error {
	err := r.queries.DeleteSubscription(ctx, sqlc.DeleteSubscriptionParams{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}

func (r *postgresRepository) List(ctx context.Context) ([]*domain.Subscription, error) {
	items, err := r.queries.ListSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	subs := make([]*domain.Subscription, len(items))
	for i, item := range items {
		subs[i] = &domain.Subscription{
			Owner: item.Owner,
			Repo:  item.Repo,
		}
	}
	return subs, nil
}

func (r *postgresRepository) Exists(ctx context.Context, owner, repo string) (bool, error) {
	exists, err := r.queries.ExistsSubscription(ctx, sqlc.ExistsSubscriptionParams{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return exists, nil
}