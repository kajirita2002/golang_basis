package usecase

import "github.com/google/wire"

var Set = wire.NewSet(
	NewUser,
	wire.Bind(new(IFUser), new(*UserImpl)),
)
