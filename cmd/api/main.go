package main

import (
	"github.com/kajirita2002/golang_basis/config"
	"github.com/kajirita2002/golang_basis/log"
	"github.com/kajirita2002/golang_basis/server/api/grpc"
	"github.com/kajirita2002/golang_basis/server/api/grpc/intercepctor"
	"github.com/kajirita2002/golang_basis/server/api/grpc/notifier"
)

func main() {
	srv, err := injector()
	if err != nil {
		log.Panic("failed to init service server", log.Ferror(err))
	}
	srvConf, err := config.NewService()
	if err != nil {
		log.Panic("failed to read config of GCP")
	}

	bugsnagConf, err := config.NewBugsnag()
	if err != nil {
		log.Panic("failed to read config of bugsag", log.Ferror(err))
	}
	isTemporaryEnv := srvConf.Env != config.DEV && srvConf.Env != config.PRD && srvConf.Env != config.LOAD
	grpcAddr := ":30990"
	noti := notifier.NewBugsnagNotifier(srvConf, bugsnagConf)
	grpcServer := grpc.NewServer(
		grpcAddr,
		registerServicesToGRPCServer(srv),
		grpc.WithGracefulStop(!isTemporaryEnv),
		grpc.WithInterceptors(intercepctor.ErrorHandleInterceptor(noti)),
		grpc.WithPanicNotifier(noti),
	)
	if err = grpcServer.ListenAndServe(); err != nil {
		log.Error("gRPC server is terminated", log.Ferror(err))
	}
}
