package messageservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/anthdm/hollywood/actor"
	"github.com/google/uuid"
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
		cm.conversations[msg.ConversationID] = ctx.SpawnChild(NewConversation([]string{msg.UserID, msg.ReceiverID}, msg.ConversationID), "conversation")
		// db call
		fmt.Println(cm.conversations)
		go func() {
			c := context.Background()
			msgs, err := cm.storage.LoadConversation(c, msg.ConversationID)
			if err != nil {
				//TODO
			} else {
				ctx.Respond(&proto.SuccessOpenAndLoadConversation{Messages: msgs})
			}
		}()
	case *proto.GetUserConversation:
		c := context.Background()
		go func() {
			conversations, err := cm.storage.GetUserConversations(c, msg.Id)
			if err != nil {
				ctx.Respond(&proto.FailureGetUserConversation{})
				fmt.Println(err)
				//TODO
			} else {
				ctx.Respond(&proto.SuccessGetUserConversation{ConvSummary: conversations})
			}
		}()
	case *proto.FillConversationID:
		c := context.Background()
		go func() {
			ok, id, err := cm.storage.DoConversationExists(c, msg.SenderID, msg.ReceiverID)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					utils.Logger.Error("DB ERROR fillconversationDB x1", "ERR:", err)
					ctx.Respond(&proto.FailureOfFillConversationID{})
					return
				}
			}
			if !ok {
				fmt.Println(msg.SenderID, msg.ReceiverID)
				cnv := &proto.CreateConversation{
					Id:         uuid.New().String(),
					SenderID:   msg.SenderID,
					ReceiverID: msg.ReceiverID,
				}
				err = cm.storage.CreateConversation(c, cnv)
				if err != nil {
					ctx.Respond(&proto.FailureOfFillConversationID{})
					utils.Logger.Error("DB ERROR fillconversationDB x2", "ERR", err)
				}
				ctx.Respond(&proto.SuccessOfFillConversationID{Id: cnv.Id})
				return
			}
			ctx.Respond(&proto.SuccessOfFillConversationID{Id: id})
		}()
	case *proto.GetPresence:
		go func() {
			resp := ctx.Request(ctx.Parent(), msg, utils.WaitTime)
			res, _ := resp.Result()
			if message, ok := res.(*proto.Presence); ok {
				ctx.Respond(message)
			} else {
				ctx.Respond(message)
			}
		}()
	case *proto.SendMessage: // if no ok it means that sb is not online in chat
		go func() {
			if _, ok := cm.conversations[msg.Message.ConversationID]; !ok {
				panic("XDDD")
				//or sendign through profile also active this, like click send also actobve openadnload conv but in other message
			}
			orgSender := ctx.Sender()

			//after ctx.Request->result ctx.Sender() is changing to actor that answering on request
			resp := ctx.Request(cm.conversations[msg.Message.ConversationID], msg, utils.WaitTime)
			res, err := resp.Result()
			if err != nil {
				panic(err.Error() + "cnv manager")
			}
			if message, ok := res.(*proto.SuccessSend); ok {
				ctx.Send(orgSender, message)
			} else {
				ctx.Send(orgSender, message)
			}
			fmt.Println("SENDER po", ctx.Sender())
		}()
	case *proto.StoreMessage:
		c := context.Background()
		err := cm.storage.InsertMessage(c, msg.Message)
		if err != nil {
			ctx.Respond(&proto.FailureStoreMessage{})
		} else {
			ctx.Respond(&proto.SuccessStoreMessage{})
		}
	case *proto.DeliverMessage:
		fmt.Println("odebralem od cnv", msg.Message, ctx.Parent())
		ctx.Send(ctx.Parent(), msg)
	default:
		_ = msg
	}
}
