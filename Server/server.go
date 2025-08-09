package Server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"time"

	"github.com/anthdm/hollywood/actor"
	"github.com/google/uuid"
	db "github.com/janicaleksander/bcs/Database"
	"github.com/janicaleksander/bcs/Proto"
)

// GENERAL TODO check why in some places when i have messageservcie down loading is too long
var Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

const (
	PingPongTime = 3 * time.Second
)

// TODO make a two maybe maps one for connected by app commander lvl person
// and second to 0 1 2 users/soldiers connected by device->unit->server
// this is problem e.g. in sending message from 5lvl from app to 0lvl to device in unit
// but this 5 see this 0 in his lists of person in system (maybe)
type Server struct {
	storage     db.Storage
	listenAddr  string            // IP of (one) main server
	connections map[string]string // PID  to uuid
	active      map[string]bool   //uuid->bool
	activeChan  chan struct {
		uuid string
		PID  *actor.PID
	}
}

func NewServer(listenAddr string, storage db.Storage) actor.Producer {
	return func() actor.Receiver {
		return &Server{
			storage:     storage,
			listenAddr:  listenAddr,
			connections: make(map[string]string),
			active:      make(map[string]bool),
			activeChan: make(chan struct {
				uuid string
				PID  *actor.PID
			}, 1024),
		}
	}
}

