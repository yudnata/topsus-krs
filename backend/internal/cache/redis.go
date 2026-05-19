package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

// Client membungkus Redis dengan fallback graceful jika tidak tersedia.
type Client struct {
	rdb     *redis.Client
	enabled bool
}

func New(redisURL string) *Client {
	if redisURL == "" {
		log.Println("⚠ REDIS_URL kosong — cache dinonaktifkan (fallback DB only)")
		return &Client{enabled: false}
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("⚠ Redis URL tidak valid: %v — cache dinonaktifkan", err)
		return &Client{enabled: false}
	}

	rdb := redis.NewClient(opt)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("⚠ Redis tidak tersedia: %v — fallback DB only", err)
		_ = rdb.Close()
		return &Client{enabled: false}
	}

	log.Println("✓ Redis connected")
	return &Client{rdb: rdb, enabled: true}
}

func (c *Client) Enabled() bool {
	return c.enabled
}

func (c *Client) Ping(ctx context.Context) error {
	if !c.enabled {
		return errors.New("redis disabled")
	}
	return c.rdb.Ping(ctx).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, bool, error) {
	if !c.enabled {
		return "", false, nil
	}
	val, err := c.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return val, true, nil
}

func (c *Client) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if !c.enabled {
		return nil
	}
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	if !c.enabled || len(keys) == 0 {
		return nil
	}
	return c.rdb.Del(ctx, keys...).Err()
}

func (c *Client) GetJSON(ctx context.Context, key string, dest any) (bool, error) {
	raw, hit, err := c.Get(ctx, key)
	if err != nil || !hit {
		return false, err
	}
	if err := json.Unmarshal([]byte(raw), dest); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, string(b), ttl)
}

func (c *Client) Close() error {
	if c.rdb == nil {
		return nil
	}
	return c.rdb.Close()
}
