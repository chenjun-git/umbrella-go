package pubsub

import (
	"context"

	"github.com/fzzy/radix/extra/pubsub"
)

type Receiver func(ctx context.Context) *pubsub.SubReply

type ReceiverWrapper func(next Receiver, ctx context.Context) *pubsub.SubReply

func (pw ReceiverWrapper) Wrap(next Receiver) Receiver {
	return func(ctx context.Context) *pubsub.SubReply {
		return pw(next, ctx)
	}
}

func identityReceiverWrapper(next Receiver, ctx context.Context) *pubsub.SubReply {
	return next(ctx)
}

func ChainReceiverWrappers(rws ...ReceiverWrapper) ReceiverWrapper {
	rws = removeNilReceiverWrapper(rws)
	switch len(rws) {
	case 0:
		return identityReceiverWrapper
	case 1:
		return rws[0]
	default:
		return func(next Receiver, ctx context.Context) *pubsub.SubReply {
			n := next
			for i := len(rws) - 1; i >= 0; i-- {
				n = rws[i].Wrap(n)
			}
			return n(ctx)
		}
	}
}

func removeNilReceiverWrapper(rws []ReceiverWrapper) []ReceiverWrapper {
	var result []ReceiverWrapper
	for _, rw := range rws {
		if rw != nil {
			result = append(result, rw)
		}
	}

	return result
}
