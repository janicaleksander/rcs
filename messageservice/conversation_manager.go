package messageservice

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/utils"
)

type ConversationManager struct {
}

func NewConversationManager() actor.Producer {
	return func() actor.Receiver {
		return &ConversationManager{}
	}
}

func (cm *ConversationManager) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		utils.Logger.Info("Conversation manager has initialized")
	case actor.Started:
		utils.Logger.Info("Conversation manager has started")
	case actor.Stopped:
		utils.Logger.Info("Conversation manager has stooped")

	default:
		_ = msg
	}
}
