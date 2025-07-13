package application

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	"reflect"
)

func NewApp(test chan *Proto.Payload) actor.Producer {
	return func() actor.Receiver {
		return &App{
			test: test,
		}
	}
}

type App struct {
	//actors
	serverPID *actor.PID

	//chan's
	test chan *Proto.Payload
}

func (a *App) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
	case actor.Initialized:
	case actor.Stopped:
	case *Proto.NeededServerConfiguration:
		a.serverPID = actor.NewPID(msg.ServerPID.Address, msg.ServerPID.Id)
	case *Proto.Payload:
		a.test <- &Proto.Payload{Data: msg.Data}
	default:
		Server.Logger.Warn("Server got unknown message", "Type:", reflect.TypeOf(msg).String())
	}
}
