package messageservice

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

type PresencePlace uint8

/*
const (

	OUTBOX  = 0
	INBOX   = 1

)//trzeba jakos zrobic zeby to by wiadomo jaki inbox dla jakiego user
*/
type PresenceManager struct {
	presence map[string]*proto.PresenceType
}

func NewPresenceManager() actor.Producer {
	return func() actor.Receiver {
		return &PresenceManager{
			presence: make(map[string]*proto.PresenceType),
		}
	}
}

// TODO Presence will be change: when sb is entering a conversation
// there will be button home to in inbox area exit all chats and then move to outbox
// or faster clicked on go back button in inbox (also outbox change)
func (pm *PresenceManager) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		utils.Logger.Info("Presence manager has initialized")
	case actor.Started:
		utils.Logger.Info("Presence manager has started")
	case actor.Stopped:
		utils.Logger.Info("Presence manager has stooped")
	case *proto.UpdatePresence:
		pm.presence[msg.Id] = msg.Presence
	case *proto.GetPresence:
		presence, ok := pm.presence[msg.Id]
		if !ok {
			ctx.Respond(&proto.Presence{
				Presence: &proto.PresenceType{
					Type: &proto.PresenceType_Outbox{Outbox: &proto.Outbox{}},
				},
			})
		} else {
			ctx.Respond(&proto.Presence{
				Presence: presence,
			})
		}

	default:
		_ = msg
	}

}
