package Unit

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/google/uuid"
	"github.com/janicaleksander/bcs/External"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	"github.com/janicaleksander/bcs/User"
	"reflect"
)

type Unit struct {
	id        uuid.UUID
	serverPID *actor.PID
	external  *External.External
	users     []*User.User
}

func NewUnit(ext *External.External) actor.Producer {
	return func() actor.Receiver {
		return &Unit{
			external: ext,
			users:    make([]*User.User, 1024),
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
	case *Proto.NeededServerConfiguration:
		u.serverPID = actor.NewPID(msg.ServerPID.Address, msg.ServerPID.Id)
	default:
		Server.Logger.Warn("Unit got unknown message", "Type", reflect.TypeOf(msg).String())
		_ = msg

	}

}
