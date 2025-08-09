package main

import (
	"flag"
	"os"

	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/janicaleksander/bcs/application"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/server"
	"github.com/janicaleksander/bcs/utils"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		server.Logger.Error("Error with loading .env file")
		return
	}
	appAddrFlag := flag.String("address", "", "Type here IP address of application")
	flag.Parse()
	if len(*appAddrFlag) <= 0 {
		server.Logger.Error("Type value of flag")
		return
	}
	//Setup remote access
	r := remote.New(*appAddrFlag, remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		server.Logger.Error(err.Error())
		return
	}

	server.Logger.Info("application is running on:", "Addr:", os.Getenv("APP_ADDR"))
	messageServicePID := actor.NewPID(os.Getenv("MESSAGE_SERVICE_ADDR"), "messageService/primary")
	serverPID := actor.NewPID(os.Getenv("SERVER_ADDR"), "server/primary")
	//ping server
	resp := e.Request(serverPID, &proto.IsServerRunning{}, utils.WaitTime)
	_, err = resp.Result()
	if err != nil {
		server.Logger.Error("server is not running", "err: ", err)
		return
	}
	window := application.NewWindow()
	app := application.NewWindowActor(window, serverPID, messageServicePID)
	e.Spawn(app, "app") //this is creating new app

	//window
	window.RunWindow()

	//we need a block because actor is non-blocking,
	//so after we close application we close chanel and return from this
	<-window.Done
}
