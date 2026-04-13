package controller

import (
	"context"
	"log/slog"
	subscriberpb "repo-stat/proto/subscriber"
	"repo-stat/subscriber/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubscriptionHandler struct {
	subscriberpb.UnimplementedSubscriberServer
	log          *slog.Logger
	subscription *usecase.SubscriptionUseCase
}

func NewSubscriptionHandler(log *slog.Logger, subscriptionUseCase *usecase.SubscriptionUseCase) *SubscriptionHandler {
	return &SubscriptionHandler{
		log:          log,
		subscription: subscriptionUseCase,
	}
}

func (sh *SubscriptionHandler) Create(ctx context.Context, request *subscriberpb.PostSubscriptionRequest) (*subscriberpb.PostSubscriptionResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "subscription is required")
	}

	err := sh.subscription.Create(ctx, request.Subscription.Owner, request.Subscription.Repo)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &subscriberpb.PostSubscriptionResponse{}, nil
}

func (sh *SubscriptionHandler) List(ctx context.Context, request *subscriberpb.ListSubscriptionRequest) (*subscriberpb.ListSubscriptionResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "subscription is required")
	}

	subscriptions, err := sh.subscription.List(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := make([]*subscriberpb.Subscription, 0, len(subscriptions))
	for _, sub := range response {
		response = append(response, &subscriberpb.Subscription{
			Owner: sub.Owner,
			Repo:  sub.Repo,
		})
	}

	return &subscriberpb.ListSubscriptionResponse{Subscriptions: response}, nil
}

func (sh *SubscriptionHandler) Delete(ctx context.Context, request *subscriberpb.DeleteSubscriptionRequest) (*subscriberpb.DeleteSubscriptionResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "subscription is required")
	}

	err := sh.subscription.Delete(ctx, request.Subscription.Owner, request.Subscription.Repo)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &subscriberpb.DeleteSubscriptionResponse{}, nil
}
