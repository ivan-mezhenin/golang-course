package redis

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host string `yaml:"host" env:"REDIS_HOST" env-default:"localhost"`
	Port string `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
}

type Client struct {
	cli *redis.Client
	log *slog.Logger
}

func New(cfg Config, log *slog.Logger) *Client {
	addr := cfg.Host + ":" + cfg.Port
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	return &Client{
		cli: rdb,
		log: log,
	}
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	val, err := c.cli.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (c *Client) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.cli.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Ping(ctx context.Context) error {
	return c.cli.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.cli.Close()
}
