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
	"sync"
	"time"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

const (
	PingPingTime = 3 * time.Second
)

type Server struct {
	storage     db.Storage
	listenAddr  string                // IP of (one) main server
	connections map[*actor.PID]string // PID  to uuid
}

func NewServer(listenAddr string, storage db.Storage) actor.Producer {
	return func() actor.Receiver {
		return &Server{
			storage:     storage,
			listenAddr:  listenAddr,
			connections: make(map[*actor.PID]string),
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
		go s.heartbeat(ctx)
	case actor.Stopped:
		Logger.Info("Server has stopped")
	case *Proto.IsServerRunning:
		ctx.Respond(&Proto.Running{})
	case *Proto.Disconnect: // after this switch state to loginScene
		_, ok := s.connections[ctx.Sender()]
		if ok {
			delete(s.connections, ctx.Sender())
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
			s.connections[pid] = id                          //pid to uuid
			ctx.Respond(&Proto.AcceptLogin{Info: "Login successful! ", RuleLevel: int64(role)})
		}
	case *Proto.GetLoggedInUUID: //returning an uuid of current logged-in user/unit
		pid := actor.NewPID(msg.Pid.Address, msg.Pid.Id) //client PID
		id := s.connections[pid]
		ctx.Respond(&Proto.LoggedInUUID{Id: id})

	case *Proto.GetUserAboveLVL:
		c := context.Background()
		users, err := s.storage.GetUsersWithLVL(c, 4)
		if err == nil {
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

	default:
		Logger.Warn("Server got unknown message", reflect.TypeOf(msg).String())

	}
}
func (s *Server) heartbeat(ctx *actor.Context) {
	var mutex sync.RWMutex

	for {
		time.Sleep(5 * time.Second)
		for pid := range s.connections {
			go func(p *actor.PID) {
				resp := ctx.Request(p, &Proto.Ping{}, PingPingTime)
				v, err := resp.Result()
				if _, ok := v.(*Proto.Pong); !ok || err != nil {
					mutex.Lock()
					if _, exists := s.connections[p]; exists {
						delete(s.connections, p)
					}
					mutex.Unlock()
				}
			}(pid)
		}
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
