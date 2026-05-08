package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"repo-stat/platform/redis"
)

type RateLimiter interface {
	Allow(ctx context.Context, ip string) bool
}

type inMemoryLimiter struct {
	limiters *sync.Map
	r        rate.Limit
	b        int
}

func newInMemoryLimiter(rps int, burst int) *inMemoryLimiter {
	return &inMemoryLimiter{
		limiters: &sync.Map{},
		r:        rate.Limit(rps),
		b:        burst,
	}
}

func (l *inMemoryLimiter) getLimiter(ip string) *rate.Limiter {
	v, exists := l.limiters.Load(ip)
	if exists {
		return v.(*rate.Limiter)
	}

	limiter := rate.NewLimiter(l.r, l.b)
	l.limiters.Store(ip, limiter)
	return limiter
}

func (l *inMemoryLimiter) Allow(ctx context.Context, ip string) bool {
	limiter := l.getLimiter(ip)
	return limiter.Allow()
}

func RateLimitMiddleware(log *slog.Logger, limiter RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		if !limiter.Allow(r.Context(), ip) {
			log.Warn("rate limit exceeded", "ip", ip)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "rate limit exceeded",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CacheMiddleware(log *slog.Logger, redisClient *redis.Client, cacheTTL time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || redisClient == nil {
			next.ServeHTTP(w, r)
			return
		}

		cacheKey := fmt.Sprintf("cache:%s:%s", r.URL.Path, r.URL.RawQuery)

		cachedData, err := redisClient.Get(r.Context(), cacheKey)
		if err == nil && cachedData != "" {
			log.Info("cache hit", "key", cacheKey)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(cachedData))
			return
		}

		log.Info("cache miss", "key", cacheKey)

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		if rw.status == http.StatusOK && len(rw.body) > 0 {
			if err := redisClient.Set(r.Context(), cacheKey, string(rw.body), cacheTTL); err != nil {
				log.Error("failed to save to cache", "error", err)
			}
		}
	})
}

type responseWriter struct {
	http.ResponseWriter
	body   []byte
	status int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = append(rw.body, b...)
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
