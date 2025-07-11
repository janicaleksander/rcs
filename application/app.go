package application

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	"reflect"
)

func NewApp() actor.Producer {
	return func() actor.Receiver {
		return &App{}
	}
}

type App struct {
	//actors
	serverPID *actor.PID
}

func (a *App) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
	case actor.Initialized:
	case actor.Stopped:
	case *Proto.NeededServerConfiguration:
		a.serverPID = actor.NewPID(msg.ServerPID.Address, msg.ServerPID.Id)
	default:
		Server.Logger.Warn("Server got unknown message", "Type:", reflect.TypeOf(msg).String())
	}
}
