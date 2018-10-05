package grpc

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"umbrella-go/umbrella-common/caller"
	"umbrella-go/umbrella-common/middleware/grpc"
)

const (
	callerName = "caller-name"
)

func injectCallerName(ctx context.Context, name string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if ok {
		return metadata.NewOutgoingContext(ctx, metadata.Join(md, metadata.Pairs(callerName, name)))
	} else {
		return metadata.NewOutgoingContext(ctx, metadata.Pairs(callerName, name))
	}
}

func InjectCallerNameUnary(name string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(injectCallerName(ctx, name), method, req, reply, cc, opts...)
	}
}

func InjectCallerNameStream(name string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		return streamer(injectCallerName(ctx, name), desc, cc, method, opts...)
	}
}

func extractCallerName(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	vs := md[callerName]
	if len(vs) == 0 {
		return ""
	}

	return vs[0]
}

func contextWithCallerName(ctx context.Context) context.Context {
	name := extractCallerName(ctx)

	if name != "" {
		return caller.ContextWithCallerName(ctx, name)
	} else {
		return ctx
	}
}

func ExtractCallerNameUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		return handler(contextWithCallerName(ctx), req)
	}
}

func ExtractCallerNameStream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := contextWithCallerName(ss.Context())
		return handler(srv, grpcmiddleware.ServerStreamWithContext(ss, ctx))
	}
}
