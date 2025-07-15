package main

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	"github.com/janicaleksander/bcs/application"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		Server.Logger.Error("Error loading .env file")
		return
	}
	//chan's part to send data between app actor and window:
	ChloginUser := make(chan *Proto.LoginUser, 1024)

	app := application.NewApp(ChloginUser)
	r := remote.New(os.Getenv("APP_ADDR"), remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		Server.Logger.Error(err.Error())
		return
	}

	Server.Logger.Info("App is running on:", "Addr:", os.Getenv("APP_ADDR"))
	serverPID := actor.NewPID(os.Getenv("SERVER_ADDR"), "server/primary")

	//ping server
	resp := e.Request(serverPID, &Proto.IsServerRunning{}, time.Second)
	_, err = resp.Result()
	if err != nil {
		Server.Logger.Error("Server is not running", "err: ", err)
		return
	}

	appPID := e.Spawn(app, "app") //this is creating new app

	resp = e.Request(serverPID, &Proto.ConnectToServer{
		Client: &Proto.PID{
			Address: appPID.GetAddress(),
			Id:      appPID.GetID(),
		},
	}, time.Second)
	val, err := resp.Result()
	if err != nil {
		Server.Logger.Error("Can't connect to the server!", "err: ", err)
		return
	}
	//respond to connect to server neededServerConfiguration
	e.Send(appPID, val)

	//window
	window := application.NewWindow(ChloginUser)
	window.RunWindow()

	//running here -> first scene is loading bar and change to loginPanel only if ping to server works
	<-window.Done
}
