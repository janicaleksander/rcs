package messageservice

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/database"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
)

type MessageService struct {
	storage     database.Storage
	serverPID   *actor.PID            //TODO maybe rmv this? //
	connections map[string]*actor.PID //userUUID -> appPID

	presenceManger     *actor.PID
	conversationManger *actor.PID
}

func NewMessageService(serverPID *actor.PID, db database.Storage) actor.Producer {
	return func() actor.Receiver {
		return &MessageService{
			storage:     db,
			serverPID:   serverPID,
			connections: make(map[string]*actor.PID),
		}
	}
}

// TODO move logger to utils
func (ms *MessageService) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		utils.Logger.Info("messageservice is initialized")
	case actor.Started:
		utils.Logger.Info("messageservice is running on:")
		ms.presenceManger = ctx.SpawnChild(NewPresenceManager(), "presence_manager")
		ms.conversationManger = ctx.SpawnChild(NewConversationManager(ms.storage), "conversation_manager")
	case actor.Stopped:
		utils.Logger.Info("messageservice is stopped")
	case *proto.Ping:
		ctx.Respond(&proto.Pong{})
	case *proto.RegisterClient:
		ms.connections[msg.Id] = actor.NewPID(msg.Pid.Address, msg.Pid.Id)
		ctx.Send(ms.presenceManger, &proto.UpdatePresence{
			Id: msg.Id,
			Presence: &proto.PresenceType{
				Type: &proto.PresenceType_Outbox{
					Outbox: &proto.Outbox{}},
			},
		})
	case *proto.UpdatePresence:
		ctx.Forward(ms.presenceManger)
	case *proto.OpenAndLoadConversation:
		go func() {
			resp := ctx.Request(ms.conversationManger, msg, utils.WaitTime)
			res, _ := resp.Result()
			if conversation, ok := res.(*proto.SuccessOpenAndLoadConversation); ok {
				ctx.Respond(conversation)
			}
		}()
	default:
		_ = msg
	}
}
