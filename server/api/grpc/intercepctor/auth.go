package intercepctor

import (
	"context"
	"errors"

	"github.com/kajirita2002/golang_basis/log"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const authorizationKey = "x-authorization"

var errTemp = errors.New("temp")

type DefaultAuthFunc func(ctx context.Context) (context.Context, error)

func DefaultAuthentication() grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		authorization, err := Authorization(ctx)
		if err != nil {
			return nil, err
		}

		err = verify(authorization)
		if err != nil {
			return nil, err
		}

		return ctx, nil
	}
}

// Authorization は gRPC メタデータからユーザー認証トークンを取得する
func Authorization(ctx context.Context) (string, error) {
	return fromMeta(ctx, authorizationKey)
}

func fromMeta(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errTemp
	}
	vs := md[key]
	if len(vs) == 0 {
		return "", errTemp
	}
	return vs[0], nil
}

func verify(token string) error {
	// TODO: 認証処理を書く
	log.Infof(token)
	return nil
}

// ServiceAuthorize はサービスごとの認可を行う関数を実装するインターフェースを表す。
type ServiceAuthorize interface {
	Authorize(context.Context, string) error
}

func NewAuthTokenPropagator() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			authmd := md.Get(authorizationKey)
			if len(authmd) > 0 {
				ctx = metadata.AppendToOutgoingContext(ctx, authorizationKey, authmd[0])
			}
		}

		return handler(ctx, req)
	}
}
