package messageservice

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/server"
	"github.com/janicaleksander/bcs/utils"
)

// 1 - > user has presencePlace = 0 -> logged in not in chat
// 2 -> user has presencePlace = 1 -> logged in on chat
// if we don't have uuid of user in map  == 2 point

type MessageService struct {
	serverPID   *actor.PID
	connections map[string]string               // userUUID -> appPID
	presence    map[string]*proto.PresencePlace // PID-> presence
}

func NewMessageService(serverPID *actor.PID) actor.Producer {
	return func() actor.Receiver {
		return &MessageService{
			serverPID:   serverPID,
			connections: make(map[string]string),
			presence:    make(map[string]*proto.PresencePlace),
		}
	}
}

// TODO move logger to utils
func (ms *MessageService) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		server.Logger.Info("messageservice is initialized")
	case actor.Started:
		server.Logger.Info("messageservice is running on:")
	case actor.Stopped:
		server.Logger.Info("messageservice is stopped")
	case *proto.Ping:
		ctx.Respond(&proto.Pong{})
	case *proto.RegisterClient:
		ms.connections[msg.Id] = actor.NewPID(msg.Pid.Address, msg.Pid.Id).String()
	case *proto.UpdatePresence:
		ms.presence[actor.NewPID(msg.Pid.Address, msg.Pid.Id).String()] = msg.PresencePlace
	case *proto.SendMessage:
		ms.sendMessage(ctx, msg)
	case *proto.GetUserConversation:
	default:
		_ = msg
	}
}

func (ms *MessageService) sendMessage(ctx *actor.Context, message *proto.SendMessage) {
	//1 check state of receiver
	senderID := message.Message.SenderID
	senderPID := ms.connections[senderID]
	presencePlace := ms.presence[senderPID]
	switch presencePlace.Place.(type) {
	case *proto.PresencePlace_Outbox:
		fmt.Print("OUTBOX")
	case *proto.PresencePlace_Inbox:
		fmt.Print("INBOX")
	case *proto.PresencePlace_InChat:
		fmt.Print("INCHAT")
		//presencePlace.GetInChat().GetConversationId()

	}
	//2 fill up the conversation ID
	res := ctx.Request(ms.serverPID, &proto.FillConversationID{
		SenderID:   message.Message.SenderID,
		ReceiverID: message.Receiver}, utils.WaitTime)
	resp, err := res.Result()
	if err != nil {
		//TODO
	}
	if v, ok := resp.(*proto.SuccessOfFillConversationID); ok {
		message.Message.ConversationID = v.Id
	} else {
		//TODO ERROR
	}

	//3.send

	//4.DB push
	res = ctx.Request(ms.serverPID, &proto.StoreMessage{Message: message.Message}, utils.WaitTime)
	resp, err = res.Result()
	if err != nil {
		//TODO ERROR
	}

}
