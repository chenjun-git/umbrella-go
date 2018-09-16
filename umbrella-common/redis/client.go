package redis

import (
	"context"
	"fmt"
	"strings"

	"github.com/fzzy/radix/redis"
)

type (
	Reply     = redis.Reply
	ReplyType = redis.ReplyType
	RawClient = redis.Client
)

var (
	RawDial        = redis.Dial
	RawDialTimeout = redis.DialTimeout
)

type Client struct {
	client    *redis.Client
	cmder     Cmder
	pipeliner Pipeliner
	pending   []*Request
	completed []*redis.Reply
}

type Request struct {
	Cmd  string
	Args []interface{}
}

func (req *Request) String() string {
	ss := make([]string, len(req.Args)+1)
	ss[0] = req.Cmd
	for i := 0; i < len(req.Args); i++ {
		ss[i+1] = fmt.Sprint(req.Args[i])
	}
	return strings.Join(ss, " ")
}

func NewRequest(cmd string, args ...interface{}) *Request {
	return &Request{
		Cmd:  cmd,
		Args: args,
	}
}

func (c *Client) defaultCmd(ctx context.Context, cmd string, args []interface{}) *redis.Reply {
	return c.client.Cmd(cmd, args...)
}

func (c *Client) Cmd(ctx context.Context, cmd string, args ...interface{}) *redis.Reply {
	if c.cmder != nil {
		return c.cmder(ctx, cmd, args)
	}
	return c.defaultCmd(ctx, cmd, args)
}

func (c *Client) defaultPipeline(ctx context.Context, reqs []*Request) []*redis.Reply {
	for _, req := range reqs {
		c.client.Append(req.Cmd, req.Args)
	}
	reps := make([]*redis.Reply, len(reqs))
	for i := 0; i < len(reqs); i++ {
		reps[i] = c.client.GetReply()
	}
	return reps
}

func (c *Client) Pipeline(ctx context.Context, reqs []*Request) []*redis.Reply {
	if c.pipeliner != nil {
		return c.pipeliner(ctx, reqs)
	}
	return c.defaultPipeline(ctx, reqs)
}

func (c *Client) Append(ctx context.Context, cmd string, args ...interface{}) {
	c.pending = append(c.pending, &Request{cmd, args})
}

func (c *Client) GetReply(ctx context.Context) *redis.Reply {
	if len(c.completed) > 0 {
		r := c.completed[0]
		c.completed = c.completed[1:]
		return r
	}
	c.completed = nil

	if len(c.pending) == 0 {
		return &redis.Reply{Type: redis.ErrorReply, Err: redis.PipelineQueueEmptyError}
	}

	c.completed = c.Pipeline(ctx, c.pending)
	r := c.completed[0]
	c.completed = c.completed[1:]
	return r
}

func (c *Client) Close() {
	c.client.Close()
}
