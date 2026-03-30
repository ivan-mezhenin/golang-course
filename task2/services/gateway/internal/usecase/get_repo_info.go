package usecase

import (
	"context"
	"repo-stat/services/gateway/internal/domain"
)

type RepoGetter interface {
	Get(ctx context.Context, owner, repo string) (*domain.Repository, error)
}

type GetRepoInfo struct {
	repo RepoGetter
}

func NewGetRepoInfo(repo RepoGetter) *GetRepoInfo {
	return &GetRepoInfo{
		repo: repo,
	}
}

func (gri *GetRepoInfo) Execute(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if owner == "" || repo == "" {
		return nil, domain.ErrInvalidInput
	}

	repositoryInfo, err := gri.repo.Get(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	return repositoryInfo, nil
}
