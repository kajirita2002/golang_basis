package intercepctor

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kajirita2002/golang_basis/log"
	"github.com/kajirita2002/golang_basis/server/api/grpc/notifier"
)

const (
	callersNum = 16
	skipNum    = 3
)

// RecoverInterceptor adds defer recovery function before run a grpc.UnaryHandler.
func RecoverInterceptor(ntr notifier.BugNotifier) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if rErr := recovery(recover()); rErr != nil {
				err = rErr

				if ntr != nil {
					ntr.Notify(err, ctx)
				}
			}
		}()

		resp, err = handler(ctx, req)
		return resp, err
	}
}

// getCallers returns the callers of the function that calls it.
func getCallers() string {
	callers := make([]string, 0, callersNum)
	var pc [callersNum]uintptr

	n := runtime.Callers(skipNum, pc[:])
	for _, pc := range pc[:n] {
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pc)
		name := fn.Name()
		callers = append(callers, fmt.Sprintf("%s\n\t%s:%d", file, name, line))
	}
	return strings.Join(callers, "\n")
}

func recovery(catch interface{}) error {
	if catch == nil {
		return nil
	}
	callers := getCallers()
	log.Error("grpc: recovered from panic.", log.Fany("err", catch), log.Fstring("callers", callers))
	return status.Errorf(codes.Internal, "grpc: an internal error occurred")
}
