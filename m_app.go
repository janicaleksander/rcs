package main

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	cli_app "github.com/janicaleksander/bcs/cli-app"
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
	app := cli_app.NewCLI()
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

	resp = e.Request(serverPID, &Proto.NeedServerConfiguration{}, time.Second)
	val, err := resp.Result()
	if err != nil {
		Server.Logger.Error("Can't do the request!", "err: ", err)
		return
	}
	//neededServerConfiguration
	e.Send(appPID, val)

	//run APP/CLI
	e.Send(appPID, &Proto.StartCLI{})

	select {}
}
