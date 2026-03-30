package processor

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"repo-stat/api/internal/domain"
	processorProto "repo-stat/proto/processor"
)

type Client struct {
	pp processorProto.ProcessorClient
}

func NewClient(pp processorProto.ProcessorClient) *Client {
	return &Client{
		pp: pp,
	}
}

func (c *Client) Get(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	resp, err := c.pp.GetRepository(ctx, &processorProto.GetRepoRequest{
		Owner: owner,
		Repo:  repo,
	})
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
		return nil, fmt.Errorf("processor client: %w", err)
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

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pp.Ping(ctx, &processorProto.PingRequest{})
	if err != nil {
		return domain.PingStatusDown
	}
	return domain.PingStatusUp
}
