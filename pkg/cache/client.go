package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Client interface {
	Set(key string, value []byte, expiration time.Duration) error
	Get(key string) (value []byte, err error)
}

type prefixKeyClient struct {
	prefixKey string
	client    Client
}

func WithPrefixClient(prefixKey string, client Client) Client {
	return &prefixKeyClient{
		prefixKey: prefixKey,
		client:    client,
	}
}

func (p *prefixKeyClient) Set(key string, value []byte, expiration time.Duration) error {
	return p.client.Set(p.prefixKey+key, value, expiration)
}

func (p *prefixKeyClient) Get(key string) (value []byte, err error) {
	return p.client.Get(key)
}

type redisClient struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) Client {
	return &redisClient{client: client}
}

var noCtx = context.Background()

func (r *redisClient) Set(key string, value []byte, expiration time.Duration) error {
	return r.client.Set(noCtx, key, string(value), expiration).Err()
}

func (r *redisClient) Get(key string) (value []byte, err error) {
	value, err = r.client.Get(noCtx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	} else {
		return value, nil
	}
}
