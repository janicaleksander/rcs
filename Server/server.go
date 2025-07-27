package Server

import (
	"context"
	"fmt"
	"github.com/anthdm/hollywood/actor"
	db "github.com/janicaleksander/bcs/Database"
	"github.com/janicaleksander/bcs/Proto"
	"log/slog"
	"os"
	"reflect"
)

// TODO add logger statement everywhere
var Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

// actor with remote
type Server struct {
	storage     db.Storage
	listenAddr  string                // IP of (one) main server
	connections map[*actor.PID]string // PID  to uuid
}

// TODO do refactor of Receive loop to make distinct functions
func NewServer(listenAddr string, storage db.Storage) actor.Producer {
	return func() actor.Receiver {
		return &Server{
			storage:     storage,
			listenAddr:  listenAddr,
			connections: make(map[*actor.PID]string),
		}
	}
}

func (s *Server) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Initialized:
		Logger.Info("Server initialized")
	case actor.Started:
		Logger.Info("Server has started")
	case actor.Stopped:
		Logger.Info("Server has stopped")
	//case to update connection map in connection/disconnection
	case *Proto.IsServerRunning:
		ctx.Respond(&Proto.Running{})
	case *Proto.PingServer:
		// TODO: implement maps of connections to server
		// maybe login could only add sth to connections map (app)
		// needs his ID (unit)
		//respond to get others PID of server
		ctx.Respond(&Proto.NeededServerConfiguration{
			ServerPID: &Proto.PID{
				Address: ctx.PID().GetAddress(),
				Id:      ctx.PID().GetID(),
			}})
	case *Proto.Disconnect:
		_, ok := s.connections[ctx.Sender()]
		if !ok {
			//sth
		}
		delete(s.connections, ctx.Sender())
	case *Proto.LoginUnit:
		//update use map
		//id, err := s.loginUnit(ctx, msg.Email, msg.Password)
		//if err == nil {
		//	s.clients[id] = ctx.Sender().GetAddress()
		//}
	case *Proto.LoginUser:
		id, err := s.loginUser(ctx, msg.Email, msg.Password)
		if err == nil {
			pid := actor.NewPID(msg.Pid.Address, msg.Pid.Id) //client PID
			_, ok := s.connections[pid]
			if ok {
				//TODO repair this
				fmt.Println("byłem")
			}
			s.connections[pid] = id //pid to uuid
		}
	case *Proto.GetLoggedInUUID:
		pid := actor.NewPID(msg.Pid.Address, msg.Pid.Id) //client PID
		id := s.connections[pid]
		ctx.Respond(&Proto.LoggedInUUID{Id: id})

	case *Proto.GetUserAboveLVL:
		c := context.Background()
		users, err := s.storage.GetUsersWithLVL(c, 4)
		if err == nil {
			ctx.Respond(&Proto.GetUserAboveLVL{Users: users})
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

	default:
		Logger.Warn("Server got unknown message", reflect.TypeOf(msg).String())

	}
}

func (s *Server) loginUser(ctx *actor.Context, email, password string) (string, error) {
	// TODO: add jwt to login
	c := context.Background()
	id, role, err := s.storage.LoginUser(c, email, password)
	if err != nil {
		ctx.Respond(&Proto.DenyLogin{Info: err.Error()})
		return "", err
	} else {
		ctx.Respond(&Proto.AcceptLogin{Info: "Login successful! ", RuleLevel: int64(role)})
	}
	return id, nil
}
