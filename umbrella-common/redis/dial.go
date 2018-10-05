package redis

import (
	"time"

	"github.com/fzzy/radix/redis"
)

type Option func(*Client)

func WrapClient(c *redis.Client, opts ...Option) *Client {
	client := &Client{
		client: c,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func UnwrapClient(c *Client) *redis.Client {
	return c.client
}

func Dial(network, addr string, opts ...Option) (*Client, error) {
	c, err := redis.Dial(network, addr)
	if err != nil {
		return nil, err
	}
	return WrapClient(c, opts...), nil
}

func DialTimeout(network, addr string, timeout time.Duration, opts ...Option) (*Client, error) {
	c, err := redis.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}
	return WrapClient(c, opts...), nil
}
