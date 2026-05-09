package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type IdempotencyStore interface {
	IsProcessed(ctx context.Context, paymentID string) (bool, error)
	MarkProcessed(ctx context.Context, paymentID string) error
}

type RedisIdempotencyStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisIdempotencyStore(addr string, password string, db int, ttl time.Duration) *RedisIdempotencyStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisIdempotencyStore{
		client: client,
		ttl:    ttl,
	}
}

func (s *RedisIdempotencyStore) key(paymentID string) string {
	return fmt.Sprintf("notification:payment:%s", paymentID)
}

func (s *RedisIdempotencyStore) IsProcessed(ctx context.Context, paymentID string) (bool, error) {
	key := s.key(paymentID)

	exists, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if exists > 0 {
		log.Printf("[Idempotency] Duplicate notification skipped payment_id=%s", paymentID)
		return true, nil
	}

	return false, nil
}

func (s *RedisIdempotencyStore) MarkProcessed(ctx context.Context, paymentID string) error {
	key := s.key(paymentID)

	if err := s.client.Set(ctx, key, "processed", s.ttl).Err(); err != nil {
		return err
	}

	log.Printf("[Idempotency] Notification marked as processed payment_id=%s ttl=%s", paymentID, s.ttl.String())
	return nil
}
