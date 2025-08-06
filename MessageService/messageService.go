package MessageService

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/Server"
)

type MessageService struct {
	serverPID *actor.PID
}

func NewMessageService(serverPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &MessageService{
			serverPID: serverPID,
		}
	}
}

func (ms *MessageService) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		Server.Logger.Info("MessageService is initialized")
	case actor.Started:
		Server.Logger.Info("MessageService is running on:")
	case actor.Stopped:
		Server.Logger.Info("MessageService is stopped")

	default:
		_ = msg
	}
}
