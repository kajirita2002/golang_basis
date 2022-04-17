package grpc

import (
	"math"
	"net"
	"strings"
	"time"

	auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/kajirita2002/golang_basis/server/api/grpc/chain"
	"github.com/kajirita2002/golang_basis/server/api/grpc/intercepctor"
	"github.com/kajirita2002/golang_basis/server/api/grpc/requestid"

	agrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"

	"github.com/kajirita2002/golang_basis/log"
	"github.com/kajirita2002/golang_basis/server/api/grpc/notifier"
)

const (
	gracefulPeriod  = 5
	keepAlivePeriod = 20
)

type options struct {
	gracefulStop   bool
	gracefulPeriod time.Duration
	interceptors   []agrpc.UnaryServerInterceptor
	panicNotifier  notifier.BugNotifier
	kaPolicy       keepalivePolicy
}

type keepalivePolicy struct {
	minTime             time.Duration
	permitWithoutStream bool
}

// Options is a functional options for configuring grpc server.
type Options func(*options)

// WithGracefulPeriod returns an Options which configures graceful period.
func WithGracefulPeriod(period time.Duration) Options {
	return func(o *options) {
		o.gracefulPeriod = period
	}
}

// WithGracefulStop returns an Options which configures whether gracefulStop is enabled or disabled.
func WithGracefulStop(enabled bool) Options {
	return func(o *options) {
		o.gracefulStop = enabled
	}
}

// WithInterceptors returns an Options whitch configures interceptors.
func WithInterceptors(interceptors ...agrpc.UnaryServerInterceptor) Options {
	return func(o *options) {
		o.interceptors = interceptors
	}
}

// WithPanicNotifier returns an Options which configures whether bugsnag notifier
func WithPanicNotifier(ntr notifier.BugNotifier) Options {
	return func(o *options) {
		o.panicNotifier = ntr
	}
}

// WithKeepalivePolicy returns Options which sets keepalive policy.
func WithKeepalivePolicy(minTime time.Duration, permit bool) Options {
	return func(o *options) {
		o.kaPolicy = keepalivePolicy{
			minTime:             minTime,
			permitWithoutStream: permit,
		}
	}
}

// Register defines a function interface for registering grpc service.
type Register func(*agrpc.Server)

// Server represents a grpc server.
type Server struct {
	addr       string
	opts       *options
	grpcServer *agrpc.Server
}

func NewServer(addr string, register Register, opts ...Options) *Server {
	defaultOpts := &options{
		gracefulStop:   true,
		gracefulPeriod: gracefulPeriod * time.Second,
		interceptors:   nil,
		panicNotifier:  nil,
		kaPolicy: keepalivePolicy{
			minTime:             keepAlivePeriod * time.Second,
			permitWithoutStream: true,
		},
	}
	for _, opt := range opts {
		opt(defaultOpts)
	}

	s := &Server{
		addr: addr,
		opts: defaultOpts,
	}

	s.grpcServer = newServerWithInterceptors(s.opts)
	register(s.grpcServer)
	return s
}

// ListenAndServe listens on a TCP network address and serves inconming RPCs.
func (s *Server) ListenAndServe() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	log.Infof("grpcserver: running on %s", s.addr)
	err = s.grpcServer.Serve(l)
	if strings.Contains(status.Convert(err).Message(), "use of closed network connection") {
		return nil
	}
	return err
}

func newServerWithInterceptors(opts *options) *agrpc.Server {
	interceptors := []agrpc.UnaryServerInterceptor{
		intercepctor.NewAuthTokenPropagator(),
		auth.UnaryServerInterceptor(intercepctor.DefaultAuthentication()),
		intercepctor.RecoverInterceptor(opts.panicNotifier),
		requestid.UnaryServerInterceptor,
	}
	unaryInters := chain.ChainUnaryServer(interceptors...)
	serverOptions := []agrpc.ServerOption{
		agrpc.UnaryInterceptor(unaryInters),
		agrpc.MaxSendMsgSize(math.MaxInt32),
		agrpc.MaxRecvMsgSize(math.MaxInt32),
		agrpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             opts.kaPolicy.minTime,
			PermitWithoutStream: opts.kaPolicy.permitWithoutStream,
		}),
	}
	return agrpc.NewServer(serverOptions...)
}
