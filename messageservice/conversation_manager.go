package messageservice

import (
	"context"

	"github.com/anthdm/hollywood/actor"
	db "github.com/janicaleksander/bcs/database"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
)

type ConversationManager struct {
	storage       db.Storage
	conversations map[string]*actor.PID // conversationID -> conversationPID
}

func NewConversationManager(storage db.Storage) actor.Producer {
	return func() actor.Receiver {
		return &ConversationManager{
			storage:       storage,
			conversations: make(map[string]*actor.PID),
		}
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
	case *proto.OpenAndLoadConversation:
		cm.conversations[msg.ConversationID] = ctx.SpawnChild(NewConversation(), "conversation")
		// db call
		go func() {
			c := context.Background()
			msgs, err := cm.storage.LoadConversation(c, msg.ConversationID)
			if err != nil {
				//TODO
			} else {
				ctx.Respond(&proto.SuccessOpenAndLoadConversation{Messages: msgs})
			}
		}()
	default:
		_ = msg
	}
}
