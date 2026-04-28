package usecase

import (
	"context"
	"fmt"

	"repo-stat/processor/internal/domain"
)

type GetRepoUseCase struct {
	repo     Repository
	producer MessageProducer
}

func NewGetRepoUseCase(repo Repository, producer MessageProducer) *GetRepoUseCase {
	return &GetRepoUseCase{
		repo:     repo,
		producer: producer,
	}
}

func (uc *GetRepoUseCase) Get(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if owner == "" || repo == "" {
		return nil, domain.ErrInvalidInput
	}

	cached, err := uc.repo.GetRepoFromCache(ctx, owner, repo)
	if err == nil && cached != nil {
		return cached, nil
	}

	if err := uc.producer.PublishRepoRequest(ctx, owner, repo); err != nil {
		return nil, fmt.Errorf("failed to publish request to kafka: %w", err)
	}

	return &domain.Repository{
		Owner: owner,
		Repo:  repo,
	}, nil
}

func (uc *GetRepoUseCase) GetSubscriptionsInfo(ctx context.Context) (*domain.SubscriptionInfo, error) {
	subs, err := uc.repo.ListSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	var repos []domain.Repository

	for _, sub := range subs {
		cached, err := uc.repo.GetRepoFromCache(ctx, sub.Owner, sub.Repo)
		if err == nil && cached != nil {
			repos = append(repos, *cached)
			continue
		}

		_ = uc.producer.PublishRepoRequest(ctx, sub.Owner, sub.Repo)

		repos = append(repos, domain.Repository{
			Owner: sub.Owner,
			Repo:  sub.Repo,
		})
	}

	return &domain.SubscriptionInfo{Repositories: repos}, nil
}
