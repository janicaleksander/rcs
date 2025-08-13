package messageservice

import (
	"fmt"
	"time"

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
		ctx.Respond(&proto.SuccessRegisterClient{})
		ctx.Send(ms.presenceManger, &proto.UpdatePresence{
			Id: msg.Id,
			Presence: &proto.PresenceType{
				Type: &proto.PresenceType_Outbox{
					Outbox: &proto.Outbox{}},
			},
		})
	case *proto.GetPresence:

		go func() {
			resp := ctx.Request(ms.presenceManger, msg, utils.WaitTime)
			res, _ := resp.Result()
			if message, ok := res.(*proto.Presence); ok {
				ctx.Respond(message)
			} else {
				ctx.Respond(message)
			}
		}()
	case *proto.UpdatePresence:
		ctx.Forward(ms.presenceManger)
		//todo idk if i have to make this type check because i never send and if error will be ctx error
	case *proto.OpenAndLoadConversation:
		go func() {
			resp := ctx.Request(ms.conversationManger, msg, utils.WaitTime)
			res, _ := resp.Result()
			if message, ok := res.(*proto.SuccessOpenAndLoadConversation); ok {
				ctx.Respond(message)
			} else {
				utils.Logger.Error("Error in MSSVC OpenAndLoadConversation")
				ctx.Respond(message)
			}

		}()
	case *proto.GetUserConversation:
		go func() {
			resp := ctx.Request(ms.conversationManger, msg, utils.WaitTime)
			res, _ := resp.Result()
			if message, ok := res.(*proto.SuccessGetUserConversation); ok {
				ctx.Respond(message)
			} else {
				utils.Logger.Error("Error in MSSVC GetUserConversation")
				ctx.Respond(message)
			}
		}()
	case *proto.FillConversationID:
		go func() {
			resp := ctx.Request(ms.conversationManger, msg, utils.WaitTime)
			res, _ := resp.Result()
			if message, ok := res.(*proto.SuccessOfFillConversationID); ok {
				ctx.Respond(message)
			} else {
				utils.Logger.Error("Error in MSSVC FillConversationID")
				ctx.Respond(message)
			}
		}()
	case *proto.SendMessage:
		n := time.Now()
		go func() {
			org := ctx.Sender()
			resp := ctx.Request(ms.conversationManger, msg, utils.WaitTime)
			fmt.Println("PID pierwszego requesta", resp.PID())
			res, err := resp.Result()
			if message, ok := res.(*proto.SuccessSend); ok {
				ctx.Send(org, message)
			} else {
				utils.Logger.Error("Error in MSSVC SendMessage", err)
				ctx.Send(org, message)
			}
			fmt.Println("POSZŁO po ", time.Since(n))
		}()
	case *proto.DeliverMessage:
		fmt.Println("odebralem od cnv manager", msg)
		ctx.Send(ms.connections[msg.Receiver], msg)
		fmt.Println(ms.connections[msg.Receiver])
	default:

		_ = msg
	}
}

/*
case *proto.FillConversationID:
c := context.Background()
ok, id, err := s.storage.IsConversationExists(c, msg.SenderID, msg.ReceiverID)
if err != nil || !ok {
cnv := &proto.CreateConversationAndAssign{
Id:         uuid.New().String(),
SenderID:   msg.SenderID,
ReceiverID: msg.ReceiverID,
}
err = s.storage.CreateAndAssignConversation(c, cnv)
if err != nil {
//TODO ERROR
} else {
ctx.Respond(&proto.SuccessOfFillConversationID{Id: cnv.Id})
}

} else {
ctx.Respond(&proto.SuccessOfFillConversationID{Id: id})
}
case *proto.StoreMessage:
c := context.Background()
err := s.storage.InsertMessage(c, msg.Message)
if err != nil {
ctx.Respond(&proto.FailureStoreMessage{})
} else {
ctx.Respond(&proto.SuccessStoreMessage{})
}


*/
