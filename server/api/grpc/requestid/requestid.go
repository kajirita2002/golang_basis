package requestid

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type requestID string

const (
	requestIDKey requestID = "x-request-id"
)

// UnaryServerInterceptor creates a new request id interceptor for unary request.
func UnaryServerInterceptor(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	rqID := GetRequestID(ctx)
	if rqID == "" {
		rqID = Generate()
	}
	return handler(SetRequestID(ctx, rqID), req)
}

func GetRequestID(ctx context.Context) string {
	if rd := ctx.Value(requestIDKey); rd != nil {
		return rd.(string)
	}
	return ""
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func Generate() string {
	return uuid.New().String()
}
