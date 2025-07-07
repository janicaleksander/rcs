package Server

import (
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
	case *Proto.Req:
		Logger.Info("Server got unknown message", msg, reflect.TypeOf(msg).String())

	default:
		fmt.Println(msg)
		Logger.Warn("Server got unknown message", reflect.TypeOf(msg).String())
		_ = msg

	}
}
