package messageservice

import (
	"fmt"
	"time"

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
		resp := ctx.Request(ms.presenceManger, msg, utils.WaitTime)
		res, _ := resp.Result()
		if message, ok := res.(*proto.Presence); ok {
			ctx.Respond(message)
		} else {
			ctx.Respond(message)
		}
	case *proto.UpdatePresence:
		ctx.Forward(ms.presenceManger)
		//todo idk if i have to make this type check because i never send and if error will be ctx error
	case *proto.OpenAndLoadConversation:
		resp := ctx.Request(ms.conversationManger, msg, utils.WaitTime)
		res, _ := resp.Result()
		if message, ok := res.(*proto.LoadedConversation); ok {
			ctx.Respond(message)
		} else {
			utils.Logger.Error("Error in MSSVC OpenAndLoadConversation")
			ctx.Respond(message)
		}

	case *proto.GetUserConversations:
		resp := ctx.Request(ms.conversationManger, msg, utils.WaitTime)
		res, _ := resp.Result()
		if message, ok := res.(*proto.UserConversations); ok {
			ctx.Respond(message)
		} else {
			utils.Logger.Error("Error in MSSVC GetUserConversation")
			ctx.Respond(message)
		}
	case *proto.FillConversationID:
		resp := ctx.Request(ms.conversationManger, msg, utils.WaitTime)
		res, _ := resp.Result()
		if message, ok := res.(*proto.FilledConversationID); ok {
			ctx.Respond(message)
		} else {
			utils.Logger.Error("Error in MSSVC FillConversationID")
			ctx.Respond(message)
		}
	case *proto.SendMessage:
		n := time.Now()
		org := ctx.Sender()
		resp := ctx.Request(ms.conversationManger, msg, utils.WaitTime)
		fmt.Println("PID pierwszego requesta", resp.PID())
		res, err := resp.Result()
		if message, ok := res.(*proto.AcceptSend); ok {
			ctx.Send(org, message)
		} else {
			utils.Logger.Error("Error in MSSVC SendMessage", err)
			ctx.Send(org, message)
		}
		fmt.Println("POSZŁO po ", time.Since(n))
	case *proto.DeliverMessage:
		fmt.Println("odebralem od cnv manager", msg)
		ctx.Send(ms.connections[msg.Receiver], msg)
		fmt.Println(ms.connections[msg.Receiver])
	case *proto.GetUsersToNewConversation:
		resp := ctx.Request(ms.conversationManger, msg, utils.WaitTime)
		res, err := resp.Result()
		if err != nil {
			panic(err)
		}
		if message, ok := res.(*proto.UsersToNewConversation); ok {
			ctx.Respond(message)
		} else {
			utils.Logger.Error("Error in MSSVC GetUsersToNewConversation")
			ctx.Respond(message)
		}
	case *proto.CreateConversation:
		resp := ctx.Request(ms.conversationManger, msg, utils.WaitTime)
		res, err := resp.Result()
		if err != nil {
			utils.Logger.Error("Here x4")
			ctx.Respond(&proto.Error{Content: err.Error()})
			return
		}
		if _, ok := res.(*proto.AcceptCreateConversation); ok {
			ctx.Respond(&proto.AcceptCreateConversation{})
			utils.Logger.Error("Here x5")
			return
		} else {
			utils.Logger.Error("Here x6")
			ctx.Respond(&proto.Error{Content: err.Error()})
			return
		}

	default:

		_ = msg
	}
}