// TODO In the future legal way to disconnect
// TODO manage if err !=nil some log errors
func (s *Server) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		Logger.Info("Server initialized")

	case actor.Started:
		Logger.Info("Server has started")
		go s.betterHeartbeat(ctx)
	case actor.Stopped:
		close(s.activeChan)
		Logger.Info("Server has stopped")
	case *Proto.IsServerRunning:
		ctx.Respond(&Proto.Running{})
	case *Proto.Disconnect: // after this switch state to loginScene
		var pidToDelete *actor.PID
		if msg.Pid != nil {
			pidToDelete = actor.NewPID(msg.Pid.Address, msg.Pid.Id)
		} else {
			pidToDelete = ctx.Sender()
		}
		id, ok := s.connections[pidToDelete.String()]
		if ok {
			delete(s.active, id)
			delete(s.connections, pidToDelete.String())
		}
	case *Proto.IsOnline:
		if s.active[msg.Uuid] {
			ctx.Respond(&Proto.Online{})
		} else {
			ctx.Respond(&Proto.Offline{})
		}
	case *Proto.LoginUnit:
		//update use map
		//id, err := s.loginUnit(ctx, msg.Email, msg.Password)
		//if err == nil {
		//	s.clients[id] = ctx.Sender().GetAddress()
		//}
	case *Proto.LoginUser:
		id, role, err := s.loginUser(msg.Email, msg.Password)
		if err != nil {
			ctx.Respond(&Proto.DenyLogin{Info: err.Error()})
		} else {
			pid := actor.NewPID(msg.Pid.Address, msg.Pid.Id) //client PID
			s.connections[pid.String()] = id                 //pid to uuid
			fmt.Println(s.connections)
			s.activeChan <- struct {
				uuid string
				PID  *actor.PID
			}{uuid: id, PID: pid}
			ctx.Respond(&Proto.AcceptLogin{Info: "Login successful! ", RuleLevel: int64(role)})

		}
		//TODO idk if this getlogged works

	case *Proto.GetLoggedInUUID: //returning an uuid of current logged-in user
		pid := actor.NewPID(msg.Pid.Address, msg.Pid.Id) //client PID
		id := s.connections[pid.String()]
		ctx.Respond(&Proto.LoggedInUUID{Id: id})
	case *Proto.GetUserAboveLVL:
		c := context.Background()
		users, err := s.storage.GetUsersWithLVL(c, int(msg.Lvl))
		if err == nil {
			fmt.Println(users)
			ctx.Respond(&Proto.UsersAboveLVL{Users: users})
		}
	case *Proto.CreateUnit:
		c := context.Background()
		err := s.storage.InsertUnit(c, msg.Name, msg.IsConfigured, msg.UserID)
		if err != nil {
			ctx.Respond(&Proto.DenyCreateUnit{Info: err.Error()})
		} else {
			ctx.Respond(&Proto.AcceptCreateUnit{})
		}
	case *Proto.GetAllUnits:
		c := context.Background()
		units, err := s.storage.GetAllUnits(c)
		if err == nil {
			ctx.Respond(&Proto.AllUnits{Units: units})
		}

	case *Proto.GetAllUsersInUnit:
		c := context.Background()
		unitID := msg.Id
		users, err := s.storage.GetUsersInUnit(c, unitID)
		if err == nil {
			fmt.Println(users)
			ctx.Respond(&Proto.AllUsersInUnit{Users: users})
		}
	case *Proto.CreateUser:
		c := context.Background()
		err := s.storage.InsertUser(c, msg.User)
		if err != nil {
			ctx.Respond(&Proto.DenyCreateUser{Info: err.Error()})
		} else {
			ctx.Respond(&Proto.AcceptCreateUser{})
		}
	case *Proto.IsUserInUnit:
		c := context.Background()
		isInUnit, unitID, err := s.storage.IsUserInUnit(c, msg.Id)
		if err != nil {
			ctx.Respond(&Proto.UserNotInUnit{})
		} else if isInUnit {
			ctx.Respond(&Proto.UserInUnit{UnitID: unitID})
		} else {
			ctx.Respond(&Proto.UserNotInUnit{})
		}
	case *Proto.AssignUserToUnit:
		c := context.Background()
		err := s.storage.AssignUserToUnit(c, msg.UserID, msg.UnitID)
		if err != nil {
			ctx.Respond(&Proto.FailureOfAssign{})
		} else {
			ctx.Respond(&Proto.SuccessOfAssign{})
		}
	case *Proto.DeleteUserFromUnit:
		c := context.Background()
		err := s.storage.DeleteUserFromUnit(c, msg.UserID, msg.UnitID)
		if err != nil {
			ctx.Respond(&Proto.FailureOfDelete{})
		} else {
			ctx.Respond(&Proto.SuccessOfDelete{})
		}
	case *Proto.FillConversationID:
		c := context.Background()
		ok, id, err := s.storage.IsConversationExists(c, msg.SenderID, msg.ReceiverID)
		if err != nil || !ok {
			cnv := &Proto.CreateConversationAndAssign{
				Id:         uuid.New().String(),
				SenderID:   msg.SenderID,
				ReceiverID: msg.ReceiverID,
			}
			err = s.storage.CreateAndAssignConversation(c, cnv)
			if err != nil {
				//TODO ERROR
			} else {
				ctx.Respond(&Proto.SuccessOfFillConversationID{Id: cnv.Id})
			}

		} else {
			ctx.Respond(&Proto.SuccessOfFillConversationID{Id: id})
		}
	case *Proto.StoreMessage:
		c := context.Background()
		err := s.storage.InsertMessage(c, msg.Message)
		if err != nil {
			ctx.Respond(&Proto.FailureStoreMessage{})
		} else {
			ctx.Respond(&Proto.SuccessStoreMessage{})
		}
	case *Proto.GetUserConversation:
		c := context.Background()
		conversations, err := s.storage.GetUserConversations(c, msg.Id)
		if err != nil {
			ctx.Respond(&Proto.FailureGetUserConversation{})
			fmt.Println(err)
			//TODO
		} else {
			ctx.Respond(&Proto.SuccessGetUserConversation{ConvSummary: conversations})
		}
	default:
		Logger.Warn("Server got unknown message", reflect.TypeOf(msg).String())

	}
}

// if sth is not in active map => not active
func (s *Server) betterHeartbeat(ctx *actor.Context) {
	for value := range s.activeChan {
		//value is uuid or PID
		s.active[value.uuid] = true
		go func(pid *actor.PID) {
			for {
				resp := ctx.Request(pid, &Proto.Ping{}, PingPongTime)
				res, err := resp.Result()
				_, ok := res.(*Proto.Pong)
				if !ok || err != nil {
					ctx.Send(ctx.PID(), &Proto.Disconnect{Pid: &Proto.PID{
						Address: value.PID.Address,
						Id:      value.PID.ID,
					}})
					break
				}
			}
		}(value.PID)
	}
}

func (s *Server) loginUser(email, password string) (string, int, error) {
	// TODO: add jwt to login
	c := context.Background()
	id, role, err := s.storage.LoginUser(c, email, password)
	if err != nil {
		return "", -1, err
	}
	return id, role, nil
}
