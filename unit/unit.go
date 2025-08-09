package unit

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/external"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/server"
	"reflect"
)

type Unit struct {
	id        string // uuid
	serverPID *actor.PID
	external  *external.External
	users     []*proto.User
}

func NewUnit(serverPID *actor.PID, ext *external.External) actor.Producer {
	return func() actor.Receiver {
		return &Unit{
			serverPID: serverPID,
			external:  ext,
			users:     make([]*proto.User, 1024),
		}
	}
}

// actor with remote
func (u *Unit) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		server.Logger.Info("Unit has initialized")
	case actor.Started:
		server.Logger.Info("Unit has started")
	case actor.Stopped:
		server.Logger.Info("Unit has stopped")
	default:
		server.Logger.Warn("Unit got unknown message", "Type", reflect.TypeOf(msg).String())
		_ = msg

	}

}
