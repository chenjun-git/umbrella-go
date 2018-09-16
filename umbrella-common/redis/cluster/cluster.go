package cluster

import (
	"context"
	"time"

	"github.com/fzzy/radix/extra/cluster"
	radix "github.com/fzzy/radix/redis"

	"umbrella-go/umbrella-common/redis"
)

type (
	Opts       = cluster.Opts
	RawCluster = cluster.Cluster
)

type Cluster struct {
	cluster *cluster.Cluster
	cmder   redis.Cmder
	opts    []redis.Option
}

type Option func(*Cluster)

func WrapCmder(cw redis.CmderWrapper) Option {
	return func(c *Cluster) {
		c.cmder = cw.Wrap(c.defaultCmd)
		c.opts = append(c.opts, redis.WrapCmder(cw))
	}
}

func WrapPipeliner(pw redis.PipelinerWrapper) Option {
	return func(c *Cluster) {
		c.opts = append(c.opts, redis.WrapPipeliner(pw))
	}
}

func WrapCluster(c *cluster.Cluster, opts ...Option) *Cluster {
	cluster := &Cluster{
		cluster: c,
	}
	for _, opt := range opts {
		opt(cluster)
	}
	return cluster
}

func UnwrapCluster(c *Cluster) *cluster.Cluster {
	return c.cluster
}

func (c *Cluster) ClientForKey(key string) (*redis.Client, string, error) {
	client, addr, err := c.cluster.ClientForKey(key)
	if err != nil {
		return nil, "", err
	}
	return redis.WrapClient(client, c.opts...), addr, nil
}

func (c *Cluster) defaultCmd(ctx context.Context, cmd string, args ...interface{}) *radix.Reply {
	return c.cluster.Cmd(cmd, args...)
}

func (c *Cluster) Cmd(ctx context.Context, cmd string, args ...interface{}) *radix.Reply {
	if c.cmder != nil {
		return c.cmder(ctx, cmd, args...)
	}
	return c.defaultCmd(ctx, cmd, args...)
}

func (c *Cluster) Close() {
	c.Close()
}

func (c *Cluster) Reset() error {
	return c.Reset()
}

func NewCluster(addr string, opts ...Option) (*Cluster, error) {
	c, err := cluster.NewCluster(addr)
	if err != nil {
		return nil, err
	}
	return WrapCluster(c, opts...), nil
}

func NewClusterTimeout(addr string, timeout time.Duration, opts ...Option) (*Cluster, error) {
	c, err := cluster.NewClusterTimeout(addr, timeout)
	if err != nil {
		return nil, err
	}
	return WrapCluster(c, opts...), nil
}

func NewClusterWithOpts(o cluster.Opts, opts ...Option) (*Cluster, error) {
	c, err := cluster.NewClusterWithOpts(o)
	if err != nil {
		return nil, err
	}
	return WrapCluster(c, opts...), nil
}
