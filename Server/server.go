package Server

import (
	"context"
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
	storage    db.Storage
	listenAddr string // IP of (one) main server
	//	ln          net.Listener
	connections map[string]*actor.PID // ip address to PID
	clients     map[string]string     // uuid (user/unit)  to address
}

func NewServer(listenAddr string, storage db.Storage) actor.Producer {
	return func() actor.Receiver {
		return &Server{
			storage:     storage,
			listenAddr:  listenAddr,
			connections: make(map[string]*actor.PID),
			clients:     make(map[string]string),
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
	case *Proto.ConnectToServer:
		// TODO: implement maps of connections to server
		// maybe login could only add sth to connections map (app)
		// needs his ID (unit)
		pid := actor.NewPID(msg.Client.Address, msg.Client.Id)
		s.connections[pid.GetAddress()] = pid

		//respond to get others PID of server
		ctx.Respond(&Proto.NeededServerConfiguration{
			ServerPID: &Proto.PID{
				Address: ctx.PID().GetAddress(),
				Id:      ctx.PID().GetID(),
			}})
	case *Proto.Disconnect:
		_, ok := s.connections[ctx.Sender().GetAddress()]
		if !ok {
			//sth
		}
		delete(s.connections, ctx.Sender().GetAddress())
	case *Proto.LoginUnit:
		//update use map
		//id, err := s.loginUnit(ctx, msg.Email, msg.Password)
		//if err == nil {
		//	s.clients[id] = ctx.Sender().GetAddress()
		//}
	case *Proto.LoginUser:
		id, err := s.loginUser(ctx, msg.Email, msg.Password)
		if err == nil {
			pid := actor.NewPID(msg.Pid.Address, msg.Pid.Id)
			s.clients[id] = pid.GetAddress()
		}
	default:
		Logger.Warn("Server got unknown message", reflect.TypeOf(msg).String())

	}
}

func (s *Server) loginUser(ctx *actor.Context, email, password string) (string, error) {
	c := context.Background()
	err := s.storage.LoginUser(c, email, password)
	if err != nil {
		ctx.Respond(&Proto.Deny{})
	} else {
		ctx.Respond(&Proto.Accept{})
	}
	return "here ID from DB", nil
}
