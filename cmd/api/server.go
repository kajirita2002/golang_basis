package main

import (
	"google.golang.org/grpc"

	agrpc "github.com/kajirita2002/golang_basis/server/api/grpc"
	"github.com/kajirita2002/golang_basis/server/api/grpc/pb"
)

type ServiceServer struct {
	user pb.UserServiceServer
}

func newServer(user pb.UserServiceServer) *ServiceServer {
	return &ServiceServer{
		user: user,
	}
}

func registerServicesToGRPCServer(ss *ServiceServer) agrpc.Register {
	return func(gs *grpc.Server) {
		pb.RegisterUserServiceServer(gs, ss.user)
	}
}
