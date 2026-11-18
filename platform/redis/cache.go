package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func (c *Client) Get(
	ctx context.Context,
	key string,
	dst any,
) (bool, error) {
	data, err := c.Client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf("failed to get cache: %w", err)
	}

	return true, json.Unmarshal(data, dst)
}

func (c *Client) Set(
	ctx context.Context,
	key string,
	value any,
) error {
	return c.SetWithTTL(ctx, key, value, c.Timeout())
}

// SetWithTTL сохраняет значение с явным TTL, не завязываясь на общий CacheTTL.
func (c *Client) SetWithTTL(
	ctx context.Context,
	key string,
	value any,
	ttl time.Duration,
) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal data for cache: %w", err)
	}

	return c.Client.Set(ctx, key, data, ttl).Err()
}

func (c *Client) Del(
	ctx context.Context,
	key string,
) error {
	return c.Client.Del(ctx, key).Err()
}

func (c *Client) GenerateCacheKey(args ...string) string {
	return strings.Join(args, ":")
}

// Incr увеличивает значение по ключу
func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.Client.Incr(ctx, key).Result()
}

func (c *Client) GetDel(
	ctx context.Context,
	key string,
	dst any,
) (bool, error) {
	data, err := c.Client.GetDel(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, fmt.Errorf("failed to GetDel: %w", err)
	}

	return true, json.Unmarshal(data, dst)
}
