package messageservice

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/utils"
)

type Conversation struct {
}

func NewConversation() actor.Producer {
	return func() actor.Receiver {
		return &Conversation{}
	}
}

func (c *Conversation) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		utils.Logger.Info("Conversation has initialized")
	case actor.Started:
		utils.Logger.Info("Conversation has started")
	case actor.Stopped:
		utils.Logger.Info("Conversation has stooped")
	default:
		_ = msg
	}
}
