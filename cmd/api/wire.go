//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/kajirita2002/golang_basis/config"
	"github.com/kajirita2002/golang_basis/domain/service"
	ent_wrapper "github.com/kajirita2002/golang_basis/external/entwrapper"
	"github.com/kajirita2002/golang_basis/server/api/grpc/handler"
	"github.com/kajirita2002/golang_basis/usecase"
)

func injector() (*ServiceServer, error) {
	wire.Build(
		config.NewEnt,
		handler.Set,
		usecase.Set,
		service.Set,
		ent_wrapper.Set,
		newEntClient,
		newServer,
	)
	return nil, nil
}
