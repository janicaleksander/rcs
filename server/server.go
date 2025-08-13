package server

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/anthdm/hollywood/actor"
	db "github.com/janicaleksander/bcs/database"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/utils"
)

const (
	PingPongTime = 3 * time.Second
)

type Session struct {
	userID string
	//presenceplace?
}

// TODO make a two maybe maps one for connected by app commander lvl person
// and second to 0 1 2 users/soldiers connected by device->unit->server
// this is problem e.g. in sending message from 5lvl from app to 0lvl to device in unit
// but this 5 see this 0 in his lists of person in system (maybe)
type Server struct {
	storage            db.Storage
	listenAddr         string                // IP of (one) main server
	connections        map[string]*actor.PID // uuid to PID
	reverseConnections map[string]string     //PIDstring to -> uuid
}

func NewServer(listenAddr string, storage db.Storage) actor.Producer {
	return func() actor.Receiver {
		return &Server{
			storage:            storage,
			listenAddr:         listenAddr,
			connections:        make(map[string]*actor.PID),
			reverseConnections: make(map[string]string),
		}
	}
}

// TODO In the future legal way to disconnect
// TODO manage if err !=nil some log errors
func (s *Server) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		utils.Logger.Info("server initialized")

	case actor.Started:
		utils.Logger.Info("server has started")
		ctx.SendRepeat(ctx.PID(), &proto.HeartbeatTick{}, 10*time.Second)
	case actor.Stopped:
		utils.Logger.Info("server has stopped")
	case *proto.HeartbeatTick:
		s.startHeartbeat(ctx)
	case *proto.IsServerRunning:
		ctx.Respond(&proto.Running{})
	case *proto.Disconnect: // after this switch state to loginScene
		pid, ok := s.connections[msg.Id]
		if ok {
			delete(s.connections, msg.Id)
			delete(s.reverseConnections, pid.String())
		}

	case *proto.IsOnline:
		if _, ok := s.connections[msg.Uuid]; ok {
			ctx.Respond(&proto.Online{})
		} else {
			ctx.Respond(&proto.Offline{})
		}
	case *proto.LoginUnit:
		//update use map
		//id, err := s.loginUnit(ctx, msg.Email, msg.Password)
		//if err == nil {
		//	s.clients[id] = ctx.Sender().GetAddress()
		//}
	case *proto.LoginUser:
		id, role, err := s.loginUser(msg.Email, msg.Password)
		if err != nil {
			ctx.Respond(&proto.DenyLogin{Info: err.Error()})
		} else {
			pid := actor.NewPID(msg.Pid.Address, msg.Pid.Id) //client PID
			s.connections[id] = pid                          //pid to uuid
			s.reverseConnections[pid.String()] = id
			ctx.Respond(&proto.AcceptLogin{Id: id, RuleLevel: int64(role)})

		}
		//TODO idk if this getlogged works

	case *proto.GetLoggedInUUID: //returning an uuid of current logged-in user
		pid := actor.NewPID(msg.Pid.Address, msg.Pid.Id) //client PID
		id := s.reverseConnections[pid.String()]
		ctx.Respond(&proto.LoggedInUUID{Id: id})
	case *proto.GetUserAboveLVL:
		c := context.Background()
		users, err := s.storage.GetUsersWithLVL(c, int(msg.Lvl))
		if err == nil {
			ctx.Respond(&proto.UsersAboveLVL{Users: users})
		}
	case *proto.CreateUnit:
		c := context.Background()
		err := s.storage.InsertUnit(c, msg.Name, msg.IsConfigured, msg.UserID)
		if err != nil {
			ctx.Respond(&proto.DenyCreateUnit{Info: err.Error()})
		} else {
			ctx.Respond(&proto.AcceptCreateUnit{})
		}
	case *proto.GetAllUnits:
		c := context.Background()
		units, err := s.storage.GetAllUnits(c)
		if err == nil {
			ctx.Respond(&proto.AllUnits{Units: units})
		}

	case *proto.GetAllUsersInUnit:
		c := context.Background()
		unitID := msg.Id
		users, err := s.storage.GetUsersInUnit(c, unitID)
		if err == nil {
			fmt.Println(users)
			ctx.Respond(&proto.AllUsersInUnit{Users: users})
		}
	case *proto.CreateUser:
		c := context.Background()
		err := s.storage.InsertUser(c, msg.User)
		if err != nil {
			ctx.Respond(&proto.DenyCreateUser{Info: err.Error()})
		} else {
			ctx.Respond(&proto.AcceptCreateUser{})
		}
	case *proto.IsUserInUnit:
		c := context.Background()
		isInUnit, unitID, err := s.storage.IsUserInUnit(c, msg.Id)
		if err != nil {
			ctx.Respond(&proto.UserNotInUnit{})
		} else if isInUnit {
			ctx.Respond(&proto.UserInUnit{UnitID: unitID})
		} else {
			ctx.Respond(&proto.UserNotInUnit{})
		}
	case *proto.AssignUserToUnit:
		c := context.Background()
		err := s.storage.AssignUserToUnit(c, msg.UserID, msg.UnitID)
		if err != nil {
			ctx.Respond(&proto.FailureOfAssign{})
		} else {
			ctx.Respond(&proto.SuccessOfAssign{})
		}
	case *proto.DeleteUserFromUnit:
		c := context.Background()
		err := s.storage.DeleteUserFromUnit(c, msg.UserID, msg.UnitID)
		if err != nil {
			ctx.Respond(&proto.FailureOfDelete{})
		} else {
			ctx.Respond(&proto.SuccessOfDelete{})
		}

	default:
		utils.Logger.Warn("server got unknown message", reflect.TypeOf(msg).String())

	}
}

func (s *Server) startHeartbeat(ctx *actor.Context) {
	for ID, PID := range s.connections {
		go func(pid *actor.PID, id string) {
			resp := ctx.Request(pid, &proto.Ping{}, utils.WaitTime)
			res, err := resp.Result()
			_, ok := res.(*proto.Pong)
			if !ok || err != nil {
				utils.Logger.Error("User is not responding for some time:", id, ctx.PID().String())
				ctx.Send(ctx.PID(), &proto.Disconnect{Id: id})
			}
		}(PID, ID)
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
