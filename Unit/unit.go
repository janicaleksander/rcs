package Unit

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/External"
	"github.com/janicaleksander/bcs/Server"
	"reflect"
)

type Unit struct {
	serverPID *actor.PID
	external  *External.External
}

func NewUnit(ext *External.External) actor.Producer {
	return func() actor.Receiver {
		return &Unit{
			external: ext,
		}
	}
}

// actor with remote
func (u *Unit) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		Server.Logger.Info("Unit has initialized")
	case actor.Started:
		Server.Logger.Info("Unit has started")
	case actor.Stopped:
		Server.Logger.Info("Unit has stopped")
	default:
		Server.Logger.Warn("Unit got unknown message", "Type", reflect.TypeOf(msg).String())
		_ = msg

	}

}
