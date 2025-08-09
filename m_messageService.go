package main

import (
	"os"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/janicaleksander/bcs/messageservice"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/server"
	"github.com/janicaleksander/bcs/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		server.Logger.Error("Can't load .env file")
		return
	}

	r := remote.New(os.Getenv("MESSAGE_SERVICE_ADDR"), remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		server.Logger.Error("Error with engine configuration")
		return
	}

	serverPID := actor.NewPID(os.Getenv("SERVER_ADDR"), "server/primary")
	//ping server
	resp := e.Request(serverPID, &proto.IsServerRunning{}, utils.WaitTime)
	_, err = resp.Result()
	if err != nil {
		server.Logger.Error("server is not running", "err: ", err)
		return
	}
	server.Logger.Info("messageservice is running on: ", "Address: ", os.Getenv("MESSAGE_SERVICE_ADDR"))
	messageService := messageservice.NewMessageService(serverPID)
	e.Spawn(messageService, "messageService", actor.WithID("primary")) //this is creating new app

	//we need a block because actor is non-blocking,
	//so after we close application we close chanel and return from this
	select {}
}
