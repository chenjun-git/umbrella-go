package redis


type Cmder func(ctx context.Context, cmd string, args []interface{}) *redis.Reply

type CmderWrapper func(next Cmder, ctx context.Context, cmd string, args []interface{}) *redis.Reply

func (cw CmderWrapper) Wrap(next Cmder) Cmder {
	if cw == nil {
		return next
	}

	return func(ctx context.Context, cmd string, args []interface{}) *redis.Reply {
		return cw(next, ctx, cmd, args)
	}
}

func (cw CmderWrapper) Chain(cw2 CmderWrapper) CmderWrapper {
	if cw == nil {
		return cw2
	}
	if cw2 == nil {
		return cw
	}

	return func(next Cmder, ctx context.Context, cmd string, args []interface{}) *redis.Reply {
		return cw(cw2.Wrap(next), ctx, cmd, args)
	}
}

func identityCmderWrapper(next Cmder, ctx context.Context, cmd string, args []interface{}) *redis.Reply {
	return next(ctx, cmd, args)
}

func ChainCmderWrappers(cws ...CmderWrapper) CmderWrapper {
	cws = removeNilCmderWrapper(cws)
	switch len(cws) {
	case 0:
		return identityCmderWrapper
	case 1:
		return cws[0]
	default:
		return func(next Cmder, ctx context.Context, cmd string, args []interface{}) *redis.Reply {
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