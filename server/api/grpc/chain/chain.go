package chain

import (
	"context"

	"google.golang.org/grpc"
)

// ChainUnaryServer chains a list of unary interceptors to a single unary interceptor.
func ChainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		buildChain := func(curInter grpc.UnaryServerInterceptor, nextHandler grpc.UnaryHandler) grpc.UnaryHandler {
			return func(curCtx context.Context, curReq interface{}) (interface{}, error) {
				return curInter(curCtx, curReq, info, nextHandler)
			}
		}

		chainedHandler := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			chainedHandler = buildChain(interceptors[i], chainedHandler)
		}
		return chainedHandler(ctx, req)
	}
}
