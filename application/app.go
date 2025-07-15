package application

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	"reflect"
	"time"
)

func NewApp(login chan *Proto.LoginUser) actor.Producer {
	return func() actor.Receiver {
		return &App{
			ChLoginUser: login,
		}
	}
}

type App struct {
	//actors
	serverPID *actor.PID

	//chan's
	ChLoginUser chan *Proto.LoginUser
}

func (a *App) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
	case actor.Initialized:
		//this is loop to manage chann's communication
		go a.handleEventStream(ctx)
	case actor.Stopped:
	case *Proto.NeededServerConfiguration:
		a.serverPID = actor.NewPID(msg.ServerPID.Address, msg.ServerPID.Id)
	default:
		Server.Logger.Warn("Server got unknown message", "Type:", reflect.TypeOf(msg).String())
	}
}

func (a *App) handleEventStream(ctx *actor.Context) {
	for {
		select {
		case msg := <-a.ChLoginUser:
			msg.Pid = &Proto.PID{
				Address: ctx.PID().GetAddress(),
				Id:      ctx.PID().GetID(),
			}
			res := ctx.Request(a.serverPID, msg, 3*time.Second)
			v, err := res.Result()
			if err != nil {
				_ = v
				// sth
			}

		}
	}
}
