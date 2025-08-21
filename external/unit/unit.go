package unit

import (
	"reflect"

	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/server"
	proto2 "github.com/janicaleksander/bcs/types/proto"
)

type Unit struct {
	id        string // uuid
	serverPID *actor.PID
	devices   []*proto.Device
	users     []*proto2.User
}

func NewUnit(serverPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &Unit{
			serverPID: serverPID,
			users:     make([]*proto2.User, 1024),
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
