package caller

import "context"

type callerNameKey struct{}

func ContextWithCallerName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, callerNameKey{}, name)
}

func CallerNameFromContext(ctx context.Context) string {
	name, ok := ctx.Value(callerNameKey{}).(string)
	if !ok {
		return ""
	}
	return name
}
