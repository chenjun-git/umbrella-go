package pubsub

import (
	"context"

	"github.com/fzzy/radix/extra/pubsub"
)

type Cmder func(ctx context.Context, cmd string, patterns []interface{}) *pubsub.SubReply

type CmderWrapper func(next Cmder, ctx context.Context, cmd string, patterns []interface{}) *pubsub.SubReply

func (cw CmderWrapper) Wrap(next Cmder) Cmder {
	return func(ctx context.Context, cmd string, patterns []interface{}) *pubsub.SubReply {
		return cw(next, ctx, cmd, patterns)
	}
}

func identityCmderWrapper(next Cmder, ctx context.Context, cmd string, patterns []interface{}) *pubsub.SubReply {
	return next(ctx, cmd, patterns)
}

func ChainCmderWrappers(cws ...CmderWrapper) CmderWrapper {
	cws = removeNilCmderWrapper(cws)
	switch len(cws) {
	case 0:
		return identityCmderWrapper
	case 1:
		return cws[0]
	default:
		return func(next Cmder, ctx context.Context, cmd string, args []interface{}) *pubsub.SubReply {
			n := next
			for i := len(cws) - 1; i >= 0; i-- {
				n = cws[i].Wrap(n)
			}
			return n(ctx, cmd, args)
		}
	}
}

func removeNilCmderWrapper(cws []CmderWrapper) []CmderWrapper {
	var result []CmderWrapper
	for _, cw := range cws {
		if cw != nil {
			result = append(result, cw)
		}
	}

	return result
}
