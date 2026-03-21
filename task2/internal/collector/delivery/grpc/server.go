package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"repo-stat/internal/domain"
	"repo-stat/internal/usecase/repo"
	"repo-stat/proto"
)

type Server struct {
	proto.UnimplementedRepoServiceServer

	usecase *repo.GetRepoInfo
}

func NewServer(usecase *repo.GetRepoInfo) *Server {
	return &Server{usecase: usecase}
}

func (s *Server) GetRepoInfo(ctx context.Context, req *proto.GetRepoRequest) (*proto.GetRepoResponse, error) {
	repoData, err := s.usecase.Execute(ctx, req.Owner, req.Repo)
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
		CreatedAt:   repoData.CreatedAt.Format(time.RFC3339),
		Visibility:  repoData.Visibility,
	}, nil
}
