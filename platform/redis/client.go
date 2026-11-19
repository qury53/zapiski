package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
	timeout time.Duration
}

func (c *Client) Timeout() time.Duration {
	return c.timeout
}

func NewClient(ctx context.Context, config *Config) (*Client, error) {
	opts := redis.Options{
		Addr:       config.Addr,
		Password:   config.Password,
		DB:         config.DB,
		MaxRetries: config.MaxRetries,
	}

	client := redis.NewClient(&opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf(
			"failed to ping redis database: %w",
			err,
		)
	}

	return &Client{
		Client:  client,
		timeout: config.CacheTTL,
	}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
