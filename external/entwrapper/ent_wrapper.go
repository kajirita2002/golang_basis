package entwrapper

import (
	"github.com/google/wire"
	"github.com/kajirita2002/golang_basis/domain/repository"
)

var Set = wire.NewSet(
	NewUser,
	wire.Bind(new(repository.IFUser), new(*UserImpl)),
)
