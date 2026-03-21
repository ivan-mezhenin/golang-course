package repo

import (
	"context"
	"repo-stat/internal/domain"
)

type RepoGetter interface {
	Get(ctx context.Context, owner, repo string) (*domain.Repository, error)
}
