package pubsub

import (
	"context"
	"errors"

	"github.com/fzzy/radix/extra/pubsub"

	"umbrella-go/umbrella-common/redis"
)

const (
	SUBSCRIBE    = "SUBSCRIBE"
	PSUBSCRIBE   = "PSUBSCRIBE"
	UNSUBSCRIBE  = "UNSUBSCRIBE"
	PUNSUBSCRIBE = "PUNSUBSCRIBE"
)

type (
	SubReply     = pubsub.SubReply
	SubReplyType = pubsub.SubReplyType
	RawSubClient = pubsub.SubClient
)

const (
	ErrorReply       = pubsub.ErrorReply
	SubscribeReply   = pubsub.SubscribeReply
	UnsubscribeReply = pubsub.UnsubscribeReply
	MessageReply     = pubsub.MessageReply
)

var (
	UnknownPubsubCmdError = errors.New("unknown pubsub command")
)

type SubClient struct {
	subClient *pubsub.SubClient
	cmder     Cmder
	receiver  Receiver
}

type Option func(c *SubClient)

func WrapCmder(cw CmderWrapper) Option {
	return func(c *SubClient) {
		if cw != nil {
			c.cmder = cw.Wrap(c.defaultCmd)
		}
	}
}

func WrapReceiver(rw ReceiverWrapper) Option {
	return func(c *SubClient) {
		if rw != nil {
			c.receiver = rw.Wrap(c.defaultReceive)
		}
	}
}

func WrapSubClient(c *pubsub.SubClient, opts ...Option) *SubClient {
	client := &SubClient{
		subClient: c,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

func UnwrapSubClient(c *SubClient) *pubsub.SubClient {
	return c.subClient
}

func (c *SubClient) defaultCmd(ctx context.Context, cmd string, patterns []interface{}) *pubsub.SubReply {
	cmd = strings.ToUpper(cmd)
	switch cmd {
	case PSUBSCRIBE:
		return c.subClient.PSubscribe(patterns...)
	case PUNSUBSCRIBE:
		return c.subClient.PUnsubscribe(patterns...)
	case SUBSCRIBE:
		return c.subClient.Subscribe(patterns...)
	case UNSUBSCRIBE:
		return c.subClient.Unsubscribe(patterns...)
	default:
		panic(UnknownPubsubCmdError)
	}
}

func (c *SubClient) cmd(ctx context.Context, cmd string, patterns []interface{}) *pubsub.SubReply {
	if c.cmder != nil {
		return c.cmder(ctx, cmd, patterns)
	}
	return c.defaultCmd(ctx, cmd, patterns)
}

func (c *SubClient) PSubscribe(ctx context.Context, patterns ...interface{}) *pubsub.SubReply {
	return c.cmd(ctx, PSUBSCRIBE, patterns)
}

func (c *SubClient) PUnsubscribe(ctx context.Context, patterns ...interface{}) *pubsub.SubReply {
	return c.cmd(ctx, PUNSUBSCRIBE, patterns)
}

func (c *SubClient) Subscribe(ctx context.Context, patterns ...interface{}) *pubsub.SubReply {
	return c.cmd(ctx, SUBSCRIBE, patterns)
}

func (c *SubClient) Unsubscribe(ctx context.Context, patterns ...interface{}) *pubsub.SubReply {
	return c.cmd(ctx, UNSUBSCRIBE, patterns)
}

func (c *SubClient) defaultReceive(ctx context.Context) *pubsub.SubReply {
	return c.subClient.Receive()
}

func (c *SubClient) Receive(ctx context.Context) *pubsub.SubReply {
	if c.receiver != nil {
		return c.receiver(ctx)
	}
	return c.defaultReceive(ctx)
}

func NewSubClient(client *redis.Client, opts ...Option) *SubClient {
	return WrapSubClient(pubsub.NewSubClient(redis.UnwrapClient(client)), opts...)
}
