package grpc

import (
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"umbrella-go/umbrella-common/caller"
	pb "umbrella-go/umbrella-common/caller/grpc/test"
	"umbrella-go/umbrella-common/middleware/grpc"
)

type testServer struct{}

func (ts testServer) Echo(ctx context.Context, m *pb.EchoMsg) (*pb.EchoMsg, error) {
	return m, nil
}

func (ts testServer) EchoStream(stream pb.Echo_EchoStreamServer) error {
	m, err := stream.Recv()
	for err == nil {
		err = stream.Send(m)
		if err != nil {
			return err
		}
		m, err = stream.Recv()
	}

	if err == io.EOF {
		return nil
	} else {
		return err
	}
}

func newEchoServer(ui grpc.UnaryServerInterceptor, si grpc.StreamServerInterceptor) (*grpc.Server, net.Addr, error) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, nil, err
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(ui),
		grpc.StreamInterceptor(si),
	)
	pb.RegisterEchoServer(s, testServer{})

	go s.Serve(lis)
	return s, lis.Addr(), nil
}

func assertCallerNameUnary(t *testing.T, name string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		callerName := caller.CallerNameFromContext(ctx)
		assert.Equal(t, name, callerName)
		return handler(ctx, req)
	}
}

func assertCallerNameStream(t *testing.T, name string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		callerName := caller.CallerNameFromContext(ss.Context())
		assert.Equal(t, name, callerName)
		return handler(srv, ss)
	}
}

func sendOne(c pb.EchoClient, ctx context.Context, msg *pb.EchoMsg) error {
	// TODO 深入理解
	s, err := c.EchoStream(ctx)
	if err != nil {
		return err
	}
	if err = s.Send(msg); err != nil {
		return err
	}
	if err = s.CloseSend(); err != nil {
		return err
	}
	for {
		_, err := s.Recv()
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
	}
}

func TestCallerName(t *testing.T) {
	ui := grpcmiddleware.ChainUnaryServer(ExtractCallerNameUnary(), assertCallerNameUnary(t, "test"))
	si := grpcmiddleware.ChainStreamServer(ExtractCallerNameStream(), assertCallerNameStream(t, "test"))
	server, addr, err := newEchoServer(ui, si)
	if err != nil {
		t.Fatal(err)
	}
	// TODO 深入理解
	defer server.GracefulStop()

	conn, err := grpc.Dial(addr.String(),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(InjectCallerNameUnary("test")),
		grpc.WithStreamInterceptor(InjectCallerNameStream("test")),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := pb.NewEchoClient(conn)
	m := &pb.EchoMsg{"hi"}

	c.Echo(context.Background(), m)
	sendOne(c, context.Background(), m)
}
