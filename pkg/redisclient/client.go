// Package redisclient handles Redis connection setup only.
// No queue logic, no business logic — just a connected client.
package redisclient

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// New creates a Redis client from a URL and verifies the connection.
// URL format: redis://:password@host:port/db
// Simple local: redis://localhost:6379

type Client struct {
	rdb *redis.Client
}

func New(redisURL string) (*Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("connect to redis: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

func (c *Client) Redis() *redis.Client {
	return c.rdb
}

func (c *Client) Close() error {
	return c.rdb.Close()
}
