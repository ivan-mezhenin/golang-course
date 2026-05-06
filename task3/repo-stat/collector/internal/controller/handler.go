package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"repo-stat/collector/internal/domain"
	repo "repo-stat/collector/internal/usecase"
	proto "repo-stat/proto/collector"
)

type Handler struct {
	proto.UnimplementedCollectorServer

	usecase *repo.GetRepoInfo
}

func NewHandler(usecase *repo.GetRepoInfo) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) GetRepoInfo(ctx context.Context, req *proto.GetRepoRequest) (*proto.GetRepoResponse, error) {
	repoData, err := h.usecase.Execute(ctx, req.Owner, req.Repo)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case domain.ErrRepoNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case domain.ErrGitHubRateLimited:
			return nil, status.Error(codes.ResourceExhausted, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error: "+err.Error())
		}
	}

	return &proto.GetRepoResponse{
		Name:        repoData.Name,
		Description: repoData.Description,
		Stars:       int32(repoData.Stars),
		Forks:       int32(repoData.Forks),
		CreatedAt:   repoData.CreatedAt,
		Visibility:  repoData.Visibility,
	}, nil
}
