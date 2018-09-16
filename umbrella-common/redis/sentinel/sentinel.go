package sentinel

import (
	"github.com/fzzy/radix/extra/sentinel"
)

type (
	RawClient = sentinel.Client
)

type Client struct {
	client *sentinel.Client
	opts   []redis.Option
}

func (c *Client) CarefullyPutMaster(name string, client *redis.Client, potentialErr *error) {
	c.client.CarefullyPutMaster(name, redis.UnwrapClient(client), potentialErr)
}

func (c *Client) GetMaster(name string) (*redis.Client, error) {
	client, err := c.client.GetMaster(name)
	if err != nil {
		return nil, err
	}
	return redis.WrapClient(client, c.opts...), nil
}

func (c *Client) PutMaster(name string, client *redis.Client) {
	c.client.PutMaster(name, redis.UnwrapClient(client))
}

func (c *Client) Close() {
	c.Close()
}

func WrapClient(c *sentinel.Client, opts ...redis.Option) *Client {
	return &Client{
		client: c,
		opts:   opts,
	}
}

func UnwrapClient(c *Client) *sentinel.Client {
	return c.client
}

func NewClient(network, address string, poolSize int, names ...string) (*Client, error) {
	c, err := sentinel.NewClient(network, address, poolSize, names...)
	if err != nil {
		return nil, err
	}
	return WrapClient(c), nil
}

func NewClientWithOptions(network, address string, poolSize int, names []string, opts ...redis.Option) (*Client, error) {
	c, err := sentinel.NewClient(network, address, poolSize, names...)
	if err != nil {
		return nil, err
	}
	return WrapClient(c, opts...), nil
}
