package collector

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"repo-stat/processor/internal/domain"
	proto "repo-stat/proto/collector"
)

type Client struct {
	grpc proto.CollectorClient
}

func NewClient(grpc proto.CollectorClient) (*Client, error) {
	return &Client{
		grpc: grpc,
	}, nil
}

func (c *Client) Get(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	req := &proto.GetRepoRequest{
		Owner: owner,
		Repo:  repo,
	}
	resp, err := c.grpc.GetRepoInfo(ctx, req)
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
				return nil, fmt.Errorf("%w: %s", domain.ErrGitHubAPIError, st.Message())
			}
		}

		return nil, fmt.Errorf("collector: %w", err)
	}

	return &domain.Repository{
		Name:        resp.Name,
		Description: resp.Description,
		Stars:       resp.Stars,
		Forks:       resp.Forks,
		CreatedAt:   resp.CreatedAt,
		Visibility:  resp.Visibility,
	}, nil
}
