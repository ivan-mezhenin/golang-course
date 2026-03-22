package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"repo-stat/proto"
)

type Client struct {
	service proto.RepoServiceClient
	conn    *grpc.ClientConn
}

func NewClient(collectorAddr string) (*Client, error) {
	conn, err := grpc.Dial(
		collectorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial collector: %w", err)
	}

	return &Client{
		service: proto.NewRepoServiceClient(conn),
		conn:    conn,
	}, nil
}

func (c *Client) GetRepoInfo(ctx context.Context, owner, repo string) (*proto.GetRepoResponse, error) {
	req := &proto.GetRepoRequest{
		Owner: owner,
		Repo:  repo,
	}
	return c.service.GetRepoInfo(ctx, req)
}

func (c *Client) Close() error {
	return c.conn.Close()
}
