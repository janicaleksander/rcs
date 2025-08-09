package MessageService

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	"github.com/janicaleksander/bcs/Utils"
)

// 1 - > user has presencePlace = 0 -> logged in not in chat
// 2 -> user has presencePlace = 1 -> logged in on chat
// if we don't have uuid of user in map  == 2 point

type MessageService struct {
	serverPID   *actor.PID
	connections map[string]string               // userUUID -> appPID
	presence    map[string]*Proto.PresencePlace // PID-> presence
}

func NewMessageService(serverPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &MessageService{
			serverPID:   serverPID,
			connections: make(map[string]string),
			presence:    make(map[string]*Proto.PresencePlace),
		}
	}
}

// TODO move logger to utils
func (ms *MessageService) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		Server.Logger.Info("MessageService is initialized")
	case actor.Started:
		Server.Logger.Info("MessageService is running on:")
	case actor.Stopped:
		Server.Logger.Info("MessageService is stopped")
	case *Proto.Ping:
		ctx.Respond(&Proto.Pong{})
	case *Proto.RegisterClient:
		ms.connections[msg.Id] = actor.NewPID(msg.Pid.Address, msg.Pid.Id).String()
	case *Proto.UpdatePresence:
		ms.presence[actor.NewPID(msg.Pid.Address, msg.Pid.Id).String()] = msg.PresencePlace
	case *Proto.SendMessage:
		ms.sendMessage(ctx, msg)
	case *Proto.GetUserConversation:
	default:
		_ = msg
	}
}

func (ms *MessageService) sendMessage(ctx *actor.Context, message *Proto.SendMessage) {
	//1 check state of receiver
	senderID := message.Message.SenderID
	senderPID := ms.connections[senderID]
	presencePlace := ms.presence[senderPID]
	switch presencePlace.Place.(type) {
	case *Proto.PresencePlace_Outbox:
		fmt.Print("OUTBOX")
	case *Proto.PresencePlace_Inbox:
		fmt.Print("INBOX")
	case *Proto.PresencePlace_InChat:
		fmt.Print("INCHAT")
		//presencePlace.GetInChat().GetConversationId()

	}
	//2 fill up the conversation ID
	res := ctx.Request(ms.serverPID, &Proto.FillConversationID{
		SenderID:   message.Message.SenderID,
		ReceiverID: message.Receiver}, Utils.WaitTime)
	resp, err := res.Result()
	if err != nil {
		//TODO
	}
	if v, ok := resp.(*Proto.SuccessOfFillConversationID); ok {
		message.Message.ConversationID = v.Id
	} else {
		//TODO ERROR
	}

	//3.send

	//4.DB push
	res = ctx.Request(ms.serverPID, &Proto.StoreMessage{Message: message.Message}, Utils.WaitTime)
	resp, err = res.Result()
	if err != nil {
		//TODO ERROR
	}

}
