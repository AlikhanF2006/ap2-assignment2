package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"order-service/internal/domain"
)

var ErrCacheMiss = errors.New("order not found in cache")

type RedisOrderCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisOrderCache(addr string, password string, db int, ttl time.Duration) *RedisOrderCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisOrderCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *RedisOrderCache) key(id string) string {
	return fmt.Sprintf("order:%s", id)
}

func (c *RedisOrderCache) Get(ctx context.Context, id string) (*domain.Order, error) {
	key := c.key(id)

	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Printf("CACHE MISS order_id=%s", id)
			return nil, ErrCacheMiss
		}
		return nil, err
	}

	var order domain.Order
	if err := json.Unmarshal([]byte(value), &order); err != nil {
		return nil, err
	}

	log.Printf("CACHE HIT order_id=%s", id)
	return &order, nil
}

func (c *RedisOrderCache) Set(ctx context.Context, order *domain.Order) error {
	key := c.key(order.ID)

	data, err := json.Marshal(order)
	if err != nil {
		return err
	}

	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return err
	}

	log.Printf("CACHE SET order_id=%s ttl=%s", order.ID, c.ttl.String())
	return nil
}

func (c *RedisOrderCache) Delete(ctx context.Context, id string) error {
	key := c.key(id)

	if err := c.client.Del(ctx, key).Err(); err != nil {
		return err
	}

	log.Printf("CACHE INVALIDATED order_id=%s", id)
	return nil
}
