package grpcmiddleware

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// ChainUnaryServer creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryServer(one, two, three) will execute one before two before three, and three
// will see context changes of one and two.
func ChainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	interceptors = removeNilUnaryServerInterceptor(interceptors)
	switch len(interceptors) {
	case 0:
		return nil
	case 1:
		return interceptors[0]
	default:
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			buildChain := func(current grpc.UnaryServerInterceptor, next grpc.UnaryHandler) grpc.UnaryHandler {
				return func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
					return current(currentCtx, currentReq, info, next)
				}
			}
			chain := handler
			for i := len(interceptors) - 1; i >= 0; i-- {
				chain = buildChain(interceptors[i], chain)
			}
			return chain(ctx, req)
		}
	}
}

func removeNilUnaryServerInterceptor(interceptors []grpc.UnaryServerInterceptor) []grpc.UnaryServerInterceptor {
	var result []grpc.UnaryServerInterceptor
	for _, interceptor := range interceptors {
		if interceptor != nil {
			result = append(result, interceptor)
		}
	}

	return result
}

// ChainStreamServer creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryServer(one, two, three) will execute one before two before three.
// If you want to pass context between interceptors, use WrapServerStream.
func ChainStreamServer(interceptors ...grpc.StreamServerInterceptor) grpc.StreamServerInterceptor {
	interceptors = removeNilStreamServerInterceptor(interceptors)
	switch len(interceptors) {
	case 0:
		return nil
	case 1:
		return interceptors[0]
	default:
		return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			buildChain := func(current grpc.StreamServerInterceptor, next grpc.StreamHandler) grpc.StreamHandler {
				return func(currentSrv interface{}, currentStream grpc.ServerStream) error {
					return current(currentSrv, currentStream, info, next)
				}
			}
			chain := handler
			for i := len(interceptors) - 1; i >= 0; i-- {
				chain = buildChain(interceptors[i], chain)
			}
			return chain(srv, stream)
		}
	}
}

func removeNilStreamServerInterceptor(interceptors []grpc.StreamServerInterceptor) []grpc.StreamServerInterceptor {
	var result []grpc.StreamServerInterceptor
	for _, interceptor := range interceptors {
		if interceptor != nil {
			result = append(result, interceptor)
		}
	}

	return result
}

func ChainUnaryClient(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	interceptors = removeNilUnaryClientInterceptor(interceptors)
	switch len(interceptors) {
	case 0:
		return nil
	case 1:
		return interceptors[0]
	default:
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			buildChain := func(current grpc.UnaryClientInterceptor, next grpc.UnaryInvoker) grpc.UnaryInvoker {
				return func(currentCtx context.Context, currentMethod string, currentReq, currentRepl interface{}, currentConn *grpc.ClientConn, currentOpts ...grpc.CallOption) error {
					return current(currentCtx, currentMethod, currentReq, currentRepl, currentConn, next, currentOpts...)
				}
			}
			chain := invoker
			for i := len(interceptors) - 1; i >= 0; i-- {
				chain = buildChain(interceptors[i], chain)
			}
			return chain(ctx, method, req, reply, cc, opts...)
		}
	}
}

func removeNilUnaryClientInterceptor(interceptors []grpc.UnaryClientInterceptor) []grpc.UnaryClientInterceptor {
	var result []grpc.UnaryClientInterceptor
	for _, interceptor := range interceptors {
		if interceptor != nil {
			result = append(result, interceptor)
		}
	}

	return result
}

// ChainStreamClient creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainStreamClient(one, two, three) will execute one before two before three.
func ChainStreamClient(interceptors ...grpc.StreamClientInterceptor) grpc.StreamClientInterceptor {
	interceptors = removeNilStreamClientInterceptor(interceptors)
	switch len(interceptors) {
	case 0:
		return nil
	case 1:
		return interceptors[0]
	default:
		return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			buildChain := func(current grpc.StreamClientInterceptor, next grpc.Streamer) grpc.Streamer {
				return func(currentCtx context.Context, currentDesc *grpc.StreamDesc, currentConn *grpc.ClientConn, currentMethod string, currentOpts ...grpc.CallOption) (grpc.ClientStream, error) {
					return current(currentCtx, currentDesc, currentConn, currentMethod, next, currentOpts...)
				}
			}
			chain := streamer
			for i := len(interceptors) - 1; i >= 0; i-- {
				chain = buildChain(interceptors[i], chain)
			}
			return chain(ctx, desc, cc, method, opts...)
		}
	}
}

func removeNilStreamClientInterceptor(interceptors []grpc.StreamClientInterceptor) []grpc.StreamClientInterceptor {
	var result []grpc.StreamClientInterceptor
	for _, interceptor := range interceptors {
		if interceptor != nil {
			result = append(result, interceptor)
		}
	}

	return result
}

// WithUnaryServerChain is a grpc.Server config option that accepts multiple unary interceptors.
// Basically syntactic sugar.
func WithUnaryServerChain(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	return grpc.UnaryInterceptor(ChainUnaryServer(interceptors...))
}

// WithStreamServerChain is a grpc.Server config option that accepts multiple stream interceptors.
// Basically syntactic sugar.
func WithStreamServerChain(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	return grpc.StreamInterceptor(ChainStreamServer(interceptors...))
}

func WithUnaryClientChain(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithUnaryInterceptor(ChainUnaryClient(interceptors...))
}

func WithStreamClientChain(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return grpc.WithStreamInterceptor(ChainStreamClient(interceptors...))
}

type wrappedServerStream struct {
	grpc.ServerStream
	context context.Context
}

func ServerStreamWithContext(ss grpc.ServerStream, ctx context.Context) grpc.ServerStream {
	wss, ok := ss.(*wrappedServerStream)
	if ok {
		return &wrappedServerStream{
			ServerStream: wss.ServerStream,
			context:      ctx,
		}
	}

	return &wrappedServerStream{
		ServerStream: ss,
		context:      ctx,
	}
}

func (wss *wrappedServerStream) Context() context.Context {
	return wss.context
}
