package main

import (
	"flag"
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/janicaleksander/bcs/Application"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		Server.Logger.Error("Error with loading .env file")
		return
	}
	appAddrFlag := flag.String("address", "", "Type here IP address of application")
	flag.Parse()
	if len(*appAddrFlag) <= 0 {
		Server.Logger.Error("Type value of flag")
		return
	}
	//Setup remote access
	r := remote.New(*appAddrFlag, remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		Server.Logger.Error(err.Error())
		return
	}

	Server.Logger.Info("Application is running on:", "Addr:", os.Getenv("APP_ADDR"))
	serverPID := actor.NewPID(os.Getenv("SERVER_ADDR"), "server/primary")
	//ping server
	resp := e.Request(serverPID, &Proto.IsServerRunning{}, Application.WaitTime)
	_, err = resp.Result()
	if err != nil {
		Server.Logger.Error("Server is not running", "err: ", err)
		return
	}
	window := Application.NewWindow()
	app := Application.NewWindowActor(window, serverPID)
	e.Spawn(app, "app") //this is creating new app

	//window
	window.RunWindow()

	//we need a block because actor is non-blocking,
	//so after we close application we close chanel and return from this
	<-window.Done
}
