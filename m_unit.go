package main

import (
	"fmt"
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
	err := godotenv.Load()
	if err != nil {
		Server.Logger.Error("Error loading .env file")
		return
	}
	ext := External.NewExternal()
	unit := Unit.NewUnit(ext)
	r := remote.New(os.Getenv("127.0.0.1:2002"), remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		Server.Logger.Error(err.Error())
		return
	}
	Server.Logger.Info("Unit is running on: ", "Addr: ", "127.0.0.1:2002")
	e.Spawn(unit, "unit")
	pid := actor.NewPID(os.Getenv("SERVER_ADDR"), "server/primary")
	for {
		e.Send(pid, &Proto.Req{Name: "x"})
		fmt.Println(pid.String())
		time.Sleep(time.Second * 2)
	}
	select {}
}
