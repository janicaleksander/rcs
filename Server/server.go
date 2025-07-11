package Server

import (
	"context"
	"fmt"
	"github.com/anthdm/hollywood/actor"
	db "github.com/janicaleksander/bcs/Database"
	"github.com/janicaleksander/bcs/Proto"
	"log/slog"
	"net"
	"os"
	"reflect"
)

var Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

// actor with remote
type Server struct {
	storage     db.Storage
	listenAddr  string // IP of (one) main server
	ln          net.Listener
	connections map[*actor.PID]net.Conn // Units PID to Unit IP addresses
}

func NewServer(listenAddr string, storage db.Storage) actor.Producer {
	return func() actor.Receiver {
		return &Server{
			storage:     storage,
			listenAddr:  listenAddr,
			connections: make(map[*actor.PID]net.Conn),
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
	case *Proto.NeedServerConfiguration:
		ctx.Respond(&Proto.NeededServerConfiguration{
			ServerPID: &Proto.PID{
				Address: ctx.PID().GetAddress(),
				Id:      ctx.PID().GetID(),
			}})
	case *Proto.LoginUser:
		s.loginUser(ctx, msg.Email, msg.Password)
	default:
		Logger.Warn("Server got unknown message", reflect.TypeOf(msg).String())

	}
}

func (s *Server) loginUser(ctx *actor.Context, email, password string) {
	c := context.Background()
	err := s.storage.LoginUser(c, email, password)
	fmt.Println(err)
	if err != nil {
		ctx.Respond(&Proto.Deny{})
	} else {
		ctx.Respond(&Proto.Accept{})
	}
}
