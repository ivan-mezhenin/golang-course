package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"repo-stat/processor/internal/domain"
	"repo-stat/processor/internal/usecase"

	"github.com/segmentio/kafka-go"
)

type SubscriptionProducer struct {
	writer *kafka.Writer
	repo   usecase.Repository
	log    *slog.Logger
}

func NewSubscriptionProducer(brokers []string, repo usecase.Repository, log *slog.Logger) *SubscriptionProducer {
	return &SubscriptionProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "subscription-updates",
			Balancer: &kafka.LeastBytes{},
		},
		repo: repo,
		log:  log,
	}
}

func (p *SubscriptionProducer) Start(ctx context.Context, interval time.Duration) {
	p.log.Info("starting subscription producer", "interval", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.log.Info("stopping subscription producer")
			return
		case <-ticker.C:
			p.sendSubscriptions(ctx)
		}
	}
}

func (p *SubscriptionProducer) sendSubscriptions(ctx context.Context) {
	subs, err := p.repo.ListSubscriptions(ctx)
	if err != nil {
		p.log.Error("failed to list subscriptions", "error", err)
		return
	}

	if len(subs) == 0 {
		p.log.Debug("no subscriptions to send")
		return
	}

	subList := make([]domain.Subscription, 0, len(subs))
	for _, sub := range subs {
		subList = append(subList, domain.Subscription{
			Owner: sub.Owner,
			Repo:  sub.Repo,
		})
	}

	data, err := json.Marshal(subList)
	if err != nil {
		p.log.Error("failed to marshal subscriptions", "error", err)
		return
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("subscriptions"),
		Value: data,
	})
	if err != nil {
		p.log.Error("failed to write subscriptions to kafka", "error", err)
		return
	}

	p.log.Info("sent subscriptions to kafka", "count", len(subs))
}

func (p *SubscriptionProducer) Close() error {
	return p.writer.Close()
}
