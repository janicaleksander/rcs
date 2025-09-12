package messageservice

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

// TODO   we have to stop actor when presence change
// ctx.Engine.Poison()
type Conversation struct {
	id        string
	receivers []string
}

func NewConversation(receivers []string, id string) actor.Producer {
	return func() actor.Receiver {
		return &Conversation{
			receivers: receivers,
			id:        id,
		}
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
	case *proto.SendMessage:
		//if store success then ->
		resp := ctx.Request(ctx.Parent(), &proto.StoreMessage{Message: msg.Message}, utils.WaitTime)
		res, err := resp.Result()
		if err != nil {
			panic(err.Error() + "cnv")
		}
		if _, ok := res.(*proto.AcceptStoreMessage); ok {
			ctx.Respond(&proto.AcceptSend{})
		} else {
			ctx.Respond(&proto.Error{Content: err.Error()})
			utils.Logger.Error("SOME ERROR in sending ")
			return
		}
		c.sendMessage(ctx, msg)
	default:
		_ = msg

	}
}

func (c *Conversation) sendMessage(ctx *actor.Context, msg *proto.SendMessage) {
	// i think i will do this: sender makes a messange in instant its added to []Message on client
	// then is sending throguh ctx request to send to receiver
	for _, receiver := range c.receivers { // cause i cant push directly to message ???
		if true || receiver != msg.Message.SenderID {
			resp := ctx.Request(ctx.Parent(), &proto.GetPresence{Id: receiver}, utils.WaitTime)
			//TODO do this _ err
			res, err := resp.Result()
			if err != nil {
				panic("error conversation sendMessage")
			}
			if message, ok := res.(*proto.Presence); ok {
				switch message.Presence.Type.(type) {
				case *proto.PresenceType_Outbox:
					//only to db in receive loop
				case *proto.PresenceType_Inbox:
					//send
					ctx.Send(ctx.Parent(), &proto.DeliverMessage{
						Receiver: receiver,
						Message:  msg.Message,
					})
					fmt.Println("WYSYLAM", msg.Message, "do", ctx.Parent())
				default:
					utils.Logger.Error("Brak ustawionego typu")
				}
			}
		}
	}
}

//we have to kill this actor somehow
//e.g i change conversation->then i have to change in map that i dont have opened itd
//but it depend what i will choose:
// send from info also add to thisa array->somehow track if user has opened this to not remove this conversation
// do not do this and every change od conv kill prev children
