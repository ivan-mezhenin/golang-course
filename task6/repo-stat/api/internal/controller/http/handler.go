package http

import (
	"context"
	"log/slog"
	"net/http"

	"repo-stat/api/config"
	"repo-stat/api/internal/adapter/processor"
	"repo-stat/api/internal/adapter/subscriber"
	"repo-stat/api/internal/usecase"
	"repo-stat/platform/redis"
)

func NewHandler(ctx context.Context, log *slog.Logger, cfg config.Config) (http.Handler, error) {

	subscriberClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("cannot init subscriber adapter", "error", err)
		return nil, err
	}

	processorAdapter, err := processor.NewClient(cfg.Services.Processor, log)
	if err != nil {
		log.Error("cannot init processor adapter", "error", err)
		return nil, err
	}

	repoUseCase := usecase.NewGetRepoInfo(processorAdapter)
	pingUseCase := usecase.NewPing(subscriberClient, processorAdapter)
	subscriptionUseCase := usecase.NewSubscriptionUseCase(subscriberClient)

	mux := http.NewServeMux()
	AddRoutes(mux, log, pingUseCase, repoUseCase, subscriptionUseCase)

	log.Info("HTTP handlers initialized successfully")

	var handler http.Handler = mux

	redisClient := redis.New(cfg.Redis, log)
	if err := redisClient.Ping(ctx); err != nil {
		log.Warn("redis is not available, using in-memory fallback", "error", err)
		redisClient = nil
	} else {
		log.Info("redis connected successfully")
	}

	limiter := newInMemoryLimiter(cfg.RateLimit.RequestsPerSecond, cfg.RateLimit.Burst)

	handler = CacheMiddleware(log, redisClient, cfg.Cache.TTL(), handler)
	handler = RateLimitMiddleware(log, limiter, handler)

	return handler, nil
}
