package messageservice

import (
	"context"
	"fmt"
	"github.com/anthdm/hollywood/actor"
	"github.com/janicaleksander/bcs/database"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
	"time"
)

// TODO   we have to stop actor when presence change
// ctx.Engine.Poison()
type Conversation struct {
	id        string
	receivers []string
	storage   database.Storage
}

func NewConversation(receivers []string, id string, storage database.Storage) actor.Producer {
	return func() actor.Receiver {
		return &Conversation{
			receivers: receivers,
			id:        id,
			storage:   storage,
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
		n := time.Now()
		ctxx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := c.storage.InsertMessage(ctxx, msg.Message)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			c.sendMessage(ctx, msg)
			ctx.Respond(&proto.AcceptSend{})
		}

		fmt.Println("w conv", time.Since(n))

	default:
		fmt.Println("XD1")
		_ = msg

	}
}

func (c *Conversation) sendMessage(ctx *actor.Context, msg *proto.SendMessage) {
	n := time.Now()
	// i think i will do this: sender makes a messange in instant its added to []Message on client
	// then is sending throguh ctx request to send to receiver
	for _, receiver := range c.receivers { // cause i cant push directly to message ???
		if receiver != msg.Message.SenderID {

			ctx.Send(ctx.Parent(), &proto.DeliverMessage{
				Receiver: receiver,
				Message:  msg.Message,
			})
		}
	}
	fmt.Println("XD", time.Since(n))
}

//we have to kill this actor somehow
//e.g i change conversation->then i have to change in map that i dont have opened itd
//but it depend what i will choose:
// send from info also add to thisa array->somehow track if user has opened this to not remove this conversation
// do not do this and every change od conv kill prev children
