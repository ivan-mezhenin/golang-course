package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"repo-stat/subscriber/internal/domain"

	"github.com/segmentio/kafka-go"
)

type SubscriptionProducer struct {
	writer *kafka.Writer
	log    *slog.Logger
}

func NewSubscriptionProducer(brokers []string, log *slog.Logger) *SubscriptionProducer {
	return &SubscriptionProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "subscription-updates",
			Balancer: &kafka.LeastBytes{},
		},
		log: log,
	}
}

func (p *SubscriptionProducer) PublishSubscriptions(ctx context.Context, subs []*domain.Subscription) error {
	// Convert to the format expected by processor
	subList := make([]domain.Subscription, 0, len(subs))
	for _, sub := range subs {
		subList = append(subList, *sub)
	}

	data, err := json.Marshal(subList)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("subscriptions"),
		Value: data,
	})
}

func (p *SubscriptionProducer) Close() error {
	return p.writer.Close()
}
