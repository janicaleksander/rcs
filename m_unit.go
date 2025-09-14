package main

import (
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	db "github.com/janicaleksander/bcs/database"
	"github.com/janicaleksander/bcs/external/unit"
	"github.com/janicaleksander/bcs/types/proto"
	"github.com/janicaleksander/bcs/utils"
)

func main() {
	//Setup remote access
	config := struct {
		Server struct {
			Addr       string `toml:"addr"`
			ServerAddr string `toml:"serverAddr"`
			UnitID     string `toml:"unitID"`
		}
	}{}
	_, err := toml.DecodeFile("configproduction/unit.toml", &config)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	if len(strings.TrimSpace(config.Server.Addr)) == 0 ||
		len(strings.TrimSpace(config.Server.ServerAddr)) == 0 ||
		len(strings.TrimSpace(config.Server.UnitID)) == 0 {
		utils.Logger.Error("bad unit cfg file")
		return
	}
	r := remote.New(config.Server.ServerAddr, remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	utils.Logger.Info("unit is running on:", "Addr:", config.Server.Addr)
	serverPID := actor.NewPID(config.Server.ServerAddr, "server/primary")
	//ping server
	resp := e.Request(serverPID, &proto.IsServerRunning{}, utils.WaitTime)
	_, err = resp.Result()
	if err != nil {
		utils.Logger.Error("server is not running", "err: ", err)
		return
	}
	dbManager, err := db.GetDBManager(db.WithConnectionTimeout(10))
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}
	dbase := dbManager.GetDB()
	pg := &db.Postgres{Conn: dbase}
	_ = e.Spawn(unit.NewUnit(config.Server.UnitID, serverPID, pg), "device", actor.WithID(config.Server.UnitID))
	select {}
}
