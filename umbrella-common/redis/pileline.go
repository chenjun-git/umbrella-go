package redis

import (
	"context"

	"git.meiqia.com/triones/compass/redis"
)


type Pipeliner func(ctx context.Context, reqs []*Request) []*redis.Reply

type PipelinerWrapper func(next Pipeliner, ctx context.Context, reqs []*Request) []*redis.Reply

func (pw PipelinerWrapper) Wrap(next Pipeliner) Pipeliner {
	if pw == nil {
		return next
	}

	return func(ctx context.Context, reqs []*Request) []*redis.Reply {
		return pw(next, ctx, reqs)
	}
}

func (pw PipelinerWrapper) Chain(pw2 PipelinerWrapper) PipelinerWrapper {
	if pw == nil {
		return pw2
	}
	if pw2 == nil {
		return pw
	}

	return func(next Pipeliner, ctx context.Context, reqs []*Request) []*redis.Reply {
		return pw(pw2.Wrap(next), ctx, reqs)
	}
}

func identityPipelinerWrapper(next Pipeliner, ctx context.Context, reqs []*Request) []*redis.Reply {
	return next(ctx, reqs)
}

func ChainPipelinerWrappers(pws ...PipelinerWrapper) PipelinerWrapper {
	pws = removeNilPipelinerWrapper(pws)
	switch len(pws) {
	case 0:
		return identityPipelinerWrapper
	case 1:
		return pws[0]
	default:
		return func(next Pipeliner, ctx context.Context, reqs []*Request) []*redis.Reply {
			n := next
			for i := len(pws) - 1; i >= 0; i-- {
				n = pws[i].Wrap(n)
			}
			return n(ctx, reqs)
		}
	}
}

func removeNilPipelinerWrapper(pws []PipelinerWrapper) []PipelinerWrapper {
	var result []PipelinerWrapper
	for _, pw := range pws {
		if pw != nil {
			result = append(result, pw)
		}
	}

	return result
}