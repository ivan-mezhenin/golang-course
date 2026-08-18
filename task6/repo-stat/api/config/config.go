package config

import (
	"repo-stat/platform/env"
	"repo-stat/platform/httpserver"
	"repo-stat/platform/logger"
	"repo-stat/platform/redis"
	"time"
)

type App struct {
	AppName string `yaml:"app_name" env:"APP_NAME" env-default:"repo-stat-api"`
}

type Services struct {
	Subscriber string `yaml:"subscriber" env:"SUBSCRIBER_ADDRESS" env-default:"localhost:8081"`
	Processor  string `yaml:"processor" env:"PROCESSOR_ADDRESS" env-default:"localhost:8083"`
}

type Cache struct {
	TTLSeconds int `yaml:"ttl_seconds" env:"CACHE_TTL_SECONDS" env-default:"60"`
}

type RateLimit struct {
	RequestsPerSecond int `yaml:"requests_per_second" env:"RATE_LIMIT_RPS" env-default:"5"`
	Burst             int `yaml:"burst" env:"RATE_LIMIT_BURST" env-default:"10"`
}

type Config struct {
	App       App               `yaml:"app"`
	Services  Services          `yaml:"services"`
	HTTP      httpserver.Config `yaml:"http"`
	Logger    logger.Config     `yaml:"logger"`
	Redis     redis.Config      `yaml:"redis"`
	Cache     Cache             `yaml:"cache"`
	RateLimit RateLimit         `yaml:"rate_limit"`
}

func MustLoad(path string) Config {
	var cfg Config
	env.MustLoad(path, &cfg)
	return cfg
}

func (c *Cache) TTL() time.Duration {
	return time.Duration(c.TTLSeconds) * time.Second
}
