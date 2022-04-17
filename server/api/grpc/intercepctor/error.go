package intercepctor

import (
	"context"

	"github.com/bugsnag/bugsnag-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/kajirita2002/golang_basis/code"
	"github.com/kajirita2002/golang_basis/log"
	"github.com/kajirita2002/golang_basis/server/api/grpc/notifier"
	"github.com/kajirita2002/golang_basis/server/api/grpc/pb"
	acontext "github.com/kajirita2002/golang_basis/server/api/grpc/requestid"
)

func ErrorHandleInterceptor(bugNotifier notifier.BugNotifier) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		id := acontext.GetRequestID(ctx)
		res, err := handler(ctx, req)
		if err == nil {
			return res, nil
		}

		lg := log.With(
			log.Fstring("x-request-id", id),
			log.Fstring("method", info.FullMethod),
			log.Fany("request", req),
		)

		var userID string
		md, _ := metadata.FromIncomingContext(ctx)
		userIDs := md.Get("userid")
		if len(userIDs) > 0 {
			userID = userIDs[0]
		}

		c := code.GetCode(err)
		gcode := toGrpcCode(c)

		switch gcode {
		case codes.Internal, codes.Unknown:
			lg.Error(err.Error(), log.Ferror(err))
			if berr := bugNotifier.Notify(err, ctx, bugsnag.User{Id: userID}); berr != nil {
				log.Error("interface/grpc/interceptor: failed to notify to bugsnag", log.Ferror(berr))
			}
		default:
			lg.Warn(err.Error(), log.Ferror(err))
		}

		st := status.New(gcode, string(c))
		dcode := toGrpcDetailError(c)
		dt, err := st.WithDetails(
			&pb.ErrorDetail{Code: dcode, Message: string(c)},
		)
		if err != nil {
			return res, status.Error(codes.Internal, "failed to attach error detail")
		}

		return res, dt.Err()
	}
}

// toGrpcCode converts internal error code to gRPC error code.
func toGrpcCode(c code.Code) codes.Code {
	var gcode codes.Code
	switch c {
	case code.InvalidArgument:
		gcode = codes.InvalidArgument
	case code.Forbidden, code.ContentExpired:
		gcode = codes.PermissionDenied
	case code.NotFound:
		gcode = codes.NotFound
	case code.Unexpected:
		gcode = codes.Internal
	case code.Unknown:
		gcode = codes.Unknown
	default:
		gcode = codes.Unknown
	}
	return gcode
}

func toGrpcDetailError(c code.Code) pb.ErrorCode {
	switch c {
	case code.Forbidden:
		return pb.ErrorCode_FORBIDDEN
	case code.NotFound:
		return pb.ErrorCode_NOT_FOUND
	case code.InvalidArgument:
		return pb.ErrorCode_INVALID_ARGUMENT
	case code.ContentExpired:
		return pb.ErrorCode_CONTENT_EXPIRED
	case code.Unexpected:
		return pb.ErrorCode_UNEXPECTED
	default:
		return pb.ErrorCode_UNKNOWN
	}
}
