package notifier

import (
	"github.com/bugsnag/bugsnag-go"

	"github.com/kajirita2002/golang_basis/config"
)

type BugNotifier interface {
	Notify(error, ...interface{}) error
}

func NewBugsnagNotifier(srvCfg config.Service, bugsnagCfg config.Bugsnag) *bugsnag.Notifier {
	return bugsnag.New(bugsnag.Configuration{
		APIKey:          bugsnagCfg.APIKey,
		ReleaseStage:    srvCfg.Env,
		ProjectPackages: []string{"main", "github.com/kajirita2002/golang_basis/**"},
		PanicHandler:    func() {},
	})
}
