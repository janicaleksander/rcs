package cli_app

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	"reflect"
)

type Application struct{}

type CLI struct {
	serverPID   *actor.PID
	application Application
}

func NewCLI() actor.Producer {
	return func() actor.Receiver {
		return &CLI{}
	}
}

func (c *CLI) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
	case actor.Initialized:
	case actor.Stopped:
	case *Proto.NeededServerConfiguration:
		c.serverPID = actor.NewPID(msg.ServerPID.Address, msg.ServerPID.Id)
	case *Proto.StartCLI:

	default:
		Server.Logger.Warn("Server got unknown message", "Type:", reflect.TypeOf(msg).String())
	}
}
