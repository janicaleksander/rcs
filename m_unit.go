package main

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/janicaleksander/bcs/External"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	"github.com/janicaleksander/bcs/Unit"
	"github.com/joho/godotenv"
	"os"
	"time"
)

func main() {
	//some way to first configure a unit e.g. created by general,

	//we can assume that we provide as VENV uuid of unit
	err := godotenv.Load()
	if err != nil {
		Server.Logger.Error("Error loading .env file")
		return
	}

	ext := External.NewExternal() //maybe its inside unit?
	unit := Unit.NewUnit("some id from venv", ext)

	r := remote.New(os.Getenv("UNIT_ADDR"), remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		Server.Logger.Error(err.Error())
		return
	}

	Server.Logger.Info("Unit is running on:", "Addr:", os.Getenv("UNIT_ADDR"))
	serverPID := actor.NewPID(os.Getenv("SERVER_ADDR"), "server/primary")

	//ping server
	resp := e.Request(serverPID, &Proto.IsServerRunning{}, time.Second)
	_, err = resp.Result()
	if err != nil {
		Server.Logger.Error("Servers is not running", "err: ", err)
		return
	}

	// TODO: prob in the future when server is down
	// still we can operate on unit
	unitPID := e.Spawn(unit, "unit") //this is creating new unit

	resp = e.Request(serverPID, &Proto.ConnectToServer{
		Client: &Proto.PID{
			Address: unitPID.GetAddress(),
			Id:      unitPID.GetID(),
		},
	}, time.Second)
	val, err := resp.Result()
	if err != nil {
		Server.Logger.Error("Can't connect to the server!", "err: ", err)
		return
	}
	//Respond to ConnectToServer neededServerConfiguration
	e.Send(unitPID, val)

	e.Send(unitPID, &Proto.LoginUnit{Id: "Id from venv"})
	select {}
}
