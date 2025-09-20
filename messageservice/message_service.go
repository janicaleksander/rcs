package messageservice

import (
	"reflect"

	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/database"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

// TODO respond is a sugar on the top of send
// when you have another request between sth and respond
// u have to memorize first the orgPID sender from ctx.Sender()
// then make a request for other actor
// and if want to still respond to orgReceiver use orgPID
type MessageService struct {
	storage     database.Storage
	connections map[string]*actor.PID //userUUID -> appPID

	presenceManger     *actor.PID
	conversationManger *actor.PID
}

func NewMessageService(db database.Storage) actor.Producer {
	return func() actor.Receiver {
		return &MessageService{
			storage:     db,
			connections: make(map[string]*actor.PID),
		}
	}
}

// TEST delete go func from actors
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
	case *proto.RegisterClientInMessageService:
		ms.connections[msg.Id] = actor.NewPID(msg.Pid.Address, msg.Pid.Id)
		ctx.Respond(&proto.AcceptRegisterClient{})
		ctx.Send(ms.presenceManger, &proto.UpdatePresence{
			Id: msg.Id,
			Presence: &proto.PresenceType{
				Type: &proto.PresenceType_Outbox{
					Outbox: &proto.Outbox{}},
			},
		})
	case *proto.GetPresence:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, ms.presenceManger, msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})

		} else {
			ctx.Respond(res)
		}
	case *proto.UpdatePresence:
		ctx.Forward(ms.presenceManger)
		//todo idk if i have to make this type check because i never send and if error will be ctx error
	case *proto.OpenAndLoadConversation:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, ms.conversationManger, msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(res)
		}
	case *proto.OpenConversation:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, ms.conversationManger, msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(res)
		}
	case *proto.GetUserConversations:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, ms.conversationManger, msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(res)
		}
	case *proto.FillConversationID:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, ms.conversationManger, msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(res)
		}
	case *proto.SendMessage:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, ms.conversationManger, msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(res)
		}
	case *proto.DeliverMessage:
		ctx.Send(ms.connections[msg.Receiver], msg)
	case *proto.GetUsersToNewConversation:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, ms.conversationManger, msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(res)
		}
	case *proto.CreateConversation:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, ms.conversationManger, msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(res)
		}
	default:
		utils.Logger.Info("Unsupported type of message", reflect.TypeOf(msg).String())

	}
}
