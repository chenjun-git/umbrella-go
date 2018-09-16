package grpcmiddleware

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"umbrella-go/umbrella-common/lang"
	proto "umbrella-go/umbrella-common/proto"
)

type errorGetter interface {
	GetError() *proto.Error
}

type ErrorMsgGetter func(code int, languages []string) string

func MakeUnaryServerErrorTranslator(errorMsgGetter ErrorMsgGetter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		languages := lang.FromIncomingContext(ctx)
		ctx = lang.ContextSetLanguages(ctx, languages)
		resp, err = handler(ctx, req)
		if resp, ok := resp.(errorGetter); ok {
			err := resp.GetError()

			if err != nil {
				// 设置err的message信息
				if err.Message == "" {
					if msg := errorMsgGetter(int(err.Code), languages); msg != "" {
						err.Message = msg
					} else {
						err.Message = "Unknown error"
					}
				}
			}
		}
		return resp, err
	}
}
