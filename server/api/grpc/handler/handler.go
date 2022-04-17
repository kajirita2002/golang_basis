package handler

import (
	"github.com/google/wire"
	"github.com/kajirita2002/golang_basis/server/api/grpc/pb"
)

// Set is a provider set for handler.
var Set = wire.NewSet(
	NewUser,
	wire.Bind(new(pb.UserServiceServer), new(*User)),
)
