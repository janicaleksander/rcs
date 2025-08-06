package main

import (
	"flag"
	"os"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/janicaleksander/bcs/Application"
	"github.com/janicaleksander/bcs/MessageService"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		Server.Logger.Error("Can't load .env file")
		return
	}
	messageServiceAddress := flag.String("address", "", "Type a IP address of this service")
	flag.Parse()
	if len(*messageServiceAddress) <= 0 {
		Server.Logger.Error("Type a flag")
		return
	}
	r := remote.New(*messageServiceAddress, remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		Server.Logger.Error("Error with engine configuration")
		return
	}

	serverPID := actor.NewPID(os.Getenv("SERVER_ADDR"), "server/primary")
	//ping server
	resp := e.Request(serverPID, &Proto.IsServerRunning{}, Application.WaitTime)
	_, err = resp.Result()
	if err != nil {
		Server.Logger.Error("Server is not running", "err: ", err)
		return
	}
	Server.Logger.Info("MessageService is running on: ", "Address: ", *messageServiceAddress)
	messageService := MessageService.NewMessageService(serverPID)
	e.Spawn(messageService, "messageService") //this is creating new app

	//we need a block because actor is non-blocking,
	//so after we close application we close chanel and return from this
	select {}
}
