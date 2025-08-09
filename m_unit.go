package main

import (
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/janicaleksander/bcs/Unit"
	"github.com/janicaleksander/bcs/application"
	"github.com/janicaleksander/bcs/external"
	"github.com/janicaleksander/bcs/proto"
	"github.com/janicaleksander/bcs/server"
	"github.com/janicaleksander/bcs/utils"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		server.Logger.Error("Error with loading .env file")
		return
	}

	r := remote.New(os.Getenv("UNIT_ADDR"), remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		server.Logger.Error(err.Error())
		return
	}

	server.Logger.Info("Unit is running on:", "Addr:", os.Getenv("UNIT_ADDR"))
	serverPID := actor.NewPID(os.Getenv("SERVER_ADDR"), "server/primary")

	//ping server
	resp := e.Request(serverPID, &proto.IsServerRunning{}, utils.WaitTime)
	_, err = resp.Result()
	if err != nil {
		server.Logger.Error("Servers is not running", "err: ", err)
		return
	}
	ext := external.NewExternal()
	unit := unit.NewUnit(serverPID, ext)

	e.Spawn(unit, "unit") //this is creating new unit

	// TODO: New idea is to have CLI to login unit to server by e.g. sending unit id
	//	e.Send(unitPID, &proto.LoginUnit{Id: "Id from venv"})
	select {}
}
