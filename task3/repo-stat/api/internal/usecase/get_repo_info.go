package usecase

import (
	"context"
	"repo-stat/api/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, domain.ErrRepoNotFound
			case codes.InvalidArgument:
				return nil, domain.ErrInvalidInput
			case codes.ResourceExhausted:
				return nil, domain.ErrGitHubRateLimited
			default:
				return nil, domain.ErrGitHubAPIError
			}
		}
		return nil, domain.ErrGitHubAPIError
	}

	return repositoryInfo, nil
}
