package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"repo-stat/collector/config"
	"repo-stat/collector/internal/adapter/github"
	"repo-stat/collector/internal/adapter/kafka"
	"repo-stat/collector/internal/adapter/subscriber"
	"repo-stat/collector/internal/domain"
	"repo-stat/platform/logger"
)

func run(ctx context.Context) error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "config file")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	log := logger.MustMakeLogger(cfg.Logger.LogLevel)

	log.Info("starting collector server...")

	// GitHub adapter
	ghAdapter := github.NewAdapter()

	// Kafka Response Producer
	responseProducer := kafka.NewResponseProducer([]string{cfg.Services.Kafka})

	// Kafka Task Consumer
	taskConsumer := kafka.NewTaskConsumer(
		[]string{cfg.Services.Kafka},
		"collector-task-group",
		ghAdapter,
		responseProducer,
		log,
	)

	go startSubscriptionUpdater(ctx, cfg, log, responseProducer)

	taskConsumer.Start(ctx)

	return nil
}

func startSubscriptionUpdater(ctx context.Context, cfg config.Config, log *slog.Logger, producer *kafka.ResponseProducer) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	subClient, err := subscriber.NewClient(cfg.Services.Subscriber, log)
	if err != nil {
		log.Error("failed to create subscriber client for updater", "error", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Debug("starting scheduled subscription update")

			subs, err := subClient.GetSubscriptionsInfo(ctx)
			if err != nil {
				log.Error("failed to get subscriptions from subscriber", "error", err)
				continue
			}

			for _, sub := range subs {
				if err := producer.Publish(ctx, domain.RepoResponse{Owner: sub.Owner, Repo: sub.Repo}); err != nil {
					log.Error("failed to publish scheduled request", "owner", sub.Owner, "repo", sub.Repo, "error", err)
				}
			}

			log.Info("scheduled update completed", "subscriptions_count", len(subs))
		}
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "collector error: %v\n", err)
		os.Exit(1)
	}
}
