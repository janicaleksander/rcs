package messageservice

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/google/uuid"
	db "github.com/janicaleksander/bcs/database"
	"github.com/janicaleksander/bcs/types/proto"
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
	case *proto.CreateConversation: // ?
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		exists, _, err := cm.storage.DoConversationExists(c, msg.SenderID, msg.ReceiverID)
		if exists || (err != nil && !errors.Is(err, sql.ErrNoRows)) {
			ctx.Respond(&proto.Error{Content: err.Error()})
			return
		}
		cnv := &proto.Conversation{
			Id:         msg.Id,
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
		}
		err = cm.storage.CreateConversation(c, cnv)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
			return
		} else {
			ctx.Respond(&proto.AcceptCreateConversation{})
			return
		}
	case *proto.OpenConversation:
		if _, ok := cm.conversations[msg.ConversationID]; !ok {
			cm.conversations[msg.ConversationID] = ctx.SpawnChild(NewConversation([]string{msg.UserID, msg.ReceiverID}, msg.ConversationID, cm.storage), "conversation")
		}
		ctx.Respond(&proto.OpenedConversation{})
	case *proto.OpenAndLoadConversation:
		if _, ok := cm.conversations[msg.ConversationID]; !ok {
			cm.conversations[msg.ConversationID] = ctx.SpawnChild(NewConversation([]string{msg.UserID, msg.ReceiverID}, msg.ConversationID, cm.storage), "conversation")
		}
		// db call
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		msgs, err := cm.storage.LoadConversation(c, msg.ConversationID)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.LoadedConversation{Messages: msgs})
		}
	case *proto.GetUserConversations:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		conversations, err := cm.storage.GetUserConversations(c, msg.Id)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.UserConversations{ConvSummary: conversations})
		}
	case *proto.FillConversationID:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ok, id, err := cm.storage.DoConversationExists(c, msg.SenderID, msg.ReceiverID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				ctx.Respond(&proto.Error{Content: err.Error()})
				return
			}
		}
		if ok {
			ctx.Respond(&proto.FilledConversationID{Id: id})
		} else {
			cnv := &proto.Conversation{
				Id:         uuid.New().String(),
				SenderID:   msg.SenderID,
				ReceiverID: msg.ReceiverID,
			}
			err = cm.storage.CreateConversation(c, cnv)
			if err != nil {
				ctx.Respond(&proto.Error{Content: err.Error()})
			} else {
				ctx.Respond(&proto.FilledConversationID{Id: cnv.Id})
			}
		}
	case *proto.GetPresence:
		res, err := utils.MakeRequest(utils.NewRequest(ctx, ctx.Parent(), msg))
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(res)
		}
	case *proto.SendMessage: // if no ok it means that sb is not online in chat ?? TODO
		orgSender := ctx.Sender()
		if _, ok := cm.conversations[msg.Message.ConversationID]; !ok {
			panic("XDDD")
			//or sendign through profile also active this, like click send also actobve openadnload conv but in other message
		}
		//after ctx.Request->result ctx.Sender() is changing to actor that answering on request
		res, err := utils.MakeRequest(utils.NewRequest(ctx, cm.conversations[msg.Message.ConversationID], msg))
		if err != nil {
			ctx.Send(orgSender, &proto.Error{Content: err.Error()})
		} else {
			ctx.Send(orgSender, res)
		}

	case *proto.DeliverMessage:
		ctx.Send(ctx.Parent(), msg)
	case *proto.GetUsersToNewConversation:
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		users, err := cm.storage.SelectUsersToNewConversation(c, msg.UserID)
		if err != nil {
			ctx.Respond(&proto.Error{Content: err.Error()})
		} else {
			ctx.Respond(&proto.UsersToNewConversation{Users: users})
		}
	default:
		utils.Logger.Info("Unsupported type of message", reflect.TypeOf(msg).String())
	}
}
