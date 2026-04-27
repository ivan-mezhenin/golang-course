package kafka

import (
	"context"
	"encoding/json"

	"repo-stat/processor/internal/domain"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "repo-requests",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) PublishRepoRequest(ctx context.Context, owner, repo string) error {
	message := domain.RepoRequest{
		Owner: owner,
		Repo:  repo,
	}

	value, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(owner + "/" + repo),
		Value: value,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
