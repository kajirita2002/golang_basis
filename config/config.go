package config

import (
	"github.com/kelseyhightower/envconfig"

	"github.com/kajirita2002/golang_basis/external/ent"
)

const (
	DEV  = "development"
	PRD  = "production"
	LOAD = "load"
)

type Service struct {
	Env string `envconfig:"ENV" default:"test"`
}

type Bugsnag struct {
	APIKey string `envconfig:"BUGSNAG_API_KEY" default:"api-key"`
}

type Ent struct {
	DriverName     string `envconfig:"DRIVER_NAME" defalut:"mysql"`
	DataSourceName string `envconfig:"DATA_SOURCE_NAME" default:""`
}

func NewService() (Service, error) {
	ent.NewClient()
	conf := Service{}
	if err := envconfig.Process("BASIS", &conf); err != nil {
		return Service{}, err
	}

	return conf, nil
}

func NewBugsnag() (Bugsnag, error) {
	conf := Bugsnag{}
	if err := envconfig.Process("BASIS", &conf); err != nil {
		return Bugsnag{}, err
	}
	return conf, nil
}

func NewEnt() (Ent, error) {
	conf := Ent{}
	if err := envconfig.Process("BASIS", &conf); err != nil {
		return Ent{}, err
	}
	return conf, nil
}
