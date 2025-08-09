package main

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/janicaleksander/bcs/Application"
	"github.com/janicaleksander/bcs/External"
	"github.com/janicaleksander/bcs/Proto"
	"github.com/janicaleksander/bcs/Server"
	"github.com/janicaleksander/bcs/Unit"
	"github.com/janicaleksander/bcs/Utils"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		Server.Logger.Error("Error with loading .env file")
		return
	}

	r := remote.New(os.Getenv("UNIT_ADDR"), remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		Server.Logger.Error(err.Error())
		return
	}

	Server.Logger.Info("Unit is running on:", "Addr:", os.Getenv("UNIT_ADDR"))
	serverPID := actor.NewPID(os.Getenv("SERVER_ADDR"), "server/primary")

	//ping server
	resp := e.Request(serverPID, &Proto.IsServerRunning{}, Utils.WaitTime)
	_, err = resp.Result()
	if err != nil {
		Server.Logger.Error("Servers is not running", "err: ", err)
		return
	}
	ext := External.NewExternal()
	unit := Unit.NewUnit(serverPID, ext)

	e.Spawn(unit, "unit") //this is creating new unit

	// TODO: New idea is to have CLI to login unit to server by e.g. sending unit id
	//	e.Send(unitPID, &Proto.LoginUnit{Id: "Id from venv"})
	select {}
}
