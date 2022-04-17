package main

import (
	"github.com/kajirita2002/golang_basis/config"
	"github.com/kajirita2002/golang_basis/external/ent"
	"github.com/kajirita2002/golang_basis/log"
)

func newEntClient(cfg config.Ent) *ent.Client {
	cli, err := ent.Open(cfg.DriverName, cfg.DataSourceName)
	if err != nil {
		log.Error("failed to open mysql: %v", log.Ferror(err))
	}
	return cli
}
