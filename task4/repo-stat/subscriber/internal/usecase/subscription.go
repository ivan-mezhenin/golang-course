package usecase

import (
	"context"
	"fmt"

	"repo-stat/subscriber/internal/domain"
)

type SubscriptionUseCase struct {
	repo SubscriptionRepository
}

func NewSubscriptionUseCase(repo SubscriptionRepository) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		repo: repo,
	}
}

func (uc *SubscriptionUseCase) Create(ctx context.Context, owner, repo string) error {
	if owner == "" || repo == "" {
		return fmt.Errorf("owner and repo cannot be empty")
	}

	sub := domain.NewSubscription(owner, repo)

	exists, err := uc.repo.Exists(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to check subscription existence: %w", err)
	}
	if exists {
		return fmt.Errorf("subscription for %s/%s already exists", owner, repo)
	}

	if err := uc.repo.Create(ctx, sub); err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

func (uc *SubscriptionUseCase) Delete(ctx context.Context, owner, repo string) error {
	if owner == "" || repo == "" {
		return fmt.Errorf("owner and repo cannot be empty")
	}

	return uc.repo.Delete(ctx, owner, repo)
}

func (uc *SubscriptionUseCase) List(ctx context.Context) ([]*domain.Subscription, error) {
	return uc.repo.List(ctx)
}
